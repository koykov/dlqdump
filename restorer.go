package dlqdump

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/blqueue"
)

// Restorer represents dump restore handler.
// Restorer may be scheduled (see Config.CheckInterval).
type Restorer struct {
	config *Config
	status blqueue.Status

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
func (r *Restorer) Restore() error {
	r.once.Do(r.init)
	if status := r.getStatus(); status == blqueue.StatusClose || status == blqueue.StatusFail {
		return blqueue.ErrQueueClosed
	}

	if atomic.LoadUint32(&r.lock) == 1 {
		return nil
	}
	atomic.StoreUint32(&r.lock, 1)
	defer atomic.StoreUint32(&r.lock, 0)

	var (
		err error
		ver Version
	)
	for {
		if r.getStatus() == blqueue.StatusClose {
			return blqueue.ErrQueueClosed
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
			r.config.MetricsWriter.Fail(r.config.Key, "read error")
			continue
		}
		if ver != r.config.Version {
			r.config.MetricsWriter.Fail(r.config.Key, "version mismatch")
			continue
		}

		// Decode item.
		var x interface{}
		if x, err = r.config.Decoder.Decode(r.buf); err != nil {
			r.config.MetricsWriter.Fail(r.config.Key, "decode error")
			continue
		}
		// Spin until destination queue rate is too big.
		for r.config.Queue.Rate() > r.config.AllowRate {
			if r.getStatus() == blqueue.StatusClose {
				return blqueue.ErrQueueClosed
			}
			time.Sleep(r.config.PostponeInterval)
		}
		// Put item to the destination queue.
		if err = r.config.Queue.Enqueue(x); err != nil {
			r.config.MetricsWriter.Fail(r.config.Key, "enqueue fail")
			continue
		}
		r.config.MetricsWriter.Restore(r.config.Key, len(r.buf))
	}
	return nil
}

// Close gracefully stops the restorer.
func (r *Restorer) Close() error {
	return r.CloseWithTimeout(time.Second * 30)
}

// CloseWithTimeout stops the queue with timeout.
func (r *Restorer) CloseWithTimeout(timeout time.Duration) error {
	now := time.Now()
	for atomic.LoadUint32(&r.lock) == 1 {
		if time.Since(now) > timeout {
			return ErrTimeout
		}
	}
	r.setStatus(blqueue.StatusClose)
	return nil
}

// ForceClose immediately stops the queue.
func (r *Restorer) ForceClose() error {
	r.setStatus(blqueue.StatusClose)
	return nil
}

// Init the restorer.
func (r *Restorer) init() {
	c := r.config

	// Check mandatory params.
	if len(c.Key) == 0 {
		r.Err = blqueue.ErrNoKey
		r.setStatus(blqueue.StatusFail)
		return
	}
	if c.Decoder == nil {
		r.Err = ErrNoDecoder
		r.setStatus(blqueue.StatusFail)
		return
	}
	if c.Reader == nil {
		r.Err = ErrNoReader
		r.setStatus(blqueue.StatusFail)
		return
	}
	if c.Queue == nil {
		r.Err = ErrNoQueue
		r.setStatus(blqueue.StatusFail)
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

	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetrics{}
	}

	// Restorer is ready!
	r.setStatus(blqueue.StatusActive)

	// Init background check ticker.
	ticker := time.NewTicker(c.CheckInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if r.getStatus() == blqueue.StatusClose {
					ticker.Stop()
					return
				}
				_ = r.Restore()
			}
		}
	}()
}

func (r *Restorer) setStatus(status blqueue.Status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *Restorer) getStatus() blqueue.Status {
	return blqueue.Status(atomic.LoadUint32((*uint32)(&r.status)))
}
