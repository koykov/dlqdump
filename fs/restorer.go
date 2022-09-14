package fs

import "github.com/koykov/blqueue"

type Restorer struct{}

func (r Restorer) Restore(dst blqueue.Interface) error {
	_ = dst
	return nil
}
