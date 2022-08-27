package dlqdump

import "time"

const (
	defaultFileMask  = "%Y-%m-%d--%H-%M-%S--%i.bin"
	defaultTimeLimit = time.Second * 30
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

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
