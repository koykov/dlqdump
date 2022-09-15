package fs

import (
	"time"

	"github.com/koykov/blqueue"
	"github.com/koykov/dlqdump"
)

const (
	// Default rate limit that allows restore.
	defaultRestoreAllowRateLimit = .95
	// Default delay if restore allow rate limit exceeds.
	defaultRestoreDisallowDelay = time.Second
)

type Restorer struct {
	// Decoder helper to convert bytes to item.
	// Mandatory param.
	Decoder dlqdump.Decoder
	// Destination indicates the queue to put data from dump files.
	// Mandatory param.
	Destination blqueue.Interface
	// Queue rate that forbids or allows put data from dump the Destination queue.
	// If this param omit defaultRestoreAllowRateLimit (95%) will use instead.
	AllowRateLimit float32
	// DisallowDelay indicates how many need wait before new attempt if RestoreAllowRateLimit was exceeded.
	// If this param omit defaultRestoreDisallowDelay (1 second) will use instead.
	DisallowDelay time.Duration
}

func (r Restorer) Restore(dst blqueue.Interface) error {
	_ = dst
	return nil
}
