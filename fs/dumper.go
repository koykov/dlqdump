package fs

import (
	"github.com/koykov/dlqdump"
)

type Dumper struct {
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

func (d Dumper) Dump(p []byte) (n int, err error) {
	_ = p
	return
}

func (d Dumper) Size() dlqdump.MemorySize {
	return 0
}

func (d Dumper) Flush() error {
	return nil
}
