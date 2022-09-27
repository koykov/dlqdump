package fs

type Reader struct {
	// MatchMask represents pattern to match the names of dump files.
	// Mandatory param.
	MatchMask string
}

func (r Reader) Read(buf []byte) (int, error) {
	_ = buf
	return 0, nil
}
