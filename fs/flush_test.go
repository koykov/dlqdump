package fs

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/decoder"
)

func readAllFS(t *testing.T, dir string) [][]byte {
	t.Helper()
	var got [][]byte
	r := &Reader{MatchMask: filepath.Join(dir, "*.bin")}
	for {
		_, buf, err := r.Read(nil)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		x, err := (decoder.Fallthrough{}).Decode(buf)
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		got = append(got, append([]byte(nil), x.([]byte)...))
	}
	return got
}

func TestFlush(t *testing.T) {
	cases := []struct {
		name  string
		items [][]byte
	}{
		{name: "empty", items: nil},
		{name: "single", items: [][]byte{[]byte("abc")}},
		{name: "multiple", items: [][]byte{[]byte("a"), []byte("bb"), []byte("ccc")}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			w := &Writer{Directory: dir}
			ver := dlqdump.NewVersion(1, 0, 0, 0)
			for _, p := range c.items {
				if _, err := w.Write(ver, p); err != nil {
					t.Fatalf("write: %v", err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("flush: %v", err)
			}
			// Flush must be idempotent and safe on an already flushed (or empty) writer.
			if err := w.Flush(); err != nil {
				t.Fatalf("second flush: %v", err)
			}
			if got := w.Size(); got != 0 {
				t.Errorf("Size() after flush = %d, want 0", got)
			}

			got := readAllFS(t, dir)
			if len(got) != len(c.items) {
				t.Fatalf("read %d items, want %d", len(got), len(c.items))
			}
			for i := range c.items {
				if string(got[i]) != string(c.items[i]) {
					t.Errorf("item %d = %q, want %q", i, got[i], c.items[i])
				}
			}
		})
	}
}
