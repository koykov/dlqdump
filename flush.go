package dlqdump

import (
	"os"
	"time"

	"github.com/koykov/clock"
	"github.com/koykov/fastconv"
)

const (
	flushReasonSize flushReason = iota
	flushReasonTimeLimit
	flushReasonForce

	writeBufferSize = 16
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

	size := len(q.buf)
	if q.buf, err = clock.AppendFormat(q.buf, q.config.FileMask, time.Now()); err != nil {
		return
	}
	filename := q.buf[size:]
	off1 := len(q.buf)
	q.buf = append(q.buf, q.config.Directory...)
	q.buf = append(q.buf, os.PathSeparator)
	q.buf = append(q.buf, filename...)
	filepath := fastconv.B2S(q.buf[off1:])
	q.buf = append(q.buf, ".tmp"...)
	filepathTmp := fastconv.B2S(q.buf[off1:])

	var f *os.File
	if f, err = os.Create(filepathTmp); err != nil {
		return
	}
	p := q.buf[:size]
	for len(p) >= writeBufferSize {
		if _, err = f.Write(p[:writeBufferSize]); err != nil {
			return
		}
		p = p[writeBufferSize:]
	}
	if len(p) > 0 {
		if _, err = f.Write(p); err != nil {
			return
		}
	}
	q.buf = q.buf[:0]
	if err = f.Close(); err != nil {
		return
	}

	err = os.Rename(filepathTmp, filepath)

	q.config.MetricsWriter.QueueFlush(q.config.Key, reason.String(), size)

	return
}

func (r flushReason) String() string {
	switch r {
	case flushReasonSize:
		return "size limit"
	case flushReasonTimeLimit:
		return "time limit"
	case flushReasonForce:
		return "force"
	default:
		return "unknown"
	}
}
