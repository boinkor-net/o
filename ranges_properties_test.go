package o_test

import (
	"fmt"
	"testing"

	"github.com/antifuchs/o"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestPropFIFOandLIFOMatch(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 1000
	properties := gopter.NewProperties(params)
	properties.Property("Ranges in scanners match", prop.ForAll(
		func(cap, overage uint) string {
			ring := o.NewRing(cap)
			insert := cap + overage
			for i := uint(0); i < insert; i++ {
				ring.ForcePush()
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
			if len(lifo) == 0 {
				// nothing else to check
				return ""
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
		gen.UIntRange(0, 2000).WithLabel("ring size"),
		gen.UIntRange(1, 100).WithLabel("overflow"),
	))
	properties.TestingRun(t)
}

func TestPropReserve(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(params)
	properties.Property("Ranges in scanners match", prop.ForAll(
		func(cap, fill, read, reserve uint) string {
			ring := o.NewRing(cap)
			var startIdx uint
			for i := uint(0); i < fill; i++ {
				startIdx = ring.Mask(ring.ForcePush() + 1)
			}
			for i := uint(0); i < read; i++ {
				ring.Shift()
			}
			startSize := ring.Size()
			overflows := startSize+reserve > cap

			first, second, err := ring.PushN(reserve)
			reservedAny := !first.Empty() || !second.Empty()
			if overflows && err == nil {
				return "expected error"
			}
			if !overflows && err != nil {
				return "unexpected error"
			}

			if overflows && reservedAny {
				return fmt.Sprintf("would overflow, but reserved %d elements:\n%#v %#v",
					first.Length()+second.Length(), first, second)
			}
			if !overflows && first.Length()+second.Length() != reserve {
				return fmt.Sprintf("did not reserve %d elements:\n%#v %#v",
					reserve, first, second)
			}
			if reservedAny && startIdx != first.Start {
				return fmt.Sprintf("expected reservation to start at %d, but %#v",
					startIdx, first)
			}
			if !second.Empty() && first.End != cap {
				return fmt.Sprintf("bad end bound on first range: %d expected, but %#v",
					cap, first)
			}
			if !second.Empty() && second.Start != 0 {
				return fmt.Sprintf("bad start bound on second range: 0 expected, but %#v",
					second)
			}
			if !second.Empty() && !overflows && second.End != reserve-first.Length() {
				return fmt.Sprintf("bad end bound on second range: %d expected, but %#v %#v",
					reserve-first.Length(), first, second)
			}
			if !second.Empty() && overflows && second.End != cap-startSize-first.Length() {
				return fmt.Sprintf("bad end bound on overflowing second range: %d expected, but %#v %#v",
					cap-startSize-first.Length(), first, second)
			}
			return ""
		},
		gen.UIntRange(0, 2000).WithLabel("ring size"),
		gen.UIntRange(0, 100).WithLabel("elements to fill in"),
		gen.UIntRange(0, 100).WithLabel("elements to read"),
		gen.UIntRange(0, 100).WithLabel("elements to reserve"),
	))
	properties.TestingRun(t)
}
