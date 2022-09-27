package dlqdump

import "testing"

func TestVersionParse(t *testing.T) {
	type tc struct {
		raw        string
		m, n, p, r uint16
	}
	var tcs = []tc{
		{"", 0, 0, 0, 0},
		{"0", 0, 0, 0, 0},
		{"1", 1, 0, 0, 0},
		{"1.0", 1, 0, 0, 0},
		{"1.0.1", 1, 0, 1, 0},
		{"1.0.1.7", 1, 0, 1, 7},
		{"5.12.134", 5, 12, 134, 0},
	}
	for _, c := range tcs {
		t.Run(c.raw, func(t *testing.T) {
			ver := ParseVersion(c.raw)
			if ver.Major() != c.m || ver.Minor() != c.n || ver.Patch() != c.p || ver.Revision() != c.r {
				t.FailNow()
			}
		})
	}
}
