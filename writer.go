package dlqdump

type Writer interface {
	Write(Version, []byte) (int, error)
	Size() MemorySize
	Flush() error
}
