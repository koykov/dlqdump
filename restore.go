package dlqdump

import (
	"context"
	"sync/atomic"
	"time"
)

func (q *Queue) restore(ctx context.Context) {
	ticker := time.NewTicker(q.config.RestoreDisallowDelay)
	for {
		select {
		case <-ticker.C:
			q._restore()
		case <-ctx.Done():
			for atomic.LoadUint32(&q.rlck) == 1 {
			}
			ticker.Stop()
			return
		}
	}
}

func (q *Queue) _restore() {
	atomic.StoreUint32(&q.rlck, 1)
	defer atomic.StoreUint32(&q.rlck, 0)
	// todo implement restore
}
