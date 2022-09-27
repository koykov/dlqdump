package dlqdump

import (
	"time"

	"github.com/koykov/blqueue"
)

const (
	// Default time limit to flush the data.
	defaultTimeLimit = time.Second * 30
	// Default rate limit that allows restore.
	defaultAllowRate = .95
	// Default delay if restore allow rate limit exceeds.
	defaultWaitInterval = time.Second
)

type Config struct {
	// Write version. Must be changed at any change of Encoder param.
	Version Version
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
	// Writer helper to dump data to various destinations.
	// Mandatory param.
	Writer Writer

	// Helper to achieve data from dump.
	// Mandatory param.
	Restorer Restorer
	// Decoder helper to convert bytes to item.
	// Mandatory param.
	Decoder Decoder
	// Destination queue to restore dump.
	// Mandatory param.
	Queue blqueue.Interface

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
