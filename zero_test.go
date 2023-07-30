package o_test

import (
	"testing"

	"github.com/antifuchs/o"
	"github.com/stretchr/testify/assert"
)

func newRing() o.Ring {
	return o.NewRing(0)
}

func TestZeroMeaningless(t *testing.T) {
	r := newRing()
	for i := 0; i < 2; i++ {
		assert.False(t, r.Empty())
		assert.True(t, r.Full())
		assert.Equal(t, uint(0), r.Size())
		assert.Equal(t, uint(0), r.Capacity())
		r.Consume()
	}
}

func TestZeroPush(t *testing.T) {
	r := newRing()
	var i uint

	n, err := r.Push()
	assert.Equal(t, o.ErrFull, err)
	assert.Equal(t, i, n)
}

func TestZeroShift(t *testing.T) {
	r := newRing()
	_, err := r.Shift()
	assert.Error(t, err)

	var i uint
	n, err := r.Push()
	assert.Equal(t, o.ErrFull, err)
	assert.Equal(t, uint(0), n)

	i, err = r.Shift()
	assert.Equal(t, o.ErrEmpty, err)
	assert.Equal(t, uint(0), i)
}

func BenchmarkZeroRing(b *testing.B) {
	r := newRing()
	var i uint
	for ; i < 1<<uint(b.N); i++ {
		r.Push()
	}
	for i = 0; i < 1<<uint(b.N); i++ {
		r.Shift()
	}
}
