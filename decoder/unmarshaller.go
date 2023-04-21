package decoder

// UnmarshallerInterface is the interface that wraps the basic Unmarshal method.
type UnmarshallerInterface interface {
	Unmarshal([]byte) error
}

// Unmarshaller provides New() function that returns the object implements UnmarshallerInterface interface.
type Unmarshaller struct {
	// New is a function that must return object implements UnmarshallerInterface interface.
	// It may not be changed.
	New func() UnmarshallerInterface
}

func (d Unmarshaller) Decode(p []byte) (any, error) {
	if d.New == nil {
		return nil, ErrNoNewFunc
	}
	m := d.New()
	if m == nil {
		return nil, ErrEmptyUnmarshaller
	}
	err := m.Unmarshal(p)
	return m, err
}
