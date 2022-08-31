package dlqdump

import (
	"context"
	"sync/atomic"
	"time"
)

type rworker struct {
	config *Config
	lck    uint32
	ctx    context.Context
	cancel context.CancelFunc
}

func newRestoreWorker(config *Config) *rworker {
	w := rworker{config: config}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	return &w
}

func (w *rworker) observe() {
	ticker := time.NewTicker(w.config.RestoreDisallowDelay)
	for {
		select {
		case <-ticker.C:
			// ...
		case <-w.ctx.Done():
			w.wait()
			ticker.Stop()
			return
		}
	}
}

func (w *rworker) do() {
	w.lock()
	defer w.unlock()
	// todo implement me
}

func (w *rworker) lock() {
	atomic.StoreUint32(&w.lck, 1)
}

func (w *rworker) unlock() {
	atomic.StoreUint32(&w.lck, 0)
}

func (w *rworker) wait() {
	for atomic.LoadUint32(&w.lck) == 1 {
	}
}

func (w *rworker) stop() {
	w.cancel()
	w.wait()
}
