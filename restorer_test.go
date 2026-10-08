package dlqdump_test

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/decoder"
	"github.com/koykov/dlqdump/fs"
)

// destQ is a mock destination queue for the Restorer.
type destQ struct {
	mu        sync.Mutex
	got       []string
	rate      float32
	rateCalls int32
}

func (d *destQ) Enqueue(x any) error {
	d.mu.Lock()
	d.got = append(d.got, string(x.([]byte)))
	d.mu.Unlock()
	return nil
}
func (d *destQ) Size() int         { return 0 }
func (d *destQ) Capacity() int     { return 1000 }
func (d *destQ) Close() error      { return nil }
func (d *destQ) ForceClose() error { return nil }
func (d *destQ) Rate() float32 {
	atomic.AddInt32(&d.rateCalls, 1)
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.rate
}
func (d *destQ) setRate(r float32) {
	d.mu.Lock()
	d.rate = r
	d.mu.Unlock()
}
func (d *destQ) delivered() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.got...)
}

func writeDump(t *testing.T, dir, prefix string, items []string) {
	t.Helper()
	w := &fs.Writer{Directory: dir, FileMask: prefix + "--%Y-%m-%d--%H-%M-%S--%N.bin"}
	ver := dlqdump.NewVersion(1, 0, 0, 0)
	for _, s := range items {
		if _, err := w.Write(ver, []byte(s)); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
}

func newTestRestorer(t *testing.T, dir string, dq *destQ, closeTimeout time.Duration) *dlqdump.Restorer {
	t.Helper()
	r, err := dlqdump.NewRestorer(&dlqdump.Config{
		Version:          dlqdump.NewVersion(1, 0, 0, 0),
		CheckInterval:    time.Hour, // background ticker must not interfere with the tests
		PostponeInterval: time.Millisecond,
		AllowRate:        .95,
		CloseTimeout:     closeTimeout,
		Reader:           &fs.Reader{MatchMask: filepath.Join(dir, "*.bin")},
		Decoder:          decoder.Fallthrough{},
		Queue:            dq,
	})
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func waitRestorer(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition was not met in time")
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestRestorerClose(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		block bool
		want  []string
	}{
		{
			name:  "delivers_all",
			items: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "close_waits_for_in_flight_item",
			items: []string{"a", "b", "c"},
			block: true,
			want:  []string{"a"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			writeDump(t, dir, "f1", c.items)
			dq := &destQ{}
			if c.block {
				dq.setRate(1)
			}
			r := newTestRestorer(t, dir, dq, 5*time.Second)

			done := make(chan error, 1)
			go func() { done <- r.Restore() }()

			if c.block {
				// Wait until the first item is read and Restore is waiting for the destination queue.
				waitRestorer(t, func() bool { return atomic.LoadInt32(&dq.rateCalls) > 0 })

				closed := make(chan error, 1)
				go func() { closed <- r.Close() }()
				// Close must not return while the in-flight item can't be delivered.
				select {
				case <-closed:
					t.Fatal("Close returned before the in-flight item was delivered")
				case <-time.After(100 * time.Millisecond):
				}
				// Let the destination queue accept items.
				dq.setRate(0)
				select {
				case err := <-closed:
					if err != nil {
						t.Fatalf("close: %v", err)
					}
				case <-time.After(3 * time.Second):
					t.Fatal("Close did not return after the destination queue freed up")
				}
				<-done
			} else {
				<-done
				if err := r.Close(); err != nil {
					t.Fatalf("close: %v", err)
				}
			}

			if got := dq.delivered(); !equalStrings(got, c.want) {
				t.Fatalf("delivered %v, want %v", got, c.want)
			}
		})
	}
}

func TestRestorerCloseStopsReadingNewFiles(t *testing.T) {
	dir := t.TempDir()
	writeDump(t, dir, "f1", []string{"a", "b", "c"})
	writeDump(t, dir, "f2", []string{"z"})

	dq := &destQ{}
	dq.setRate(1)
	r := newTestRestorer(t, dir, dq, 5*time.Second)
	go r.Restore()
	waitRestorer(t, func() bool { return atomic.LoadInt32(&dq.rateCalls) > 0 })

	closed := make(chan error, 1)
	go func() { closed <- r.Close() }()
	dq.setRate(0)
	select {
	case err := <-closed:
		if err != nil {
			t.Fatalf("close: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Close did not return")
	}

	// The second file must never be touched (not read, not removed).
	matches, _ := filepath.Glob(filepath.Join(dir, "f2*.bin"))
	if len(matches) != 1 {
		t.Fatalf("second dump file was read/removed: %v", matches)
	}
}

func TestRestorerCloseTimeout(t *testing.T) {
	dir := t.TempDir()
	writeDump(t, dir, "f", []string{"a"})
	dq := &destQ{}
	dq.setRate(1)
	r := newTestRestorer(t, dir, dq, 10*time.Millisecond)
	go r.Restore()
	waitRestorer(t, func() bool { return atomic.LoadInt32(&dq.rateCalls) > 0 })

	if err := r.CloseWithTimeout(10 * time.Millisecond); err != dlqdump.ErrTimeout {
		t.Errorf("CloseWithTimeout = %v, want %v", err, dlqdump.ErrTimeout)
	}
	dq.setRate(0)
}

func TestReaderCloseKeepsFile(t *testing.T) {
	dir := t.TempDir()
	writeDump(t, dir, "f", []string{"a", "b"})

	r := &fs.Reader{MatchMask: filepath.Join(dir, "*.bin")}
	if _, _, err := r.Read(nil); err != nil {
		t.Fatalf("read: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// The file must be kept (no data loss); it will be re-read on the next restore (at-least-once).
	matches, _ := filepath.Glob(filepath.Join(dir, "*.bin"))
	if len(matches) != 1 {
		t.Fatalf("dump file was removed on Close: %v", matches)
	}
	_, buf, err := r.Read(nil)
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	if string(buf) != "a" {
		t.Errorf("re-read %q, want %q", buf, "a")
	}
}
