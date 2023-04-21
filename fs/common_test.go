package fs

import (
	"bytes"
	"errors"
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

type testq struct {
	t *testing.T
	c int
}

func (q *testq) Enqueue(x any) error {
	raw, ok := x.([]byte)
	if !ok {
		return errors.New("type mismatch")
	}
	s := string(raw)
	switch q.c {
	case 0:
		if !bytes.Equal(raw, testVars[q.c].([]byte)) {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	case 1:
		if !bytes.Equal(raw, *testVars[q.c].(*[]byte)) {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	case 2:
		if s != testVars[q.c].(string) {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	case 3:
		if s != *testVars[q.c].(*string) {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	case 4:
		if s != testVars[q.c].(*bytes.Buffer).String() {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	case 5:
		if !bytes.Equal(raw, testVars[q.c].(*m8r).payload) {
			q.t.Fatalf("index %d: value mismatch", q.c)
		}
	}
	q.c++
	return nil
}
func (q *testq) Size() int     { return 0 }
func (q *testq) Capacity() int { return 0 }
func (q *testq) Rate() float32 { return 0 }
func (q *testq) Close() error  { return nil }
func (q *testq) reset()        { q.t, q.c = nil, 0 }

var testVars []any

func init() {
	p := []byte("asdfgh456")
	s := "qwerty789"
	var buf bytes.Buffer
	buf.WriteString("qweasdzxc2468")
	m := m8r{payload: []byte("ertfghvbn123")}
	testVars = []any{
		[]byte("foobar123"),
		&p,
		"zxcvbn051",
		&s,
		&buf,
		&m,
	}
}
