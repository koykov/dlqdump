package dlqdump

import "errors"

var (
	ErrUnknownMarshaller = errors.New("unknown marshaller")
	ErrNoDestinationDir  = errors.New("no destination directory provided")
)
