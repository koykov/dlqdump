package dlqdump

import "github.com/koykov/blqueue"

type Restorer interface {
	Restore(dst blqueue.Interface) error
}
