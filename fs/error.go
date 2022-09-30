package fs

import "errors"

var (
	ErrNoDestinationDir = errors.New("no destination directory provided")
	ErrDirNoWR          = errors.New("directory doesn't exists or writable")
)
