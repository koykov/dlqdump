package dlqdump

import "io"

type Reader interface {
	io.Reader
	GetVersion() Version
}
