package fs

import "github.com/koykov/dlqdump"

type Reader struct {
	// MatchMask represents pattern to match the names of dump files.
	// Mandatory param.
	MatchMask string
}

func (r Reader) Read(dst []byte) (dlqdump.Version, []byte, error) {
	return 0, dst, nil
}
