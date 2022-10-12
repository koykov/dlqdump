package dlqdump

// MetricsWriter is the interface that wraps the basic metrics methods.
type MetricsWriter interface {
	// Dump registers how many bytes dumped to the queue.
	Dump(size int)
	// Flush registers how many bytes flushed to the queue and what reason is.
	Flush(reason string, size int)
	// Restore registers how many bytes restored to the queue.
	Restore(size int)
	// Fail registers fail reason for given queue.
	Fail(reason string)
}
