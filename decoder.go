package dlqdump

// Decoder is the interface that wraps the basic Write Decode.
//
// Decode decodes value from p. It returns decoded value and any error encountered.
type Decoder interface {
	Decode(p []byte) (any, error)
}
