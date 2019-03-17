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
		gen.UInt().SuchThat(func(x uint) bool { return x > 0 }).WithLabel("ring size"),
		gen.UIntRange(1, 257*90).WithLabel("number of entries made"),
	),
	)
	properties.TestingRun(t)
}

func TestPropMatchingRanges(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)
	properties.Property("Ranges in scanners match", prop.ForAll(
		func(cap, overage uint) string {
			ring := o.NewRing(cap)
			insert := cap + overage
			for i := uint(0); i < insert; i++ {
				o.ForcePush(ring)
			}
			if ring.Size() != cap {
				return "Size does not match cap"
			}

			fifo := make([]uint, 0, cap)
			lifo := make([]uint, 0, cap)

			s := o.ScanLIFO(ring)
			for i := 0; s.Next(); i++ {
				lifo = append(lifo, s.Value())
			}

			s = o.ScanFIFO(ring)
			for i := 0; s.Next(); i++ {
				fifo = append(fifo, s.Value())
			}
			if len(lifo) != len(fifo) {
				return "Length mismatch between lifo&fifo order"
			}
			last := lifo[0]
			for nth, _ := range lifo {
				if lifo[nth] != fifo[len(fifo)-1-nth] {
					return fmt.Sprintf("fifo / lifo mismatch:\n%#v\n%#v", lifo, fifo)
				}
				if nth > 0 && ring.Mask(last+1) != lifo[nth] {
					return fmt.Sprintf("indexes not continuous: %#v", lifo)
				}
				last = lifo[nth]
			}
			return ""
		},
		gen.UIntRange(1, 2000).SuchThat(func(x uint) bool { return x > 0 }).WithLabel("ring size"),
		gen.UIntRange(1, 100).WithLabel("ring size"),
	))
	properties.TestingRun(t)
}
