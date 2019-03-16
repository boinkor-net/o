package o_test

import (
	"fmt"
	"sort"

	"github.com/antifuchs/o"
)

func ExampleNewRingForSlice() {
	// create some backing store:
	store := make([]int, 90)

	// put a ring on it:
	ring := o.NewRingForSlice(sort.IntSlice(store))

	// it is empty:
	fmt.Println(ring.Empty())

	// Output:
	// true
}
