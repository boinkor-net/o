package o

// Range is a normalized set of numbers representing continuous range
// of indexes that are occupied in the Ring. Start must always be <=
// End.
//
// They can be used in go slice bounds like so:
//
//     [range.Start:range.End]
type Range struct {
	Start uint // The first element of the range
	End   uint // The first element that is not part of the range.
}

// Empty is true if the range does not contain any indexes.
func (r Range) Empty() bool {
	return r.Start == r.End
}

// Length returns the number of elements in the range.
func (r Range) Length() uint {
	return r.End - r.Start
}

// Inspect returns a set of indexes that represent the bounds of the
// elements occupied in the ring.
//
// Returned indexes
//
// Since a ring buffer consists of indexes that might wrap around to
// zero, callers of Inspect must use all returned Ranges to get an
// accurate picture of the occupied elements. The second range may be
// empty (Start & Length = 0) if there is nothing occupied on the left
// part of the buffer.
func Inspect(ring Ring) (first Range, second Range) {
	if ring.Empty() {
		return
	}
	first.Start = ring.start()
	end1 := ring.end()

	first.End = end1 + 1
	if end1 <= first.Start {
		second.End = end1
		first.End = ring.capacity()
	}
	return
}

// Consume resets the ring to its empty state, returning a set of
// indexes that can be used to construct a copy of the elements that
// were occupied in the ring prior to resetting.
//
// See also Inspect.
func Consume(ring Ring) (first Range, second Range) {
	defer ring.reset()
	return Inspect(ring)
}

// Scanner implements iterating over the elements in a Ring without
// removing them. A scanner can go in either LIFO (oldest element
// first) or FIFO (newest element first) direction.
type Scanner struct {
	ring Ring
	cur  uint
	fifo bool
}

// ScanLIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in LIFO (oldest to newest) direction.
func ScanLIFO(ring Ring) *Scanner {
	return &Scanner{ring, ring.start(), false}
}

// ScanFIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in FIFO (newest to oldest) direction.
func ScanFIFO(ring Ring) *Scanner {
	return &Scanner{ring, ring.capacity()*2 + ring.start(), true}
}

// Next advances the Scanner to the next available element. If no next
// element is available, it returns false.
func (s *Scanner) Next() bool {
	var next uint
	var ok bool
	if s.fifo {
		next = s.cur - 1
		ok = next > s.ring.start()+s.ring.capacity()-1
	} else {
		next = s.cur + 1
		ok = next <= s.ring.start()+s.ring.Size()
	}
	if !ok {
		return false
	}
	s.cur = next
	return true
}

// Value returns the index value of the Scanner's current position.
func (s *Scanner) Value() uint {
	if s.fifo {
		return s.ring.Mask(s.cur)
	}
	return s.ring.Mask(s.cur - 1)
}
