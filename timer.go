package dlqdump

import (
	"sync/atomic"
	"time"
)

type timerSignal uint8

const (
	timerReach timerSignal = iota
	timerReset
	timerStop
)

// Internal timer implementation.
type timer struct {
	c chan timerSignal
	s uint32
}

func newTimer() *timer {
	t := timer{c: make(chan timerSignal, 1)}
	return &t
}

// Background waiter method.
func (t *timer) wait(queue *Queue) {
	time.AfterFunc(queue.config.FlushInterval, func() {
		queue.timer.reach()
		queue.SetBit(flagTimer, false)
	})
	for {
		signal, ok := <-t.c
		if !ok {
			return
		}
		switch signal {
		case timerReach:
			_ = queue.flush(flushReasonInterval)
		case timerReset:
			break
		case timerStop:
			atomic.StoreUint32(&t.s, 1)
			close(t.c)
			return
		}
	}
}

// Send time reach signal.
func (t *timer) reach() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReach
}

// Send reset signal.
func (t *timer) reset() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReset
}

// Send stop signal.
func (t *timer) stop() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerStop
}
