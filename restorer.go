package dlqdump

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/queue"
)

// Restorer represents dump restore handler.
// Restorer may be scheduled (see Config.CheckInterval).
type Restorer struct {
	config *Config
	status queue.Status

	once sync.Once
	lock uint32
	buf  []byte

	Err error
}

// NewRestorer makes new restorer instance and initialize it according config params.
func NewRestorer(config *Config) (*Restorer, error) {
	r := &Restorer{
		config: config.Copy(),
	}
	r.once.Do(r.init)
	return r, nil
}

// Restore makes an attempt of restoring operation.
//
// Restore reads the dump and forwards each item to the destination queue. On Close it stops reading new
// items, but the item that has already been read is still delivered: Close is synchronous and waits until
// the in-flight item is sent (see CloseWithTimeout/CloseWithContext).
func (r *Restorer) Restore() error {
	r.once.Do(r.init)
	if status := r.getStatus(); status == queue.StatusClose || status == queue.StatusFail {
		return queue.ErrQueueClosed
	}

	if !atomic.CompareAndSwapUint32(&r.lock, 0, 1) {
		// Another restore attempt is in progress.
		return nil
	}
	defer atomic.StoreUint32(&r.lock, 0)

	var (
		err error
		ver Version
	)
	for {
		// Once closing, stop reading new items. The item read on the previous iteration has already been
		// delivered, so nothing is lost.
		if r.getStatus() == queue.StatusClose {
			break
		}

		// Check reader for new encoded items.
		r.buf = r.buf[:0]
		ver, r.buf, err = r.config.Reader.Read(r.buf)
		if err == io.EOF {
			// EOF reached, finish current restore attempt.
			err = nil
			break
		}
		if err != nil {
			r.config.MetricsWriter.Fail("read error")
			continue
		}
		if ver != r.config.Version {
			r.config.MetricsWriter.Fail("version mismatch")
			continue
		}

		// Decode item.
		var x any
		if x, err = r.config.Decoder.Decode(r.buf); err != nil {
			r.config.MetricsWriter.Fail("decode error")
			continue
		}
		// Wait until the destination queue can accept the already-read item. On Close the item is not
		// dropped: the wait continues (bounded by CloseWithTimeout/CloseWithContext) until it is delivered.
		for r.config.Queue.Rate() > r.config.AllowRate {
			time.Sleep(r.config.PostponeInterval)
		}
		// Put item to the destination queue.
		if err = r.config.Queue.Enqueue(x); err != nil {
			r.config.MetricsWriter.Fail("enqueue fail")
			break
		}
		r.config.MetricsWriter.Restore(len(r.buf))
	}
	return nil
}

// Close gracefully stops the restorer.
func (r *Restorer) Close() error {
	return r.CloseWithTimeout(r.config.CloseTimeout)
}

// CloseWithTimeout stops the restorer with timeout.
func (r *Restorer) CloseWithTimeout(timeout time.Duration) error {
	// Signal an in-flight Restore() to stop reading new items once the current item is delivered.
	r.setStatus(queue.StatusClose)
	deadline := time.Now().Add(timeout)
	for atomic.LoadUint32(&r.lock) == 1 {
		if time.Now().After(deadline) {
			return ErrTimeout
		}
		time.Sleep(time.Millisecond)
	}
	return r.config.Reader.Close()
}

// CloseWithContext stops the restorer and respects ctx cancellation while waiting
// for an in-flight Restore() to finish.
func (r *Restorer) CloseWithContext(ctx context.Context) error {
	r.setStatus(queue.StatusClose)
	for atomic.LoadUint32(&r.lock) == 1 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond):
		}
	}
	return r.config.Reader.Close()
}

// ForceClose immediately stops the queue.
func (r *Restorer) ForceClose() error {
	r.setStatus(queue.StatusClose)
	return nil
}

// Init the restorer.
func (r *Restorer) init() {
	c := r.config

	// Check mandatory params.
	if c.Decoder == nil {
		r.Err = ErrNoDecoder
		r.setStatus(queue.StatusFail)
		return
	}
	if c.Reader == nil {
		r.Err = ErrNoReader
		r.setStatus(queue.StatusFail)
		return
	}
	if c.Queue == nil {
		r.Err = ErrNoQueue
		r.setStatus(queue.StatusFail)
		return
	}

	// Check non-mandatory params and set default values if needed.
	if c.CheckInterval == 0 {
		c.CheckInterval = defaultCheckInterval
	}
	if c.PostponeInterval == 0 {
		c.PostponeInterval = c.CheckInterval
	}
	if c.AllowRate == 0 {
		c.AllowRate = defaultAllowRate
	}
	if c.CloseTimeout == 0 {
		c.CloseTimeout = defaultCloseTimeout
	}

	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetrics{}
	}

	// Restorer is ready!
	r.setStatus(queue.StatusActive)

	// Init background check ticker.
	ticker := time.NewTicker(c.CheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if r.getStatus() == queue.StatusClose {
					ticker.Stop()
					return
				}
				_ = r.Restore()
			}
		}
	}()
}

func (r *Restorer) setStatus(status queue.Status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *Restorer) getStatus() queue.Status {
	return queue.Status(atomic.LoadUint32((*uint32)(&r.status)))
}
