package encoder

import (
	"github.com/koykov/bytealg"
)

// marshallerTo is the interface that wraps the basic MarshalTo method.
type marshallerTo interface {
	Size() int
	MarshalTo([]byte) (int, error)
}

// MarshallerTo is an encoder that can encode objects implementing marshallerTo interface.
type MarshallerTo struct{}

func (e MarshallerTo) Encode(dst []byte, x interface{}) ([]byte, error) {
	if m, ok := x.(marshallerTo); ok {
		off := len(dst)
		dst = bytealg.GrowDelta(dst, m.Size())
		_, err := m.MarshalTo(dst[off:])
		return dst, err
	}
	return dst, ErrIncompatibleEncoder
}
