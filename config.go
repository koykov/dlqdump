package dlqdump

import (
	"time"

	"github.com/koykov/queue"
)

const (
	// Default time limit to flush the data.
	defaultFlushInterval = time.Second * 30
	// Default rate limit that allows restore.
	defaultAllowRate = .95
	// Default delay if restore allow rate limit exceeds.
	defaultCheckInterval = time.Second
)

type Config struct {
	/*
		Common params.
	*/
	// Write version. Must be changed at any change of Encoder param.
	Version Version
	// Unique queue key. Indicates queue in logs and metrics.
	// Mandatory param.
	Key string
	// Metrics writer handler.
	MetricsWriter MetricsWriter
	// Logger handler.
	Logger queue.Logger

	/*
		Queue params.
		Params of this group will ignore by Restorer.
	*/
	// Max queue capacity in bytes.
	// When dumped data will reach size, queue will flush the data.
	Capacity MemorySize
	// Wait duration until flush the data.
	// After first incoming item will start the timer to flush the data when timer reach.
	// If this param omit defaultFlushInterval (30 seconds) will use by default.
	FlushInterval time.Duration
	// Encoder helper to convert item to bytes.
	// Will use universal encoder if omitted.
	Encoder Encoder
	// Writer helper to dump data to various destinations.
	// Mandatory param.
	Writer Writer

	/*
		Restorer params.
		Params of this group will ignore by Queue.
	*/
	// Interval between restore attempts.
	// If this param omit defaultWaitInterval (1 second) will use by default.
	CheckInterval time.Duration
	// Wait duration if queue rate exceeds AllowRate.
	// If this param omit CheckInterval will use instead.
	PostponeInterval time.Duration
	// Queue rate that allows restore.
	// If this param omit defaultAllowRate (95%) will use by default.
	AllowRate float32
	// Helper to achieve data from dump.
	// Mandatory param.
	Reader Reader
	// Decoder helper to convert bytes to item.
	// Mandatory param.
	Decoder Decoder
	// Destination queue to restore dump.
	// Mandatory param.
	Queue queue.Interface
}

// Copy copies config instance to protect queue from changing params after start.
// It means that after starting queue all config modifications will have no effect.
func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
