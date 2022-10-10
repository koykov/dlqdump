package dlqdump

// Reader is the interface that wraps the basic Read method.
//
// Read reads next encoded entry from the dump to dst. It returns version, dst contains entry and any error encountered.
type Reader interface {
	Read(dst []byte) (Version, []byte, error)
}
