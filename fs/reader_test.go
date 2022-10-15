package fs

import (
	"testing"
	"time"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/decoder"
)

func TestReader(t *testing.T) {
	conf := dlqdump.Config{
		Version:       dlqdump.NewVersion(1, 0, 0, 0),
		CheckInterval: time.Hour, // will never happen during testing
		Reader: &Reader{
			MatchMask: "testdata/*.bin",
			OnEOF:     func(_ string) error { return nil }, // dummy EOF function to keep testing dump file
		},
		Decoder: decoder.Fallthrough{},
		Queue:   &testq{t: t},
	}
	rst, err := dlqdump.NewRestorer(&conf)
	if err != nil {
		t.Fatal(err)
	}

	if err = rst.Restore(); err != nil {
		t.Fatal(err)
	}
	_ = rst.ForceClose()
}
