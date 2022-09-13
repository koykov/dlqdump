package io

import "github.com/koykov/blqueue"

type FS struct{}

func (d FS) Dump(p []byte) (n int, err error) {
	_ = p
	return
}

func (d FS) Flush() error {
	return nil
}

func (d FS) Restore(dst blqueue.Interface) error {
	_ = dst
	return nil
}
