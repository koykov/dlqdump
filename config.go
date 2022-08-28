package dlqdump

import (
	"time"

	"github.com/koykov/blqueue"
)

const (
	defaultFileMask  = "%Y-%m-%d--%H-%M-%S--%i.bin"
	defaultTimeLimit = time.Second * 30
)

type Config struct {
	Version       uint32
	Key           string
	Size          MemorySize
	TimeLimit     time.Duration
	Encoder       Encoder
	Decoder       Decoder
	Directory     string
	FileMask      string
	MetricsWriter MetricsWriter
	Logger        blqueue.Logger
}

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
