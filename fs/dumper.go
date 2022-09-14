package fs

import (
	"sync"
	"sync/atomic"

	"github.com/koykov/dlqdump"
)

const defaultFileMask = "%Y-%m-%d--%H-%M-%S--%i.bin"

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

	once sync.Once
	bs   dlqdump.MemorySize
	dir  string
	mask string
	sz   uint64

	err error
}

func (d *Dumper) Dump(p []byte) (n int, err error) {
	d.once.Do(d.init)
	if d.err != nil {
		return 0, d.err
	}

	_ = p
	// todo implement me

	return
}

func (d *Dumper) Size() dlqdump.MemorySize {
	return dlqdump.MemorySize(atomic.LoadUint64(&d.sz))
}

func (d *Dumper) Flush() error {
	d.once.Do(d.init)
	if d.err != nil {
		return d.err
	}
	return nil
}

func (d *Dumper) init() {
	d.err = nil
	if len(d.Directory) == 0 {
		d.err = dlqdump.ErrNoDestinationDir
		return
	}
	if !isDirWR(d.Directory) {
		d.err = ErrDirNoWR
		return
	}
	if len(d.FileMask) == 0 {
		d.FileMask = defaultFileMask
	}

	d.bs = d.Buffer
	d.dir = d.Directory
	d.mask = d.FileMask
	d.sz = 0
}
