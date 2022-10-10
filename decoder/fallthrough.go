package decoder

// Fallthrough represents decoder that returns source data as result.
type Fallthrough struct{}

func (r Fallthrough) Decode(p []byte) (interface{}, error) {
	return p, nil
}
