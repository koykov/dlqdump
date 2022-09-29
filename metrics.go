package dlqdump

type MetricsWriter interface {
	Dump(queue string, size int)
	Flush(queue, reason string, size int)
	Restore(queue string, size int)
}
