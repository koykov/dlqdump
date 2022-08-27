package dlqdump

type Marshaller interface {
	Marshal() ([]byte, error)
}

type MarshallerTo interface {
	Size() int
	MarshalTo([]byte) (int, error)
}
