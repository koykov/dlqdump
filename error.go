package dlqdump

import "errors"

var (
	ErrIncompatibleEncoder = errors.New("incompatible encoder")
	ErrNoDestinationDir    = errors.New("no destination directory provided")
	ErrNoEncoder           = errors.New("no encoder provided")
	ErrNoDecoder           = errors.New("no decoder provided")
	ErrNoWriter            = errors.New("no writer provided")
	ErrNoReader            = errors.New("no reader provided")
	ErrNoQueue             = errors.New("no destination queue provided")
	ErrTimeout             = errors.New("operation too long")
)
