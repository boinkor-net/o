package o

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const BasicN = 20

func TestPush(t *testing.T) {
	r := NewRing(BasicN)
	var i uint
	for ; i < 20; i++ {
		n, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	_, err := r.Push()
	assert.Error(t, err)
}

func TestShift(t *testing.T) {
	r := NewRing(BasicN)
	_, err := r.Shift()
	assert.Error(t, err)

	var i uint
	for ; i < 20; i++ {
		n, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	for i = 0; i < 20; i++ {
		n, err := r.Shift()
		assert.NoError(t, err)
		assert.Equal(t, i, n)
	}
	_, err = r.Shift()
	assert.Error(t, err)
}

func BenchmarkBasicRing(b *testing.B) {
	r := NewRing(uint(b.N))
	var i uint
	for ; i < uint(b.N); i++ {
		r.Push()
	}
	for i = 0; i < uint(b.N); i++ {
		r.Shift()
	}
}
