package o_test

import (
	"testing"

	"github.com/antifuchs/o"
	"github.com/stretchr/testify/assert"
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
				test.ra.ForcePush()
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
				test.ra.ForcePush()
			}
			s := o.ScanFIFO(test.ra)
			for s.Next() {
				results = append(results, s.Value())
			}
			assert.Equal(t, test.expected, results)
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
				test.ra.ForcePush()
			}
			before := test.ra.Size()
			first, second := test.ra.Inspect()
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
				test.ra.ForcePush()
			}
			first, second := test.ra.Consume()
			t.Logf("%#v", test.ra)
			assert.Equal(t, test.first, first, "first")
			assert.Equal(t, test.second, second, "second")
			assert.Equal(t, uint(0), test.ra.Size())
		})
	}
}

func TestReserve(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		cap           uint
		fill          int
		read          int
		add           uint
		first, second o.Range
		err           error
	}{
		{
			name:  "basic5/13",
			cap:   5,
			add:   13,
			first: o.Range{0, 5}, second: o.Range{0, 0},
			err: o.ErrFull,
		},
		{
			name:  "mask4/13",
			cap:   4,
			add:   13,
			first: o.Range{0, 4}, second: o.Range{0, 0},
			err: o.ErrFull,
		},
		{
			name:  "zero",
			cap:   4,
			add:   0,
			first: o.Range{0, 0}, second: o.Range{0, 0},
		},
		{
			name:  "centered",
			cap:   5,
			fill:  4,
			read:  2,
			add:   13,
			first: o.Range{4, 5}, second: o.Range{0, 2},
			err: o.ErrFull,
		},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ring := o.NewRing(test.cap)
			for i := 0; i < test.fill; i++ {
				ring.ForcePush()
			}
			for i := 0; i < test.read; i++ {
				ring.Shift()
			}

			first, second, err := o.Reserve(ring, test.add)
			t.Log("Reserve:", first, second, err)
			assert.Equal(t, test.first, first, "first")
			assert.Equal(t, test.second, second, "second")
			assert.Equal(t, test.err, err)
		})
	}
}
