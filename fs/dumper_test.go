package fs

import (
	"testing"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/encoder"
)

func TestDumper(t *testing.T) {
	t.Run("force", func(t *testing.T) {
		q, err := dlqdump.NewDLQ(&dlqdump.DLQConfig{
			Version:  0,
			Key:      "stage0",
			Capacity: dlqdump.Byte * 512,
			Encoder:  encoder.Basic{},
			Dumper: &Dumper{
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
		q, err := dlqdump.NewDLQ(&dlqdump.DLQConfig{
			Version:  0,
			Key:      "stage0",
			Capacity: dlqdump.Byte * 32,
			Encoder:  encoder.Basic{},
			Dumper: &Dumper{
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
		q, err := dlqdump.NewDLQ(&dlqdump.DLQConfig{
			Version:       0,
			Key:           "stage0",
			Capacity:      dlqdump.Byte * 512,
			FlushInterval: time.Millisecond * 10,
			Encoder:       encoder.Basic{},
			Dumper: &Dumper{
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
