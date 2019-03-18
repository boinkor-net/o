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

	// Now find all the indexes in last-in/first-out order:
	s := o.ScanFIFO(ring)
	for s.Next() {
		fmt.Print(s.Value(), " ")
	}
	// Output:
	// 1 0 16 15 14 13 12 11 10 9 8 7 6 5 4 3 2
}
