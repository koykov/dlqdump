package dlqdump

type MetricsWriter interface {
	QueuePut(queue string, size int)
	QueueFlush(queue, reason string, size int)
}
