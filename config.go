package dlqdump

import "time"

type Config struct {
	Version     uint32
	Key         string
	Size        MemorySize
	TimeLimit   time.Duration
	Encoder     Encoder
	Decoder     Decoder
	Destination string
}
