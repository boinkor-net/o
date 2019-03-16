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
// It is able to safely read and write in parallel, protected by a Mutex.
type Bounded struct {
	r         o.Ring
	buf       []byte
	overwrite bool
}

// New returns a bounded ring buffer of the given capacity. If
// overwrite is true, a full ring buffer will discard unread bytes and
// overwrite them upon writes. Otherwise, writes on a full ring buffer
// will fail.
func New(cap uint, overwrite bool) *Bounded {
	buf := make([]byte, cap)
	return &Bounded{r: o.NewRing(cap), buf: buf, overwrite: overwrite}
}

func (b *Bounded) Write(p []byte) (n int, err error) {
	var i uint
	for n, c := range p {
		if b.overwrite {
			i = o.ForcePush(b.r)
		} else {
			i, err = b.r.Push()
			if err == o.ErrFull {
				return n, io.ErrShortWrite
			}
		}
		b.buf[i] = c
	}
	return len(p), nil
}

func (b *Bounded) Read(p []byte) (n int, err error) {
	if b.r.Empty() {
		return 0, io.EOF
	}
	var i uint
	for {
		i, err = b.r.Shift()
		if err == o.ErrEmpty {
			return n, nil
		}
		p[n] = b.buf[i]
		n++
	}
	return
}
