package dlqdump

import (
	"time"

	"github.com/koykov/blqueue"
)

const (
	// Default rate limit that allows restore.
	defaultAllowRate = .95
	// Default delay if restore allow rate limit exceeds.
	defaultWaitInterval = time.Second
)

// DRConfig represents DRC/sDRC config.
type DRConfig struct {
	// Destination queue to restore dump.
	// Mandatory param.
	Queue blqueue.Interface
	// Decoder helper to convert bytes to item.
	// Mandatory param.
	Decoder Decoder
	// Helper to achieve data from dump.
	// Mandatory param.
	Restorer Restorer
	// Duration between check fresh data in dump.
	// This param enables scheduled feature.
	CheckInterval time.Duration
	// WaitInterval indicates how many need wait before new attempt if AllowRate was exceeded.
	// If this param omit defaultWaitInterval (1 second) will use instead.
	WaitInterval time.Duration
	// Queue rate that forbids or allows put data from dump the Destination queue.
	// If this param omit defaultAllowRate (95%) will use instead.
	AllowRate float32
}

func (c *DRConfig) Copy() *DRConfig {
	cpy := *c
	return &cpy
}
