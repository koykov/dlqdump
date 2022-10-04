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

type Reader struct {
	// MatchMask represents pattern to match the names of dump files.
	// Mandatory param.
	MatchMask string

	once sync.Once
	mux  sync.Mutex
	mask string
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
		var matches []string
		if matches, err = filepath.Glob(r.mask); err != nil {
			return 0, dst, err
		}
		if len(matches) == 0 {
			return 0, dst, io.EOF
		}

		filename := matches[0]
		if r.f, err = os.OpenFile(filename, os.O_RDONLY, 0644); err != nil {
			return 0, dst, err
		}
		r.buf = bytealg.Grow(r.buf, 8)
		if _, err = io.ReadAtLeast(r.f, r.buf, 8); err != nil {
			return 0, dst, r.wrapErr(err)
		}
		r.ver = dlqdump.Version(binary.LittleEndian.Uint64(r.buf))
	}

	r.buf = bytealg.Grow(r.buf, 4)
	if _, err = io.ReadAtLeast(r.f, r.buf, 4); err != nil {
		return r.ver, dst, r.wrapErr(err)
	}
	pl := binary.LittleEndian.Uint32(r.buf)
	r.buf = bytealg.Grow(r.buf, int(pl))
	if _, err = io.ReadAtLeast(r.f, r.buf, int(pl)); err != nil {
		return r.ver, dst, r.wrapErr(err)
	}
	dst = append(dst, r.buf...)

	return r.ver, dst, nil
}

func (r *Reader) init() {
	r.mask = r.MatchMask
}

func (r *Reader) wrapErr(err error) error {
	if err == io.EOF {
		_ = r.f.Close()
		r.f = nil
		r.ver = 0
	}
	return err
}
