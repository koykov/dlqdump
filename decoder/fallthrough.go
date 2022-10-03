package decoder

type Fallthrough struct{}

func (r Fallthrough) Decode(p []byte) (interface{}, error) {
	return p, nil
}
