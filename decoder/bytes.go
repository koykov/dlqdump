package decoder

type Bytes struct{}

func (Bytes) Decode(p []byte) (any, error) {
	if len(p) == 0 {
		return []byte{}, nil
	}
	buf := make([]byte, len(p))
	copy(buf, p)
	return buf, nil
}
