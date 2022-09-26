package fs

import (
	"github.com/koykov/blqueue"
)

type Restorer struct {
	// MatchMask represents pattern to match the names of dump files.
	// Mandatory param.
	MatchMask string
}

func (r Restorer) Restore(dst blqueue.Interface) error {
	_ = dst
	return nil
}
