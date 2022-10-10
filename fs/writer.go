package fs

import (
	"encoding/binary"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/bytealg"
	"github.com/koykov/clock"
	"github.com/koykov/dlqdump"
	"github.com/koykov/fastconv"
)

const (
	defaultFileMask = "%Y-%m-%d--%H-%M-%S--%i.bin"
	flushChunkSize  = 16
)

// Writer is file system implementation of dlqdump.Writer interface.
type Writer struct {
	// Max buffer size in bytes.
	// Writer will move buffered data to destination file on overflow.
	Buffer dlqdump.MemorySize
	// Destination directory for dump files.
	Directory string
	// Write file mask.
	// Supports strftime patterns (see https://github.com/koykov/clock#format).
	// If this param omit defaultFileMask ("%Y-%m-%d--%H-%M-%S--%i.bin") will use by default.
	FileMask string

	once sync.Once
	bs   dlqdump.MemorySize
	dir  string
	mask string
	sz   uint64

	mux sync.Mutex
	f   *os.File
	ft  string
	fd  string
	buf []byte

	err error
}

func (d *Writer) Write(version dlqdump.Version, p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	d.once.Do(d.init)
	if d.err != nil {
		return 0, d.err
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	off := len(d.buf)

	if off == 0 {
		// First write, encode version to header.
		d.buf = bytealg.GrowDelta(d.buf, 8)
		binary.LittleEndian.PutUint64(d.buf[off:], uint64(version))
		off += 8
	}

	// Write item length.
	d.buf = bytealg.GrowDelta(d.buf, 4)
	binary.LittleEndian.PutUint32(d.buf[off:], uint32(len(p)))
	// Write item body.
	d.buf = append(d.buf, p...)
	n = len(d.buf) - off
	atomic.AddUint64(&d.sz, uint64(n))

	if dlqdump.MemorySize(len(d.buf)) >= d.bs {
		// Buffer size reached limit, so flush it to the disk.
		if err = d.flushBuf(); err != nil {
			return
		}
	}

	return
}

func (d *Writer) Size() dlqdump.MemorySize {
	return dlqdump.MemorySize(atomic.LoadUint64(&d.sz))
}

func (d *Writer) Flush() (err error) {
	d.once.Do(d.init)
	if d.err != nil {
		return d.err
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	// Flush buffered data and clear buffer.
	if len(d.buf) > 0 {
		if err = d.flushBuf(); err != nil {
			return
		}
	}
	d.buf = d.buf[:0]
	atomic.StoreUint64(&d.sz, 0)

	// Close file and rename temporary file.
	if err = d.f.Close(); err != nil {
		return
	}
	err = os.Rename(d.ft, d.fd)
	d.f = nil

	return
}

func (d *Writer) init() {
	d.err = nil
	if len(d.Directory) == 0 {
		d.err = ErrNoDestinationDir
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
	if d.bs > 0 {
		d.buf = make([]byte, 0, d.bs)
	}
}

func (d *Writer) flushBuf() (err error) {
	lo, hi := 4, len(d.buf)
	if d.f == nil {
		d.buf = append(d.buf, d.dir...)
		d.buf = append(d.buf, os.PathSeparator)
		if d.buf, err = clock.AppendFormat(d.buf, d.mask, time.Now()); err != nil {
			return
		}
		filepath := fastconv.B2S(d.buf[hi:])
		d.fd = bytealg.CopyStr(filepath)
		d.buf = append(d.buf, ".tmp"...)
		filepathTmp := fastconv.B2S(d.buf[hi:])
		d.ft = bytealg.CopyStr(filepathTmp)
		if d.f, err = os.Create(filepathTmp); err != nil {
			return
		}
		lo = 0
	}

	p := d.buf[lo:hi]
	for len(p) >= flushChunkSize {
		if _, err = d.f.Write(p[:flushChunkSize]); err != nil {
			return
		}
		p = p[flushChunkSize:]
	}
	if len(p) > 0 {
		if _, err = d.f.Write(p); err != nil {
			return
		}
	}
	d.buf = d.buf[:4]

	return
}
