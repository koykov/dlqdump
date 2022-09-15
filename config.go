package dlqdump

import (
	"time"

	"github.com/koykov/blqueue"
)

const (
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
	Capacity MemorySize
	// Wait duration until flush the data.
	// After first incoming item will start the timer to flush the data when timer reach.
	// If this param omit defaultTimeLimit (30 seconds) will use instead.
	FlushInterval time.Duration

	// Encoder helper to convert item to bytes.
	// Will use universal encoder if omitted.
	Encoder Encoder

	// Dumper helper to dump data to various destinations.
	// Mandatory param.
	Dumper Dumper

	// Destination directory for dump files.
	Directory string
	// Dump file mask.
	// Supports strftime patterns (see https://github.com/koykov/clock#format).
	// If this param omit defaultFileMask ("%Y-%m-%d--%H-%M-%S--%i.bin") will use instead.
	FileMask string

	// RestoreTo indicates the queue to put data from dump files.
	// This param requires Decoder.
	RestoreTo blqueue.Interface
	// Queue rate that forbids or allows put data from dump files to the RestoreTo queue.
	// If this param omit defaultRestoreAllowRateLimit (95%) will use instead.
	RestoreAllowRateLimit float32
	// RestoreDisallowDelay indicates how many need wait before new attempt if RestoreAllowRateLimit was exceeded.
	// If this param omit defaultRestoreDisallowDelay (1 second) will use instead.
	RestoreDisallowDelay time.Duration

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
