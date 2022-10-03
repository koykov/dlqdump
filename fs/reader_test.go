package fs

import (
	"testing"

	"github.com/koykov/dlqdump"
	"github.com/koykov/dlqdump/decoder"
)

func TestReader(t *testing.T) {
	conf := dlqdump.Config{
		Version: dlqdump.ParseVersion("1.0"),
		Key:     "example",
		Reader: &Reader{
			MatchMask: "testdata/*.bin",
		},
		Decoder: decoder.Fallthrough{},
		Queue:   testq{t: t},
	}
	rst, err := dlqdump.NewRestorer(&conf)
	if err != nil {
		t.Fatal(err)
	}

	if err = rst.Restore(); err != nil {
		t.Fatal(err)
	}
}
