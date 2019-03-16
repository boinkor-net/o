package o

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForcePush(t *testing.T) {
	r := NewRing(1)
	r.Push()
	assert.Equal(t, ForcePush(r), uint(0))
}

func TestPushAndShift(t *testing.T) {
	tests := []struct {
		name  string
		cap   uint
		turns uint
	}{
		{"mask/3", 16, 3},
		{"basic/3", 19, 3},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ring := NewRing(test.cap)
			for i := uint(0); i < test.cap*test.turns; i++ {
				pushed, err := ring.Push()
				require.NoError(t, err)
				require.True(t, pushed < test.cap)

				shifted, err := ring.Shift()
				require.NoError(t, err)
				require.Equal(t, pushed, shifted, "on attempt %d", i)
			}
		})
	}
}
