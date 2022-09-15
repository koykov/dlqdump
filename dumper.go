package dlqdump

type Dumper interface {
	Dump(uint32, []byte) (int, error)
	Size() MemorySize
	Flush() error
}
