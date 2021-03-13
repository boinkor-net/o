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

			last := fifo[0]
			for nth := range fifo {
				if fifo[nth] != lifo[len(lifo)-1-nth] {
					return fmt.Sprintf("lifo / fifo mismatch:\n%#v\n%#v", fifo, lifo)
				}
				if nth > 0 && ring.Mask(last+1) != fifo[nth] {
					return fmt.Sprintf("indexes not continuous: %#v", fifo)
				}
				last = fifo[nth]
			}
			return ""
		},
		gen.UIntRange(0, 2000).WithLabel("ring size"),
		gen.UIntRange(1, 100).WithLabel("overflow"),
	))
	properties.TestingRun(t)
}

func TestPropPushN(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(params)
	properties.Property("Pushing N elements", prop.ForAll(
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

func TestPropShiftN(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 10000
	properties := gopter.NewProperties(params)
	properties.Property("Shifting N elements", prop.ForAll(
		func(cap, fill, skip, read uint) string {
			ring := o.NewRing(cap)

			for i := uint(0); i < fill; i++ {
				ring.ForcePush()
			}
			var startIdx uint
			if fill > cap && cap > 0 {
				startIdx = ring.Mask(fill - cap)
			}
			for i := uint(0); i < skip; i++ {
				idx, err := ring.Shift()
				if err == nil {
					startIdx = ring.Mask(idx + 1)
				}
			}
			startSize := ring.Size()

			first, second, err := ring.ShiftN(read)
			overflows := read > cap || read > startSize
			readAny := !first.Empty() || !second.Empty()

			if overflows && err == nil {
				return "expected error"
			}
			if !overflows && err != nil {
				return "unexpected error"
			}

			if overflows && readAny {
				return fmt.Sprintf("would overflow, but read %d elements:\n%#v %#v",
					first.Length()+second.Length(), first, second)
			}
			if !overflows && first.Length()+second.Length() != read {
				return fmt.Sprintf("did not read %d elements:\n%#v %#v",
					read, first, second)
			}
			if readAny && startIdx != first.Start {
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
			if !second.Empty() && !overflows && second.End != read-first.Length() {
				return fmt.Sprintf("bad end bound on second range: %d expected, but %#v %#v",
					read-first.Length(), first, second)
			}
			if !second.Empty() && overflows && second.End != cap-startSize-first.Length() {
				return fmt.Sprintf("bad end bound on overflowing second range: %d expected, but %#v %#v",
					cap-startSize-first.Length(), first, second)
			}
			return ""
		},
		gen.UIntRange(0, 2000).WithLabel("ring size"),
		gen.UIntRange(0, 100).WithLabel("elements to fill in"),
		gen.UIntRange(0, 100).WithLabel("elements to skip before reading"),
		gen.UIntRange(0, 100).WithLabel("elements to read"),
	))
	properties.TestingRun(t)
}
