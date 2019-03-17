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

func TestInspect(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		ra            o.Ring
		cycles        uint
		first, second o.Range
	}{
		{"basic5/13", o.NewRing(5), 13, o.Range{3, 5}, o.Range{0, 3}},
		{"basic5/6", o.NewRing(5), 6, o.Range{1, 5}, o.Range{0, 1}},
		{"basic5/4", o.NewRing(5), 4, o.Range{0, 5}, o.Range{0, 0}},
		{"basic5/0", o.NewRing(5), 0, o.Range{0, 0}, o.Range{0, 0}},
		{"mask4/13", o.NewRing(4), 13, o.Range{1, 4}, o.Range{0, 1}},
		{"mask4/6", o.NewRing(4), 6, o.Range{2, 4}, o.Range{0, 2}},
		{"mask4/4", o.NewRing(4), 4, o.Range{0, 4}, o.Range{0, 0}},
		{"mask4/0", o.NewRing(4), 0, o.Range{0, 0}, o.Range{0, 0}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i uint
			for i = 0; i < test.cycles; i++ {
				o.ForcePush(test.ra)
			}
			before := test.ra.Size()
			first, second := o.Inspect(test.ra)
			t.Logf("%#v", test.ra)
			assert.Equal(t, test.first, first, "first")
			assert.Equal(t, test.second, second, "second")
			assert.Equal(t, before, test.ra.Size())
		})
	}
}

func TestConsume(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		ra            o.Ring
		cycles        uint
		first, second o.Range
	}{
		{"basic5/13", o.NewRing(5), 13, o.Range{3, 5}, o.Range{0, 3}},
		{"basic5/6", o.NewRing(5), 6, o.Range{1, 5}, o.Range{0, 1}},
		{"basic5/4", o.NewRing(5), 4, o.Range{0, 5}, o.Range{0, 0}},
		{"basic5/0", o.NewRing(5), 0, o.Range{0, 0}, o.Range{0, 0}},
		{"mask4/13", o.NewRing(4), 13, o.Range{1, 4}, o.Range{0, 1}},
		{"mask4/6", o.NewRing(4), 6, o.Range{2, 4}, o.Range{0, 2}},
		{"mask4/4", o.NewRing(4), 4, o.Range{0, 4}, o.Range{0, 0}},
		{"mask4/0", o.NewRing(4), 0, o.Range{0, 0}, o.Range{0, 0}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i uint
			for i = 0; i < test.cycles; i++ {
				o.ForcePush(test.ra)
			}
			first, second := o.Consume(test.ra)
			t.Logf("%#v", test.ra)
			assert.Equal(t, test.first, first, "first")
			assert.Equal(t, test.second, second, "second")
			assert.Equal(t, uint(0), test.ra.Size())
		})
	}
}
