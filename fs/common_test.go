package fs

import "bytes"

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

var testVars []interface{}

func init() {
	p := []byte("asdfgh456")
	s := "qwerty789"
	var buf bytes.Buffer
	buf.WriteString("qweasdzxc2468")
	m := m8r{payload: []byte("ertfghvbn123")}
	testVars = []interface{}{
		[]byte("foobar123"),
		&p,
		"zxcvbn051",
		&s,
		&buf,
		&m,
	}
}
