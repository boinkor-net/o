// Package ringio implements a ring-buffer that is an io.Reader and an
// io.Writer with fixed-size semantics.
package ringio

import (
	"sync"

	"github.com/antifuchs/o"
)

// Bounded is an io.Reader and io.Writer that allows writing as many
// bytes as are given for the capacity before it has to be drained by
// reading from it.
//
// It is able to safely read and write in parallel, protected by a
// Mutex.
type Bounded struct {
	sync.Mutex
	r         o.Ring
	buf       []byte
	overwrite bool
}

type byteSlice []byte

func (bs byteSlice) Len() int {
	return len(bs)
}

// New returns a bounded ring buffer of the given capacity. If
// overwrite is true, a full ring buffer will discard unread bytes and
// overwrite them upon writes.
//
// If overwrite is false, writing more bytes than there is space in
// the buffer will fail with ErrFull and no bytes will be written.
func New(cap uint, overwrite bool) *Bounded {
	buf := make([]byte, cap)
	ring := o.NewRingForSlice(byteSlice(buf))
	return &Bounded{r: ring, buf: buf, overwrite: overwrite}
}

func (b *Bounded) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	n = len(p)
	reserve := uint(len(p))
	remaining := b.r.Capacity() - b.r.Size()
	if remaining < uint(len(p)) {
		if !b.overwrite {
			return 0, o.ErrFull
		}
		// consume the bytes that we're over and reset input
		// to fit:
		p = p[reserve-b.r.Capacity():]
		for i := uint(0); i <= b.r.Size(); i++ {
			_, _ = b.r.Shift()
		}
		reserve = uint(len(p))
	}
	first, second, _ := b.r.PushN(reserve)
	copy(b.buf[first.Start:first.End], p[0:first.Length()])
	copy(b.buf[second.Start:second.End], p[first.Length():])
	return
}

func (b *Bounded) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	if b.r.Empty() {
		return 0, nil
	}

	n = int(b.r.Size())
	if n > len(p) {
		n = len(p)
	}
	var first, second o.Range
	first, second, err = b.r.ShiftN(uint(n))
	copy(p[0:first.Length()], b.buf[first.Start:first.End])
	copy(p[first.Length():], b.buf[second.Start:second.End])
	return
}

func (b *Bounded) reset() {
	b.r = o.NewRingForSlice(byteSlice(b.buf))
}

// Reset throws away all data on the ring buffer.
func (b *Bounded) Reset() {
	b.Lock()
	defer b.Unlock()
	b.reset()
}

// Bytes consumes all readable data on the ring buffer and returns a
// newly-allocated byte slice containing all readable bytes.
func (b *Bounded) Bytes() []byte {
	b.Lock()
	defer b.Unlock()

	first, second := b.r.Consume()
	val := make([]byte, first.Length()+second.Length())
	copy(val, b.buf[first.Start:first.End])
	copy(val[first.End:], b.buf[second.Start:second.End])
	return val
}

// String consumes all readable data on the ring buffer and returns it
// as a string.
func (b *Bounded) String() string {
	return string(b.Bytes())
}
