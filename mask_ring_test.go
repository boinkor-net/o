package o

import (
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

const power = 16
const N = 1 << power

func TestMaskPush(t *testing.T) {
	r := NewRing(N)
	var i uint
	for ; i < N; i++ {
		n, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	_, err := r.Push()
	assert.Error(t, err)
}

func TestMaskShift(t *testing.T) {
	r := NewRing(N)
	_, err := r.Shift()
	assert.Error(t, err)

	var i uint
	for ; i < N; i++ {
		n, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	for i = 0; i < N; i++ {
		n, err := r.Shift()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	_, err = r.Shift()
	assert.Error(t, err)
}

func BenchmarkMaskRing(b *testing.B) {
	n := 1 << uint(bits.Len(uint(b.N))+1)
	r := NewRing(uint(n))
	var i uint
	for ; i < 1<<uint(b.N); i++ {
		r.Push()
	}
	for i = 0; i < 1<<uint(b.N); i++ {
		r.Shift()
	}
}
