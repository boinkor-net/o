package o_test

import (
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
	),
	)
	properties.TestingRun(t)
}
