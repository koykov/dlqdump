package encoder

import (
	"fmt"

	"github.com/koykov/bytealg"
	"github.com/koykov/dlqdump"
)

type Basic struct{}

func (e Basic) Encode(dst []byte, x interface{}) ([]byte, error) {
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
	case MarshallerTo:
		off := len(dst)
		m := x.(MarshallerTo)
		dst = bytealg.GrowDelta(dst, m.Size())
		_, err = m.MarshalTo(dst[off:])
	case Marshaller:
		m := x.(Marshaller)
		var b []byte
		if b, err = m.Marshal(); err == nil {
			dst = append(dst, b...)
		}
	case Byter:
		dst = append(dst, x.(Byter).Bytes()...)
	case fmt.Stringer:
		dst = append(dst, x.(fmt.Stringer).String()...)
	default:
		err = dlqdump.ErrIncompatibleEncoder
	}
	return dst, err
}
