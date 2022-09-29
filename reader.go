package dlqdump

type Reader interface {
	Read([]byte) (Version, []byte, error)
}
