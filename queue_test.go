package dlqdump

import (
	"bytes"
	"testing"
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

func TestQueue(t *testing.T) {
	t.Run("stage0", func(t *testing.T) {
		q, err := New(&Config{
			Version:   0,
			Key:       "stage0",
			Size:      Byte * 512,
			Directory: "dump",
		})
		if err != nil {
			t.Fatal(err)
		}
		p := []byte("asdfgh")
		s := "qwerty"
		var buf bytes.Buffer
		buf.WriteString("qweasdzxc")
		m := m8r{payload: []byte("ertfghvbn123")}
		var vars = []interface{}{
			[]byte("foobar"),
			&p,
			"zxcvbn",
			&s,
			&buf,
			&m,
		}
		for i := 0; i < len(vars); i++ {
			if err = q.Enqueue(vars[i]); err != nil {
				t.Fatal(err)
			}
		}
	})
}
