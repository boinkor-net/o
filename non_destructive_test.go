package o_test

import (
	"fmt"
	"testing"

	"github.com/antifuchs/o"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLIFO(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ra       o.Ring
		cycles   uint
		expected []uint
	}{
		{"basic5/13", o.NewRing(5), 13, []uint{3, 4, 0, 1, 2}},
		{"basic5/6", o.NewRing(5), 6, []uint{1, 2, 3, 4, 0}},
		{"mask4/13", o.NewRing(4), 13, []uint{1, 2, 3, 0}},
		{"mask4/6", o.NewRing(4), 6, []uint{2, 3, 0, 1}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			results := make([]uint, 0, len(test.expected))
			var i uint
			for i = 0; i < test.cycles; i++ {
				o.ForcePush(test.ra)
			}
			s := o.ScanLIFO(test.ra)
			for s.Next() {
				results = append(results, s.Value())
			}
			assert.Equal(t, test.expected, results)
		})
	}
}

func TestFIFO(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ra       o.Ring
		cycles   uint
		expected []uint
	}{
		{"basic5/13", o.NewRing(5), 13, []uint{2, 1, 0, 4, 3}},
		{"basic5/6", o.NewRing(5), 6, []uint{0, 4, 3, 2, 1}},
		{"mask4/13", o.NewRing(4), 13, []uint{0, 3, 2, 1}},
		{"mask4/6", o.NewRing(4), 6, []uint{1, 0, 3, 2}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			results := make([]uint, 0, len(test.expected))
			var i uint
			for i = 0; i < test.cycles; i++ {
				o.ForcePush(test.ra)
			}
			s := o.ScanFIFO(test.ra)
			for s.Next() {
				results = append(results, s.Value())
			}
			assert.Equal(t, test.expected, results)
		})
	}
}

func TestMatching(t *testing.T) {
	var toinsert = 263 // prime & greater than the max. capacity
	const min uint = 1
	const max uint = 256
	for n := min; n <= max; n++ {
		test := n
		t.Run(fmt.Sprintf("%03d", test), func(t *testing.T) {
			t.Parallel()
			ra := o.NewRing(test)
			for i := 0; i < toinsert; i++ {
				o.ForcePush(ra)
			}
			assert.Equal(t, ra.Size(), test)

			fifo := make([]uint, 0, test)
			lifo := make([]uint, 0, test)

			s := o.ScanLIFO(ra)
			for i := 0; s.Next(); i++ {
				lifo = append(lifo, s.Value())
			}

			s = o.ScanFIFO(ra)
			for i := 0; s.Next(); i++ {
				fifo = append(fifo, s.Value())
			}
			assert.Equal(t, len(lifo), len(fifo))
			assert.Equal(t, fifo[0], lifo[len(lifo)-1])
			assert.Equal(t, lifo[0], fifo[len(fifo)-1])
			// check contiguity:
			last := lifo[0]
			for nth, i := range lifo[1:] {
				require.Equal(t, ra.Mask(last+1), i, "at %d", nth)
				last = i
			}

			last = fifo[0]
			for nth, i := range fifo[1:] {
				require.Equal(t, ra.Mask(last+test-1), i, "at %d", nth)
				last = i
			}
		})
	}
}
