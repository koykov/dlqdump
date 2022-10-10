package dlqdump

// Writer is the interface that wraps the basic Write method.
type Writer interface {
	// Write returns the number of bytes written from p (0 <= n <= len(p) + 8) and any error encountered.
	Write(ver Version, p []byte) (int, error)
	// Size returns size in bytes of collected data.
	Size() MemorySize
	// Flush flushes collected data.
	Flush() error
}
