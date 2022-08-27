package dlqdump

type Byter interface {
	Bytes() []byte
}

type Marshaller interface {
	Marshal() ([]byte, error)
}

type MarshallerTo interface {
	Size() int
	MarshalTo([]byte) (int, error)
}
