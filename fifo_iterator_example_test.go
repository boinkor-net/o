package o_test

import (
	"fmt"

	"github.com/antifuchs/o"
)

func ExampleScanFIFO() {
	ring := o.NewRing(17)
	// Put stuff in the ring:
	for i := 0; i < 19; i++ {
		ring.ForcePush()
	}

	// Now find all the indexes in first-in/first-out order:
	s := o.ScanFIFO(ring)
	for s.Next() {
		fmt.Print(s.Value(), " ")
	}
	// Output:
	// 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 0 1
}
