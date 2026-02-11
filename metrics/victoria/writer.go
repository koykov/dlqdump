package victoria

import (
	"time"

	"github.com/koykov/vmchain"
)

type Writer interface {
	Dump(size int)
	Flush(reason string, size int)
	Restore(size int)
	Fail(reason string)
}

type writer struct {
	name string
	prec time.Duration
}

// NewWriter makes new instance of metrics writer.
func NewWriter(name string, options ...Option) Writer {
	mw := &writer{name: name}
	for _, fn := range options {
		fn(mw)
	}
	if mw.prec == 0 {
		mw.prec = time.Nanosecond
	}
	return mw
}

func (w *writer) Dump(size int) {
	vmchain.Counter("dlqdump_bytes_in").WithLabel("queue", w.name).Add(size)
	vmchain.Counter("dlqdump_size_in").WithLabel("queue", w.name).Inc()
}

func (w *writer) Flush(reason string, size int) {
	vmchain.Counter("dlqdump_bytes_flush").WithLabel("queue", w.name).WithLabel("reason", reason).Add(size)
}

func (w *writer) Restore(size int) {
	vmchain.Counter("dlqdump_bytes_out").WithLabel("queue", w.name).Add(size)
	vmchain.Counter("dlqdump_size_out").WithLabel("queue", w.name).Inc()
}

func (w *writer) Fail(reason string) {
	vmchain.Counter("dlqdump_fail").WithLabel("queue", w.name).WithLabel("reason", reason).Inc()
}

var _ = NewWriter
