// package ringio implements a ring-buffer that is an io.Reader and an
// io.Writer with fixed-size semantics.
package ringio

import (
	"io"
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
// Otherwise, writes on a full ring buffer will fill up the buffer
// with as much data as they can and then return io.ErrShortWrite.
func New(cap uint, overwrite bool) *Bounded {
	buf := make([]byte, cap)
	ring := o.NewRingForSlice(byteSlice(buf))
	return &Bounded{r: ring, buf: buf, overwrite: overwrite}
}

func (b *Bounded) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	n = len(p)
	remaining := uint(len(b.buf)) - b.r.Size()
	if b.overwrite && remaining < uint(len(p)) {
		// consume the bytes that we're over and reset input
		// to fit:
		p = p[len(p)-len(b.buf) : len(p)]
		for i := uint(0); i <= b.r.Size(); i++ {
			b.r.Shift()
		}
	}
	first, second, resErr := o.Reserve(b.r, uint(len(p)))
	if !b.overwrite {
		n = int(first.Length() + second.Length())
		if resErr != nil {
			err = io.ErrShortWrite
		}
	}
	copy(b.buf[first.Start:first.End], p[0:first.Length()])
	copy(b.buf[second.Start:second.End], p[first.Length():len(p)])
	return
}

func (b *Bounded) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	if b.r.Empty() {
		return 0, io.EOF
	}

	var i uint
	for {
		if n >= len(p) {
			return
		}
		i, err = b.r.Shift()
		if err == o.ErrEmpty {
			return n, nil
		}
		p[n] = b.buf[i]
		n++
	}
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
