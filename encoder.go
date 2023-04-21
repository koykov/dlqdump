package dlqdump

// Encoder is the interface that wraps the basic Encode method.
//
// Encode encodes x to dst. It returns dst and any error encountered.
type Encoder interface {
	Encode(dst []byte, x any) ([]byte, error)
}
