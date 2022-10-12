package dlqdump

import (
	"sync"
	"sync/atomic"

	"github.com/koykov/bitset"
	"github.com/koykov/queue"
)

const (
	flagTimer = 0
)

// Queue represents dumping queue.
type Queue struct {
	bitset.Bitset
	// Config instance.
	config *Config
	// Actual queue status.
	status queue.Status

	once sync.Once

	// Internal timer. Triggers flush operation according Config.FlushInterval.
	timer *timer

	mux sync.Mutex
	buf []byte

	Err error
}

// NewQueue makes new dumping queue instance and initialize it according config params.
func NewQueue(config *Config) (*Queue, error) {
	q := &Queue{
		config: config.Copy(),
	}
	q.once.Do(q.init)
	return q, q.Err
}

// Enqueue puts x to the queue.
func (q *Queue) Enqueue(x interface{}) (err error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == queue.StatusClose || status == queue.StatusFail {
		return queue.ErrQueueClosed
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	// Encode item to bytes.
	q.buf, err = q.c().Encoder.Encode(q.buf[:0], x)
	if err != nil {
		return
	}

	// Start timer on first incoming item.
	// Timer will trigger flush operation after Config.FlushInterval since current time.
	if !q.CheckBit(flagTimer) {
		q.SetBit(flagTimer, true)
		go q.timer.wait(q)
	}

	// Forward encoded item to writer.
	if _, err = q.c().Writer.Write(q.c().Version, q.buf); err != nil {
		q.m().Fail("write fail")
		return
	}
	q.m().Dump(len(q.buf))

	q.buf = q.buf[:0]

	// Check if Config.Capacity reached.
	if q.c().Writer.Size() >= q.c().Capacity {
		// Reset timer and flush with corresponding reason.
		q.timer.reset()
		err = q.flushLF(flushReasonSize)
	}

	return
}

// Size returns actual size in bytes of all queued items (since start or last flush).
func (q *Queue) Size() int {
	if q.c().Writer == nil {
		return 0
	}
	return int(q.c().Writer.Size())
}

// Capacity returns maximum queue capacity.
func (q *Queue) Capacity() int {
	return int(q.c().Capacity)
}

// Rate returns size to capacity ratio.
func (q *Queue) Rate() float32 {
	return 0
}

// Close gracefully stops the queue.
func (q *Queue) Close() error {
	if q.getStatus() == queue.StatusClose {
		return queue.ErrQueueClosed
	}

	if l := q.l(); l != nil {
		msg := "caught close signal"
		l.Printf(msg)
	}

	q.mux.Lock()
	defer q.mux.Unlock()
	q.timer.stop()
	return q.flushLF(flushReasonForce)
}

// Init the queue.
func (q *Queue) init() {
	c := q.c()

	// Check mandatory params.
	if c.Capacity == 0 {
		q.Err = queue.ErrNoSize
		q.setStatus(queue.StatusFail)
		return
	}
	if c.Encoder == nil {
		q.Err = ErrNoEncoder
		q.setStatus(queue.StatusFail)
		return
	}
	if c.Writer == nil {
		q.Err = ErrNoWriter
		q.setStatus(queue.StatusFail)
		return
	}

	// Check non-mandatory params and set default values if needed.
	if c.FlushInterval == 0 {
		c.FlushInterval = defaultFlushInterval
	}
	q.timer = newTimer()

	if c.MetricsWriter == nil {
		// Use dummy MW.
		c.MetricsWriter = DummyMetrics{}
	}

	// Queue is ready!
	q.setStatus(queue.StatusActive)
}

// Set status of the queue.
func (q *Queue) setStatus(status queue.Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

// Get status of the queue.
func (q *Queue) getStatus() queue.Status {
	return queue.Status(atomic.LoadUint32((*uint32)(&q.status)))
}

func (q *Queue) c() *Config {
	return q.config
}

func (q *Queue) m() MetricsWriter {
	return q.config.MetricsWriter
}

func (q *Queue) l() queue.Logger {
	return q.config.Logger
}
