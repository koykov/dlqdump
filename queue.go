package dlqdump

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/koykov/queue"
)

// Queue represents dumping queue.
type Queue struct {
	// Config instance.
	config *Config
	// Actual queue status.
	status queue.Status
	// Flag that the flush timer is running. Accessed atomically.
	timerOn uint32

	once sync.Once

	// Internal timer. Triggers flush operation according Config.FlushInterval.
	timer *timer

	mux sync.Mutex
	buf []byte
	// Counter of in-flight Enqueue operations. Close waits for it to drop to zero, so all items that
	// have already started to be written are flushed before the queue returns from Close.
	wlock int64

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
func (q *Queue) Enqueue(x any) (err error) {
	q.once.Do(q.init)

	// Register the operation before checking the status: Close waits for the counter to reach zero, so a
	// racing Enqueue is either completed before the final flush or rejected here. This ordering prevents
	// an item from being written after Close has already flushed the queue.
	atomic.AddInt64(&q.wlock, 1)
	defer atomic.AddInt64(&q.wlock, -1)
	if status := q.getStatus(); status == queue.StatusClose || status == queue.StatusFail {
		return queue.ErrQueueClosed
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	// Check job incoming and extract payload.
	switch x.(type) {
	case queue.Job:
		x = x.(queue.Job).Payload
	case *queue.Job:
		x = x.(*queue.Job).Payload
	}

	// Encode item to bytes.
	q.buf, err = q.c().Encoder.Encode(q.buf[:0], x)
	if err != nil {
		return
	}

	// Arm the timer on the first incoming item after a flush. CAS guarantees a single armed timer.
	if atomic.CompareAndSwapUint32(&q.timerOn, 0, 1) {
		q.timer.wait(q)
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
		// Disarm the timer and flush with corresponding reason.
		q.timer.stop()
		atomic.StoreUint32(&q.timerOn, 0)
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
//
// Close is fully synchronous: it stops accepting new items, waits for all in-flight Enqueue operations
// to finish, then flushes everything accumulated in the queue to the storage. It returns only after the
// data is flushed, so nothing is lost.
func (q *Queue) Close() error {
	q.once.Do(q.init)
	if status := q.getStatus(); status == queue.StatusFail {
		return q.Err
	}

	// Stop accepting new items. CAS makes Close idempotent.
	if !atomic.CompareAndSwapUint32((*uint32)(&q.status), uint32(queue.StatusActive), uint32(queue.StatusClose)) {
		return queue.ErrQueueClosed
	}

	if l := q.l(); l != nil {
		l.Printf("caught close signal")
	}

	// Wait until all in-flight Enqueue operations are finished: their data must reach the Writer before
	// the final flush. The mutex isn't held here to avoid a deadlock with the Enqueue waiting on it.
	for atomic.LoadInt64(&q.wlock) > 0 {
		runtime.Gosched()
	}

	// Disarm the timer so it can't fire after the final flush, and flush under the mutex: this serializes
	// with a possibly in-flight timer flush (a late one is a no-op, Writer.Size() is zero after this flush).
	q.mux.Lock()
	q.timer.stop()
	err := q.flushLF(flushReasonForce)
	q.mux.Unlock()
	return err
}

// Init the queue.
func (q *Queue) init() {
	c := q.c()

	// Check mandatory params.
	if c.Capacity == 0 {
		q.Err = queue.ErrNoCapacity
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
