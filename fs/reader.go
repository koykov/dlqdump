package fs

import (
	"encoding/binary"
	"io"
	"os"
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
		if r.f, err = os.OpenFile("", os.O_RDONLY, 0644); err != nil {
			return 0, dst, err
		}
		r.buf = bytealg.Grow(r.buf, 8)
		if _, err = io.ReadAtLeast(r.f, r.buf, 8); err != nil {
			return 0, dst, err
		}
		r.ver = dlqdump.Version(binary.LittleEndian.Uint64(r.buf))
	}

	r.buf = bytealg.Grow(r.buf, 4)
	if _, err = io.ReadAtLeast(r.f, r.buf, 4); err != nil {
		return r.ver, dst, err
	}
	pl := binary.LittleEndian.Uint32(r.buf)
	r.buf = bytealg.Grow(r.buf, int(pl))
	if _, err = io.ReadAtLeast(r.f, r.buf, int(pl)); err != nil {
		return r.ver, dst, err
	}

	return r.ver, dst, nil
}

func (r *Reader) init() {
	r.mask = r.MatchMask
}
