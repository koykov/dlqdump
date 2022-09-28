package dlqdump

type MetricsWriter interface {
	DumpPut(queue string, size int)
	DumpFlush(queue, reason string, size int)
	DumpRestore(queue string, size int)
}
