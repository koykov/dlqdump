package log

import (
	"log"
)

// MetricsWriter is Log implementation of dlqdump.MetricsWriter.
//
// Don't use in production. Only for debug purposes.
type MetricsWriter struct {
	name string
}

var _ = NewLogMetrics

func NewLogMetrics(name string) *MetricsWriter {
	m := &MetricsWriter{name}
	return m
}

func (w MetricsWriter) Dump(size int) {
	log.Printf("queue %s: %d bytes come to the queue\n", w.name, size)
}

func (w MetricsWriter) Flush(reason string, size int) {
	log.Printf("queue %s: flush %d bytes due to reason %s\n", w.name, size, reason)
}

func (w MetricsWriter) Restore(size int) {
	log.Printf("queue %s: %d bytes restored from dump\n", w.name, size)
}

func (w MetricsWriter) Fail(reason string) {
	log.Printf("queue %s: restore failed with reason '%s'\n", w.name, reason)
}
