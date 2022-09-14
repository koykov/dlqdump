package dlqdump

import "errors"

var (
	ErrIncompatibleMarshaller = errors.New("incompatible encoder")
	ErrNoDestinationDir       = errors.New("no destination directory provided")
	ErrNoEncoder              = errors.New("no encoder provided")
	ErrNoDecoder              = errors.New("no decoder provided")
)
