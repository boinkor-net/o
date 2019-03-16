// package ringio implements a ring-buffer that is an io.Reader and an
// io.Writer with fixed-size semantics.
package ringio

import (
	"io"
	"math/bits"

	"github.com/antifuchs/o"
)

// Bounded is an io.Reader and io.Writer that allows writing as many
// bytes as are given for the capacity before it has to be drained.
type Bounded struct {
	o.Ring
	buf []byte
}

func New(cap uint) *Bounded {
	buf := make([]byte, cap)
	if cap%2 == 0 {
		r := o.NewPowerOfTwo(uint(bits.Len(cap)))
		return &Bounded{Ring: r, buf: buf}
	}
	return &Bounded{Ring: o.NewBasic(cap), buf: buf}
}

func (b *Bounded) Write(p []byte) (n int, err error) {
	var i uint
	for n, c := range p {
		i, err = b.Ring.Push()
		if err == o.ErrFull {
			return n, io.ErrShortWrite
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
