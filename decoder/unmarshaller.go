package decoder

// unmarshaller is the interface that wraps the basic Unmarshal method.
type unmarshaller interface {
	Unmarshal([]byte) error
}

// Unmarshaller provides New() function that returns the object implements unmarshaller interface.
type Unmarshaller struct {
	// New is a function that must return object implements unmarshaller interface.
	// It may not be changed.
	New func() interface{}
}

func (d Unmarshaller) Decode(p []byte) (interface{}, error) {
	if d.New == nil {
		return nil, ErrNoNewFunc
	}
	m, ok := d.New().(unmarshaller)
	if !ok {
		return nil, ErrIncompatibleUnmarshaller
	}
	err := m.Unmarshal(p)
	return m, err
}
