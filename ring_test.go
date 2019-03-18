package o

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForcePush(t *testing.T) {
	r := NewRing(1)
	r.Push()
	assert.Equal(t, r.ForcePush(), uint(0))
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

func TestSlices(t *testing.T) {
	tests := []struct {
		name  string
		slice Slice
		len   uint
	}{
		{"[]int", sort.IntSlice([]int{1, 2, 3, 4}), 4},
		{"[]float64", sort.Float64Slice([]float64{1.0, 2.0, 3.0}), 3},
		{"[]string", sort.StringSlice([]string{"hi", "there", "farts", "yup", "strings"}), 5},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ring := NewRingForSlice(test.slice)
			assert.Equal(t, ring.capacity(), test.len)
		})
	}
}

// TestErrors is silly but improves coverage metrics.
func TestErrors(t *testing.T) {
	assert.Equal(t, ErrEmpty.Error(), "reading from an empty ring")
	assert.Equal(t, ErrFull.Error(), "inserting into a full ring")
}
