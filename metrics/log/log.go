package dlqdump

import (
	"log"
)

// LogMetrics is Log implementation of dlqdump.MetricsWriter.
//
// Don't use in production. Only for debug purposes.
type LogMetrics struct {
	name string
}

var _ = NewLogMetrics

func NewLogMetrics(name string) *LogMetrics {
	m := &LogMetrics{name}
	return m
}

func (m LogMetrics) Dump(size int) {
	log.Printf("queue %s: %d bytes come to the queue\n", m.name, size)
}

func (m LogMetrics) Flush(reason string, size int) {
	log.Printf("queue %s: flush %d bytes due to reason %s\n", m.name, size, reason)
}

func (m LogMetrics) Restore(size int) {
	log.Printf("queue %s: %d bytes restored from dump\n", m.name, size)
}

func (m LogMetrics) Fail(reason string) {
	log.Printf("queue %s: restore failed with reason '%s'\n", m.name, reason)
}
