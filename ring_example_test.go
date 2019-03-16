package o_test

import (
	"fmt"

	"github.com/antifuchs/o"
)

// A simple queue of strings that refuses to add elements when full.
func ExampleRing() {
	queueIndexes := o.NewRing(16)
	for i := 0; i < 16; i++ {
		next, _ := queueIndexes.Push()
		fmt.Print(next, " ")
	}
	_, err := queueIndexes.Push()
	fmt.Print(err)
	// Output:
	// 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 inserting into a full ring
}
