package fs

import "errors"

var (
	ErrDirNoWR = errors.New("directory doesn't exists or writable")
)
