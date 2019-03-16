package o

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForcePush(t *testing.T) {
	r := NewRing(1)
	r.Push()
	assert.Equal(t, ForcePush(r), uint(0))
}
