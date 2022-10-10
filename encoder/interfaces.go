package encoder

// Byter is the interface that wraps the basic Bytes method.
type Byter interface {
	Bytes() []byte
}

// Marshaller is the interface that wraps the basic Marshal method.
type Marshaller interface {
	Marshal() ([]byte, error)
}

// MarshallerTo is the interface that wraps the basic MarshalTo method.
type MarshallerTo interface {
	Size() int
	MarshalTo([]byte) (int, error)
}
