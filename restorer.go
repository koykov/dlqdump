package dlqdump

import (
	"sync"
	"sync/atomic"

	"github.com/koykov/blqueue"
)

// Restorer represents dump restore handler.
// Restorer may be scheduled (see Config.CheckInterval).
type Restorer struct {
	config *Config
	status blqueue.Status

	once sync.Once

	Err error
}

func NewRestorer(config *Config) (*Restorer, error) {
	r := &Restorer{
		config: config.Copy(),
	}
	r.once.Do(r.init)
	return r, nil
}

func (r *Restorer) init() {
	c := r.config

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

	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetrics{}
	}

	r.setStatus(blqueue.StatusActive)
}

func (r *Restorer) setStatus(status blqueue.Status) {
	atomic.StoreUint32((*uint32)(&r.status), uint32(status))
}

func (r *Restorer) getStatus() blqueue.Status {
	return blqueue.Status(atomic.LoadUint32((*uint32)(&r.status)))
}
