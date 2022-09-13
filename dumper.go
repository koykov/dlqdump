package dlqdump

type Dumper interface {
	Dump([]byte) (int, error)
	Flush() error
}
