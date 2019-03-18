package o_test

import (
	"fmt"
	"testing"

	"github.com/antifuchs/o"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestPropShiftPushes(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)
	properties.Property("Read own writes", prop.ForAll(
		func(ringSize, entries uint) bool {
			ring := o.NewRing(ringSize)
			for i := uint(0); i < entries; i++ {
				pushed, _ := ring.Push()
				shifted, _ := ring.Shift()

				if pushed != shifted {
					return false
				}
			}
			return true
		},
		gen.UInt().WithLabel("ring size"),
		gen.UIntRange(1, 257*90).WithLabel("number of entries made"),
	))
	properties.TestingRun(t)
}

func TestPropBounds(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)

	properties.Property("Read own writes", prop.ForAll(
		func(ringSize, entries uint) string {
			ring := o.NewRing(ringSize)

			wFirst, wSecond, err := ring.PushN(entries)
			if entries > ringSize {
				if err == nil {
					return "should have errored"
				}
				return ""
			}

			if entries == 0 {
				if !ring.Empty() {
					return "should be empty"
				}
				return ""
			}

			if entries == ringSize && !ring.Full() {
				return "should be full"
			}

			if entries < ringSize && ring.Full() {
				return "should not be full"
			}

			first, second := ring.Consume()
			if (wFirst.Start != first.Start && first.End != wFirst.End+1) ||
				(!second.Empty() && wSecond.Start != second.Start && second.End != wSecond.End+1) {
				return fmt.Sprintf("Expected same ranges, but\n%#v %#v\n%#v %#v",
					wFirst, wSecond, first, second)
			}

			return ""
		},
		gen.UInt().WithLabel("ring size"),
		gen.UIntRange(0, 257*90).WithLabel("number of entries made"),
	))
	properties.TestingRun(t)
}
