package dlqdump

import (
	"time"

	"github.com/koykov/blqueue"
)

const (
	// Dump file names default mask.
	defaultFileMask = "%Y-%m-%d--%H-%M-%S--%i.bin"
	// Default time limit to flush the data.
	defaultTimeLimit = time.Second * 30
)

type Config struct {
	// Dump version. Must be changed at any change of Encoder and/or Decoder params.
	Version uint32
	// Unique queue key. Indicates queue in logs and metrics.
	// Mandatory param.
	Key string
	// Max queue capacity in bytes.
	// When dumped data will reach size, queue will flush the data.
	Size MemorySize
	// Wait duration until flush the data.
	// After first incoming item will start the timer to flush the data when timer reach.
	// If this param omit defaultTimeLimit (30 seconds) will use instead.
	TimeLimit time.Duration
	// Encoder helper to convert item to bytes.
	// Will use universal encoder if omitted.
	Encoder Encoder
	// Decoder helper to convert bytes to item.
	// Mandatory if RestoreTo param specified.
	Decoder Decoder
	// Destination directory for dump files.
	Directory string
	// Dump file mask.
	// Supports strftime patterns (see https://github.com/koykov/clock#format).
	// If this param omit defaultFileMask ("%Y-%m-%d--%H-%M-%S--%i.bin") will use instead.
	FileMask string
	// Metrics writer handler.
	MetricsWriter MetricsWriter
	// Logger handler.
	Logger blqueue.Logger
}

// Copy copies config instance to protect queue from changing params after start.
// It means that after starting queue all config modifications will have no effect.
func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
