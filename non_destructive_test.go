package o

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ra       Ring
		cycles   uint
		expected []uint
	}{
		{"basic5/13", NewRing(5), 13, []uint{3, 4, 0, 1, 2}},
		{"basic5/6", NewRing(5), 6, []uint{1, 2, 3, 4, 0}},
		{"mask4/13", NewRing(4), 13, []uint{1, 2, 3, 0}},
		{"mask4/6", NewRing(4), 6, []uint{2, 3, 0, 1}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i uint
			for i = 0; i < test.cycles; i++ {
				ForcePush(test.ra)
			}
			assert.Equal(t, test.expected, All(test.ra))
		})
	}
}

func TestRev(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		ra       Ring
		cycles   uint
		expected []uint
	}{
		{"basic5/13", NewRing(5), 13, []uint{2, 1, 0, 4, 3}},
		{"basic5/6", NewRing(5), 6, []uint{0, 4, 3, 2, 1}},
		{"mask4/13", NewRing(4), 13, []uint{0, 3, 2, 1}},
		{"mask4/6", NewRing(4), 6, []uint{1, 0, 3, 2}},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i uint
			for i = 0; i < test.cycles; i++ {
				ForcePush(test.ra)
			}
			assert.Equal(t, test.expected, Rev(test.ra))
		})
	}
}

func TestStartEnd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		ra                 Ring
		cycles             uint
		start1, end1, end2 uint
	}{
		// filled beyond their capacity:
		{"basic5/13", NewRing(5), 13, 3, 5, 3},
		{"basic5/6", NewRing(5), 6, 1, 5, 1},
		{"mask4/13", NewRing(4), 13, 1, 4, 1},
		{"mask4/6", NewRing(4), 6, 2, 4, 2},
		// Filled to less than capacity:
		{"mask4/2", NewRing(4), 2, 0, 2, 0},
		{"basic5/2", NewRing(5), 2, 0, 2, 0},
	}
	for _, elt := range tests {
		test := elt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var i uint
			for i = 0; i < test.cycles; i++ {
				ForcePush(test.ra)
			}
			t.Logf("e: %d:%d, 0:%d", test.start1, test.end1, test.end2)
			t.Logf("g: %d:%d, 0:%d - %d:%d", Start1(test.ra), End1(test.ra), End2(test.ra),
				test.ra.start(), test.ra.capacity())
			assert.Equal(t, test.start1, Start1(test.ra))
			assert.Equal(t, test.end1, End1(test.ra))
			assert.Equal(t, test.end2, End2(test.ra))
		})
	}
}

func TestMatching(t *testing.T) {
	var toinsert = 263 // prime & greater than the max. capacity
	for i := uint(1); i < 256; i++ {
		test := i
		t.Run(fmt.Sprintf("%03d", test), func(t *testing.T) {
			t.Parallel()
			ra := NewRing(test)
			for i := 0; i < toinsert; i++ {
				ForcePush(ra)
			}
			assert.Equal(t, ra.Size(), test)

			t.Log(Start1(ra), End1(ra), End2(ra))

			all := All(ra)
			t.Log("Checking All()=", all)
			assert.Equal(t, uint(len(all)), test)
			assert.Equal(t, all[0], Start1(ra),
				"start index should be the same for All() and Start1()")
			assert.Equal(t, uint(len(all)), End1(ra),
				"End1() should always stop at the end of the array")
			assert.Equal(t, Start1(ra), End2(ra),
				"End2() should always return the index of Start1()")
			assert.Equal(t, (all[len(all)-1]+1)%uint(len(all)), End2(ra),
				"End2() should be one past the last index of All()")

			rev := Rev(ra)
			// The loop for going in reverse looks like:
			//     for i := ring.Mask(o.Start1(ring)-1); i > 0; i-- {}
			//     for i := o.End1(ring)-1; i >= o.Start1(ring); i-- {}
			t.Log("Checking Rev()=", rev)
			assert.Equal(t, uint(len(rev)), test)
			assert.Equal(t, rev[0], ra.Mask(Start1(ra)-1),
				"start index should be the same for Rev() and the one before Start1()")
			assert.Equal(t, rev[len(rev)-1], Start1(ra),
				"The last entry in Rev() should be Start1()")
		})
	}
}
