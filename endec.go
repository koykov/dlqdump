package dlqdump

// Encoder is the interface that wraps the basic Encode method.
//
// Encode encodes x to dst. It returns dst and any error encountered.
type Encoder interface {
	Encode(dst []byte, x interface{}) ([]byte, error)
}

// Decoder is the interface that wraps the basic Write Decode.
//
// Decode decodes value from p. It returns decoded value and any error encountered.
type Decoder interface {
	Decode(p []byte) (interface{}, error)
}
