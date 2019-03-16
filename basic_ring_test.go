package o

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	r := NewBasic(20)
	var i uint
	for ; i < 20; i++ {
		new, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	_, err := r.Push()
	assert.Error(t, err)
}

func TestShift(t *testing.T) {
	r := NewBasic(20)
	_, err := r.Shift()
	assert.Error(t, err)

	var i uint
	for ; i < 20; i++ {
		new, err := r.Push()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	for i = 0; i < 20; i++ {
		new, err := r.Shift()
		assert.NoError(t, err)
		assert.Equal(t, i, new)
	}
	_, err = r.Shift()
	assert.Error(t, err)
}

func BenchmarkBasicRing(b *testing.B) {
	r := NewBasic(uint(b.N))
	var i uint
	for ; i < uint(b.N); i++ {
		r.Push()
	}
	for i = 0; i < uint(b.N); i++ {
		r.Shift()
	}
}
