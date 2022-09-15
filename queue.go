package dlqdump

import (
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/koykov/blqueue"
	"github.com/koykov/bytealg"
)

type Queue struct {
	config *Config
	status blqueue.Status

	once  sync.Once
	timer *timer
	mux   sync.Mutex
	buf   []byte

	rw *rworker

	Err error
}

func New(config *Config) (*Queue, error) {
	q := &Queue{
		config: config.Copy(),
	}
	q.once.Do(q.init)
	return q, q.Err
}

func (q *Queue) Enqueue(x interface{}) (err error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == blqueue.StatusClose || status == blqueue.StatusFail {
		return blqueue.ErrQueueClosed
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	if len(q.buf) == 0 {
		q.buf = bytealg.GrowDelta(q.buf, 4)
		binary.LittleEndian.PutUint32(q.buf[:4], q.config.Version)
		go q.timer.wait(q)
	}

	ho := len(q.buf)
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	po := len(q.buf)

	q.buf, err = q.config.Encoder.Encode(q.buf, x)
	if err != nil {
		q.buf = q.buf[:ho]
		return
	}

	pl := len(q.buf) - po
	binary.LittleEndian.PutUint32(q.buf[ho:ho+4], uint32(pl))
	q.config.MetricsWriter.QueuePut(q.config.Key, pl)

	if MemorySize(len(q.buf)) >= q.config.Capacity {
		q.timer.reset()
		err = q.flushLF(flushReasonSize)
	}

	return
}

func (q *Queue) Size() int {
	if q.config.Dumper == nil {
		return 0
	}
	return int(q.config.Dumper.Size())
}

func (q *Queue) Capacity() int {
	return int(q.config.Capacity)
}

func (q *Queue) Rate() float32 {
	return 0
}

func (q *Queue) Close() error {
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

	if q.rw != nil {
		q.rw.stop()
	}

	return nil
}

func (q *Queue) init() {
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

func (q *Queue) setStatus(status blqueue.Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

func (q *Queue) getStatus() blqueue.Status {
	return blqueue.Status(atomic.LoadUint32((*uint32)(&q.status)))
}
