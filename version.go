package dlqdump

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represent simple version container.
type Version uint64

// NewVersion composes version from given parts.
func NewVersion(major, minor, patch, revision uint16) Version {
	var v uint64
	v = v | uint64(major)<<48
	v = v | uint64(minor)<<32
	v = v | uint64(patch)<<16
	v = v | uint64(revision)
	return Version(v)
}

// ParseVersion makes new version from source string.
//
// See version_test.go for examples.
func ParseVersion(ver string) Version {
	var m, n, p, r uint16
	c := 0
	for {
		i := strings.Index(ver, ".")
		if i == -1 {
			break
		}
		raw := ver[:i]
		u, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return 0
		}
		switch c {
		case 0:
			m = uint16(u)
		case 1:
			n = uint16(u)
		case 2:
			p = uint16(u)
		}
		c++
		ver = ver[i+1:]
	}
	u, err := strconv.ParseUint(ver, 10, 64)
	if err != nil {
		return 0
	}
	switch c {
	case 0:
		m = uint16(u)
	case 1:
		n = uint16(u)
	case 2:
		p = uint16(u)
	case 3:
		r = uint16(u)
	}
	return NewVersion(m, n, p, r)
}

func (v Version) Major() uint16 {
	return uint16(v >> 48)
}

func (v Version) Minor() uint16 {
	return uint16(v >> 32)
}

func (v Version) Patch() uint16 {
	return uint16(v >> 16)
}

func (v Version) Revision() uint16 {
	return uint16(v)
}

func (v Version) String() string {
	m, n, p, r := v.Major(), v.Minor(), v.Patch(), v.Revision()
	return fmt.Sprintf("%d.%d.%d.%d", m, n, p, r)
}
