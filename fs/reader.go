package fs

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/koykov/bytealg"
	"github.com/koykov/dlqdump"
)

// Reader is file system implementation of dlqdump.Reader interface.
type Reader struct {
	// MatchMask represents pattern to match the names of dump files.
	// Mandatory param.
	MatchMask string
	// OnEOF calls when EOF of current file reaches.
	// If this param omit os.Remove() will use by default.
	OnEOF func(filename string) error

	once sync.Once
	mux  sync.Mutex
	mask string
	eof  func(string) error
	fn   string
	f    *os.File
	ver  dlqdump.Version
	buf  []byte
}

func (r *Reader) Read(dst []byte) (dlqdump.Version, []byte, error) {
	r.once.Do(r.init)

	r.mux.Lock()
	defer r.mux.Unlock()

	var err error
	if r.f == nil {
		// No open file, so try to read first available dump file.
		var matches []string
		if matches, err = filepath.Glob(r.mask); err != nil {
			return 0, dst, err
		}
		if len(matches) == 0 {
			return 0, dst, io.EOF
		}

		r.fn = matches[0]
		if r.f, err = os.OpenFile(r.fn, os.O_RDONLY, 0644); err != nil {
			return 0, dst, err
		}
		// Read version from header (first 8 bytes).
		r.buf = bytealg.Grow(r.buf, 8)
		if _, err = io.ReadAtLeast(r.f, r.buf, 8); err != nil {
			return 0, dst, r.checkEOF(err)
		}
		r.ver = dlqdump.Version(binary.LittleEndian.Uint64(r.buf))
	}

	// Read item length bytes.
	r.buf = bytealg.Grow(r.buf, 4)
	if _, err = io.ReadAtLeast(r.f, r.buf, 4); err != nil {
		return r.ver, dst, r.checkEOF(err)
	}
	// Decode item length.
	pl := binary.LittleEndian.Uint32(r.buf)
	r.buf = bytealg.Grow(r.buf, int(pl))
	// Read item body.
	if _, err = io.ReadAtLeast(r.f, r.buf, int(pl)); err != nil {
		return r.ver, dst, r.checkEOF(err)
	}
	dst = append(dst, r.buf...)

	return r.ver, dst, nil
}

func (r *Reader) init() {
	r.mask = r.MatchMask
	if r.OnEOF == nil {
		r.OnEOF = os.Remove
	}
	r.eof = r.OnEOF
}

func (r *Reader) checkEOF(err error) error {
	if err == io.EOF {
		_ = r.f.Close()
		_ = r.eof(r.fn)
		r.fn = ""
		r.f = nil
		r.ver = 0
	}
	return err
}
