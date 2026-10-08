package dlqdump_test

import (
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/decoder"
	"github.com/koykov/dlqdump/encoder"
	"github.com/koykov/dlqdump/fs"
	"github.com/koykov/queue"
)

func newTestQueue(t *testing.T, dir string) *dlqdump.Queue {
	t.Helper()
	q, err := dlqdump.NewQueue(&dlqdump.Config{
		Version:       dlqdump.NewVersion(1, 0, 0, 0),
		Capacity:      dlqdump.Megabyte,
		FlushInterval: time.Hour,
		Encoder:       encoder.Builtin{},
		Writer:        &fs.Writer{Directory: dir, Buffer: dlqdump.Kilobyte},
	})
	if err != nil {
		t.Fatal(err)
	}
	return q
}

func readAll(t *testing.T, dir string) [][]byte {
	t.Helper()
	var got [][]byte
	r := &fs.Reader{MatchMask: filepath.Join(dir, "*.bin")}
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

func TestQueueClose(t *testing.T) {
	cases := []struct {
		name         string
		items        int
		after        func(t *testing.T, q *dlqdump.Queue) error
		wantAfterErr error
		wantItems    int
	}{
		{name: "flushes_on_close", items: 5, wantItems: 5},
		{name: "empty_close", items: 0, wantItems: 0},
		{
			name:         "enqueue_after_close",
			items:        3,
			after:        func(_ *testing.T, q *dlqdump.Queue) error { return q.Enqueue("late") },
			wantAfterErr: queue.ErrQueueClosed,
			wantItems:    3,
		},
		{
			name:         "second_close",
			items:        3,
			after:        func(_ *testing.T, q *dlqdump.Queue) error { return q.Close() },
			wantAfterErr: queue.ErrQueueClosed,
			wantItems:    3,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			q := newTestQueue(t, dir)
			for i := 0; i < c.items; i++ {
				if err := q.Enqueue(fmt.Sprintf("item-%d", i)); err != nil {
					t.Fatalf("enqueue: %v", err)
				}
			}
			if err := q.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}
			if got := q.Size(); got != 0 {
				t.Errorf("Size() after close = %d, want 0", got)
			}
			if c.after != nil {
				if err := c.after(t, q); err != c.wantAfterErr {
					t.Errorf("after close: err = %v, want %v", err, c.wantAfterErr)
				}
			}
			got := readAll(t, dir)
			if len(got) != c.wantItems {
				t.Fatalf("flushed %d items, want %d", len(got), c.wantItems)
			}
			for i, b := range got {
				if want := fmt.Sprintf("item-%d", i); string(b) != want {
					t.Errorf("item %d = %q, want %q", i, b, want)
				}
			}
		})
	}
}

func TestQueueCloseFailedInit(t *testing.T) {
	q, err := dlqdump.NewQueue(&dlqdump.Config{})
	if err == nil {
		t.Fatal("want init error, got nil")
	}
	if cerr := q.Close(); cerr != q.Err {
		t.Errorf("Close on failed queue = %v, want %v", cerr, q.Err)
	}
}

func TestQueueCloseConcurrent(t *testing.T) {
	cases := []struct {
		name       string
		goroutines int
		perG       int
	}{
		{name: "4x200", goroutines: 4, perG: 200},
		{name: "8x100", goroutines: 8, perG: 100},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			q := newTestQueue(t, dir)

			var accepted int64
			var wg sync.WaitGroup
			for g := 0; g < c.goroutines; g++ {
				wg.Add(1)
				go func(base int) {
					defer wg.Done()
					for i := 0; i < c.perG; i++ {
						if err := q.Enqueue(fmt.Sprintf("%d-%d", base, i)); err == nil {
							atomic.AddInt64(&accepted, 1)
						} else {
							return
						}
					}
				}(g)
			}

			// Close while the enqueues are in flight.
			time.Sleep(2 * time.Millisecond)
			if err := q.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}
			wg.Wait()

			got := readAll(t, dir)
			if n := atomic.LoadInt64(&accepted); int64(len(got)) != n {
				t.Fatalf("flushed %d items, accepted %d", len(got), n)
			}
		})
	}
}
