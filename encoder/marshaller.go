package encoder

import "github.com/koykov/bytealg"

// MarshallerInterface is the interface that wraps the basic Marshal method.
type MarshallerInterface interface {
	Marshal() ([]byte, error)
}

// MarshallerToInterface is the interface that wraps the basic MarshalTo method.
type MarshallerToInterface interface {
	Size() int
	MarshalTo([]byte) (int, error)
}

// Marshaller is an encoder that can encode objects implementing MarshallerInterface or MarshallerToInterface interface.
type Marshaller struct{}

func (e Marshaller) Encode(dst []byte, x interface{}) ([]byte, error) {
	if m, ok := x.(MarshallerToInterface); ok {
		off := len(dst)
		dst = bytealg.GrowDelta(dst, m.Size())
		_, err := m.MarshalTo(dst[off:])
		return dst, err
	}
	if m, ok := x.(MarshallerInterface); ok {
		b, err := m.Marshal()
		if err != nil {
			return dst, err
		}
		dst = append(dst, b...)
		return dst, nil
	}
	return dst, ErrIncompatibleEncoder
}
