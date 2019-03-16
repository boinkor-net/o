// package ringio implements a ring-buffer that is an io.Reader and an
// io.Writer with fixed-size semantics.
package ringio

import (
	"io"

	"github.com/antifuchs/o"
)

// Bounded is an io.Reader and io.Writer that allows writing as many
// bytes as are given for the capacity before it has to be drained.
type Bounded struct {
	o.Ring
	buf       []byte
	overwrite bool
}

func New(cap uint, overwrite bool) *Bounded {
	buf := make([]byte, cap)
	return &Bounded{Ring: o.NewRing(cap), buf: buf, overwrite: overwrite}
}

func (b *Bounded) Write(p []byte) (n int, err error) {
	var i uint
	for n, c := range p {
		if b.overwrite {
			i, err = b.Ring.Push()
			if err == o.ErrFull {
				return n, io.ErrShortWrite
			}
		} else {
			i = o.ForcePush(b.Ring)
		}
		b.buf[i] = c
	}
	return len(p), nil
}

func (b *Bounded) Read(p []byte) (n int, err error) {
	if b.Ring.Empty() {
		return 0, io.EOF
	}
	var i uint
	for {
		i, err = b.Ring.Shift()
		if err == o.ErrEmpty {
			return n, nil
		}
		p[n] = b.buf[i]
		n++
	}
	return
}
