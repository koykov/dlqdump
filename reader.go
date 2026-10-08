package dlqdump

import "io"

// Reader is the interface that wraps the basic Read method.
//
// Read reads next encoded entry from the dump to dst. It returns version, dst contains entry and any error encountered.
type Reader interface {
	io.Closer
	Read(dst []byte) (Version, []byte, error)
}
