package dlqdump

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/koykov/bitset"
	"github.com/koykov/blqueue"
	"github.com/koykov/bytealg"
)

type Queue struct {
	bitset.Bitset
	config *Config
	status blqueue.Status

	once  sync.Once
	timer *time.Timer
	mux   sync.Mutex
	buf   []byte

	Err error
}

func New(config *Config) (*Queue, error) {
	q := &Queue{
		config: config.Copy(),
	}
	q.once.Do(q.init)
	return q, q.Err
}

func (q *Queue) Enqueue(x interface{}) (err error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == blqueue.StatusClose || status == blqueue.StatusFail {
		return blqueue.ErrQueueClosed
	}

	q.mux.Lock()
	defer q.mux.Unlock()

	ho := len(q.buf)
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	q.buf = append(q.buf, '0')
	po := len(q.buf)

	if q.config.Encoder != nil {
		q.buf, err = q.config.Encoder.Encode(q.buf, x)
	} else {
		switch x.(type) {
		case []byte:
			q.buf = append(q.buf, x.([]byte)...)
		case *[]byte:
			q.buf = append(q.buf, *x.(*[]byte)...)
		case string:
			q.buf = append(q.buf, x.(string)...)
		case *string:
			q.buf = append(q.buf, *x.(*string)...)
		case MarshallerTo:
			off := len(q.buf)
			m := x.(MarshallerTo)
			q.buf = bytealg.GrowDelta(q.buf, m.Size())
			_, err = m.MarshalTo(q.buf[off:])
		case Marshaller:
			m := x.(Marshaller)
			var b []byte
			if b, err = m.Marshal(); err == nil {
				q.buf = append(q.buf, b...)
			}
		case fmt.Stringer:
			q.buf = append(q.buf, x.(fmt.Stringer).String()...)
		default:
			err = ErrUnknownMarshaller
		}
	}
	if err != nil {
		q.buf = q.buf[:ho]
		return
	}

	pl := len(q.buf) - po
	binary.LittleEndian.PutUint32(q.buf[ho:ho+4], uint32(pl))

	return
}

func (q *Queue) Rate() float32 {
	return 0
}

func (q *Queue) Close() error {
	if q.getStatus() == blqueue.StatusClose {
		return blqueue.ErrQueueClosed
	}
	return nil
}

func (q *Queue) init() {
	c := q.config

	if len(c.Key) == 0 {
		q.Err = blqueue.ErrNoKey
		q.setStatus(blqueue.StatusFail)
		return
	}
	if c.Size == 0 {
		q.Err = blqueue.ErrNoSize
		q.setStatus(blqueue.StatusFail)
		return
	}

	q.setStatus(blqueue.StatusActive)
}

func (q *Queue) setStatus(status blqueue.Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

func (q *Queue) getStatus() blqueue.Status {
	return blqueue.Status(atomic.LoadUint32((*uint32)(&q.status)))
}
