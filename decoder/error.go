package decoder

import "errors"

var (
	ErrNoNewFunc                = errors.New("no New() function provided")
	ErrIncompatibleUnmarshaller = errors.New("incompatible unmarshaller")
)
