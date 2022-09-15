package dlqdump

import "errors"

var (
	ErrIncompatibleEncoder = errors.New("incompatible encoder")
	ErrNoDestinationDir    = errors.New("no destination directory provided")
	ErrNoEncoder           = errors.New("no encoder provided")
	ErrNoDumper            = errors.New("no dumper provided")
)
