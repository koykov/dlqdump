package encoder

// marshaller is the interface that wraps the basic Marshal method.
type marshaller interface {
	Marshal() ([]byte, error)
}

// Marshaller is an encoder that can encode objects implementing marshaller interface.
type Marshaller struct{}

func (e Marshaller) Encode(dst []byte, x interface{}) ([]byte, error) {
	if m, ok := x.(marshaller); ok {
		b, err := m.Marshal()
		if err != nil {
			return dst, err
		}
		dst = append(dst, b...)
		return dst, nil
	}
	return dst, ErrIncompatibleEncoder
}
