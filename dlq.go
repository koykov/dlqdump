package dlqdump

import (
	"sync"
	"sync/atomic"

	"github.com/koykov/bitset"
	"github.com/koykov/blqueue"
)

const (
	flagTimer = 0
)

type DLQ struct {
	bitset.Bitset
	config *DLQConfig
	status blqueue.Status

	once  sync.Once
	timer *timer
	mux   sync.Mutex
	buf   []byte

	Err error
}

func NewDLQ(config *DLQConfig) (*DLQ, error) {
	q := &DLQ{
		config: config.Copy(),
	}
	q.once.Do(q.init)
	return q, q.Err
}

func (q *DLQ) Enqueue(x interface{}) (err error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == blqueue.StatusClose || status == blqueue.StatusFail {
		return blqueue.ErrQueueClosed
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	q.buf, err = q.config.Encoder.Encode(q.buf[:0], x)
	if err != nil {
		return
	}

	if !q.CheckBit(flagTimer) {
		q.SetBit(flagTimer, true)
		go q.timer.wait(q)
	}

	if _, err = q.config.Dumper.Dump(q.config.Version, q.buf); err != nil {
		return
	}
	q.buf = q.buf[:0]

	if q.config.Dumper.Size() >= q.config.Capacity {
		q.timer.reset()
		err = q.flushLF(flushReasonSize)
	}

	return
}

func (q *DLQ) Size() int {
	if q.config.Dumper == nil {
		return 0
	}
	return int(q.config.Dumper.Size())
}

func (q *DLQ) Capacity() int {
	return int(q.config.Capacity)
}

func (q *DLQ) Rate() float32 {
	return 0
}

func (q *DLQ) Close() error {
	if q.getStatus() == blqueue.StatusClose {
		return blqueue.ErrQueueClosed
	}

	if l := q.config.Logger; l != nil {
		msg := "queue #%s caught close signal"
		l.Printf(msg, q.config.Key)
	}

	q.mux.Lock()
	defer q.mux.Unlock()
	q.timer.stop()
	if len(q.buf) > 4 {
		return q.flushLF(flushReasonForce)
	}

	return nil
}

func (q *DLQ) init() {
	c := q.config

	if len(c.Key) == 0 {
		q.Err = blqueue.ErrNoKey
		q.setStatus(blqueue.StatusFail)
		return
	}
	if c.Capacity == 0 {
		q.Err = blqueue.ErrNoSize
		q.setStatus(blqueue.StatusFail)
		return
	}
	if c.Encoder == nil {
		q.Err = ErrNoEncoder
		q.setStatus(blqueue.StatusFail)
		return
	}
	if c.Dumper == nil {
		q.Err = ErrNoDumper
		q.setStatus(blqueue.StatusFail)
		return
	}
	if c.FlushInterval == 0 {
		c.FlushInterval = defaultTimeLimit
	}
	q.timer = newTimer()

	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetrics{}
	}

	q.setStatus(blqueue.StatusActive)
}

func (q *DLQ) setStatus(status blqueue.Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

func (q *DLQ) getStatus() blqueue.Status {
	return blqueue.Status(atomic.LoadUint32((*uint32)(&q.status)))
}
