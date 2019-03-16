package o

import (
	"math/bits"
	"testing"

	"github.com/stretchr/testify/assert"
)

const power = 16
const N = 1 << power

func TestMaskPush(t *testing.T) {
	r := NewPowerOfTwo(power)
	var i uint
	for ; i < N; i++ {
		new, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	_, err := r.Push()
	assert.Error(t, err)
}

func TestMaskShift(t *testing.T) {
	r := NewPowerOfTwo(power)
	_, err := r.Shift()
	assert.Error(t, err)

	var i uint
	for ; i < N; i++ {
		new, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	for i = 0; i < N; i++ {
		new, err := r.Shift()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	_, err = r.Shift()
	assert.Error(t, err)
}

func BenchmarkMaskRing(b *testing.B) {
	power := uint(bits.Len(uint(b.N)) + 1)
	r := NewPowerOfTwo(power)
	var i uint
	for ; i < 1<<uint(b.N); i++ {
		r.Push()
	}
	for i = 0; i < 1<<uint(b.N); i++ {
		r.Shift()
	}
}
