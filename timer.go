package dlqdump

import (
	"sync/atomic"
	"time"
)

const (
	timerReach timerSignal = iota
	timerReset
	timerStop
)

type timerSignal uint8

type timer struct {
	c chan timerSignal
	s uint32
}

func newTimer() *timer {
	t := timer{c: make(chan timerSignal, 1)}
	return &t
}

func (t *timer) wait(queue *Queue) {
	time.AfterFunc(queue.config.FlushInterval, func() {
		queue.timer.reach()
	})
	for {
		signal, ok := <-t.c
		if !ok {
			return
		}
		switch signal {
		case timerReach:
			_ = queue.flush(flushReasonTimeLimit)
		case timerReset:
			break
		case timerStop:
			atomic.StoreUint32(&t.s, 1)
			close(t.c)
			return
		}
	}
}

func (t *timer) reach() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReach
}

func (t *timer) reset() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReset
}

func (t *timer) stop() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerStop
}
