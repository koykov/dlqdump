package dlqdump

type Encoder interface {
	Encode(dst []byte, x interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(p []byte) (interface{}, error)
}
