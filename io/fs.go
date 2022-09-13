package io

import (
	"github.com/koykov/blqueue"
	"github.com/koykov/dlqdump"
)

type FS struct {
	// Max buffer size in bytes.
	// Dumper will move buffered data to destination file on overflow.
	Buffer dlqdump.MemorySize
	// Destination directory for dump files.
	Directory string
	// Dump file mask.
	// Supports strftime patterns (see https://github.com/koykov/clock#format).
	// If this param omit defaultFileMask ("%Y-%m-%d--%H-%M-%S--%i.bin") will use instead.
	FileMask string
}

func (d FS) Dump(p []byte) (n int, err error) {
	_ = p
	return
}

func (d FS) Flush() error {
	return nil
}

func (d FS) Restore(dst blqueue.Interface) error {
	_ = dst
	return nil
}
