package dlqdump

import (
	"bytes"
	"testing"
	"time"
)

type m8r struct {
	payload []byte
}

func (m m8r) Size() int {
	return len(m.payload)
}

func (m m8r) MarshalTo(p []byte) (int, error) {
	copy(p, m.payload)
	return m.Size(), nil
}

func TestFlush(t *testing.T) {
	p := []byte("asdfgh456")
	s := "qwerty789"
	var buf bytes.Buffer
	buf.WriteString("qweasdzxc2468")
	m := m8r{payload: []byte("ertfghvbn123")}
	var vars = []interface{}{
		[]byte("foobar123"),
		&p,
		"zxcvbn051",
		&s,
		&buf,
		&m,
	}

	t.Run("force", func(t *testing.T) {
		q, err := New(&Config{
			Version:   0,
			Key:       "stage0",
			Size:      Byte * 512,
			Directory: "dump",
			FileMask:  "force--%Y-%m-%d--%H-%M-%S--%N.bin",
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(vars); i++ {
			if err = q.Enqueue(vars[i]); err != nil {
				t.Fatal(err)
			}
		}
		// Force flush will trigger on close queue with unflushed data.
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("size", func(t *testing.T) {
		q, err := New(&Config{
			Version:   0,
			Key:       "stage0",
			Size:      Byte * 32,
			Directory: "dump",
			FileMask:  "size--%Y-%m-%d--%H-%M-%S--%N.bin",
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(vars); i++ {
			if err = q.Enqueue(vars[i]); err != nil {
				t.Fatal(err)
			}
		}
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("timer", func(t *testing.T) {
		q, err := New(&Config{
			Version:   0,
			Key:       "stage0",
			Size:      Byte * 512,
			TimeLimit: time.Millisecond * 10,
			Directory: "dump",
			FileMask:  "timer--%Y-%m-%d--%H-%M-%S--%N.bin",
		})
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(vars); i++ {
			if err = q.Enqueue(vars[i]); err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(time.Millisecond * 20)
		if err = q.Close(); err != nil {
			t.Fatal(err)
		}
	})
}
