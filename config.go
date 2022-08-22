package dlqdump

import "time"

const (
	defaultFileMask = "{key}__%Y-%m-%d"
)

type Config struct {
	Version   uint32
	Key       string
	Size      MemorySize
	TimeLimit time.Duration
	Encoder   Encoder
	Decoder   Decoder
	Directory string
	FileMask  string
}
