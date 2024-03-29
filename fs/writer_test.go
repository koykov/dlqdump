package fs

import (
	"testing"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/encoder"
)

type testEncoder struct {
	encoder.Builtin
	encoder.Marshaller
}

func (e testEncoder) Encode(dst []byte, x any) ([]byte, error) {
	var err error
	if dst, err = e.Builtin.Encode(dst, x); err == nil {
		return dst, nil
	}
	if dst, err = e.Marshaller.Encode(dst, x); err == nil {
		return dst, nil
	}
	return dst, encoder.ErrIncompatibleEncoder
}

func TestWriter(t *testing.T) {
	t.Run("force", func(t *testing.T) {
		q, err := dlqdump.NewQueue(&dlqdump.Config{
			Version:  dlqdump.NewVersion(1, 0, 0, 0),
			Capacity: dlqdump.Byte * 512,
			Encoder:  testEncoder{},
			Writer: &Writer{
				Directory: "testdata",
				FileMask:  "force--%Y-%m-%d--%H-%M-%S--%N.bin",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(testVars); i++ {
			if err = q.Enqueue(testVars[i]); err != nil {
				t.Fatal(err)
			}
		}
		// Force flush will trigger on close queue with unflushed data.
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("size", func(t *testing.T) {
		q, err := dlqdump.NewQueue(&dlqdump.Config{
			Version:  dlqdump.NewVersion(1, 0, 0, 0),
			Capacity: dlqdump.Byte * 32,
			Encoder:  testEncoder{},
			Writer: &Writer{
				Directory: "testdata",
				FileMask:  "size--%Y-%m-%d--%H-%M-%S--%N.bin",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(testVars); i++ {
			if err = q.Enqueue(testVars[i]); err != nil {
				t.Fatal(err)
			}
		}
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("timer", func(t *testing.T) {
		q, err := dlqdump.NewQueue(&dlqdump.Config{
			Version:       dlqdump.NewVersion(1, 0, 0, 0),
			Capacity:      dlqdump.Byte * 512,
			FlushInterval: time.Millisecond * 10,
			Encoder:       testEncoder{},
			Writer: &Writer{
				Directory: "testdata",
				FileMask:  "timer--%Y-%m-%d--%H-%M-%S--%N.bin",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(testVars); i++ {
			if err = q.Enqueue(testVars[i]); err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(time.Millisecond * 20)
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
}
