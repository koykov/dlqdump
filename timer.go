package dlqdump

import (
	"sync/atomic"
	"time"
)

// Internal timer implementation.
type timer struct {
	at *time.Timer
}

func newTimer() *timer {
	return &timer{}
}

// Start the timer: schedule a flush after Config.FlushInterval.
// Must be called with Queue.mux held.
func (t *timer) wait(queue *Queue) {
	t.at = time.AfterFunc(queue.config.FlushInterval, func() {
		// End the current period so the next item arms a fresh timer.
		atomic.StoreUint32(&queue.timerOn, 0)
		_ = queue.flush(flushReasonInterval)
	})
}

// stop cancels a pending flush.
// Must be called with Queue.mux held.
func (t *timer) stop() {
	if t.at != nil {
		t.at.Stop()
		t.at = nil
	}
}
