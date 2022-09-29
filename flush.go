package dlqdump

const (
	flushReasonSize flushReason = iota
	flushReasonInterval
	flushReasonForce
)

type flushReason uint8

func (q *Queue) flush(reason flushReason) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	if reason == flushReasonForce {
		q.timer.reset()
	}
	return q.flushLF(reason)
}

func (q *Queue) flushLF(reason flushReason) (err error) {
	if l := q.config.Logger; l != nil {
		msg := "queue #%s flush by reason '%s'"
		l.Printf(msg, q.config.Key, reason)
	}

	size := q.config.Writer.Size()
	if size > 0 {
		if err = q.config.Writer.Flush(); err != nil {
			return
		}
	}

	q.config.MetricsWriter.Flush(q.config.Key, reason.String(), int(size))

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
