package decoder

import "errors"

var (
	ErrNoNewFunc         = errors.New("no New() function provided")
	ErrEmptyUnmarshaller = errors.New("empty unmarshaller")
)
