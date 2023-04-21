package encoder

import (
	"fmt"
)

// byter is the interface that wraps the basic Bytes method.
type byter interface {
	Bytes() []byte
}

// Builtin represents encoder for basic cases:
// * string
// * byte slice
// * byter
// * Stringer
type Builtin struct{}

func (e Builtin) Encode(dst []byte, x any) ([]byte, error) {
	var err error
	switch x.(type) {
	case []byte:
		dst = append(dst, x.([]byte)...)
	case *[]byte:
		dst = append(dst, *x.(*[]byte)...)
	case string:
		dst = append(dst, x.(string)...)
	case *string:
		dst = append(dst, *x.(*string)...)
	case byter:
		dst = append(dst, x.(byter).Bytes()...)
	case fmt.Stringer:
		dst = append(dst, x.(fmt.Stringer).String()...)
	default:
		err = ErrIncompatibleEncoder
	}
	return dst, err
}
