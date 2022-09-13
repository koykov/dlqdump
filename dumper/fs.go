package dumper

type FS struct{}

func (d FS) Dump(p []byte) (n int, err error) {
	_ = p
	return
}

func (d FS) Flush() error {
	return nil
}
