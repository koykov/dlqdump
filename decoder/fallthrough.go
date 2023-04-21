package decoder

// Fallthrough represents decoder that returns source data as result.
type Fallthrough struct{}

func (d Fallthrough) Decode(p []byte) (any, error) {
	return p, nil
}
