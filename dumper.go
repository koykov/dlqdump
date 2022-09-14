package dlqdump

type Dumper interface {
	Dump([]byte) (int, error)
	Size() MemorySize
	Flush() error
}
