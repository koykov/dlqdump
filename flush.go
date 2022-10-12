package dlqdump

type flushReason uint8

const (
	flushReasonSize flushReason = iota
	flushReasonInterval
	flushReasonForce
)

// Flush data collected in the queue.
func (q *Queue) flush(reason flushReason) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	if reason == flushReasonForce {
		q.timer.reset()
	}
	return q.flushLF(reason)
}

// Lock-free version of flush method.
func (q *Queue) flushLF(reason flushReason) (err error) {
	if l := q.config.Logger; l != nil {
		msg := "flush by reason '%s'"
		l.Printf(msg, reason)
	}

	size := q.config.Writer.Size()
	if size > 0 {
		if err = q.config.Writer.Flush(); err != nil {
			return
		}
	}

	q.config.MetricsWriter.Flush(reason.String(), int(size))

	return
}

func (r flushReason) String() string {
	switch r {
	case flushReasonSize:
		return "size limit"
	case flushReasonInterval:
		return "reach interval"
	case flushReasonForce:
		return "force"
	default:
		return "unknown"
	}
}
