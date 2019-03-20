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
func (r Ring) Inspect() (first Range, second Range) {
	if r.Empty() {
		return
	}
	first.Start = r.start()
	end1 := r.end()

	first.End = end1 + 1
	if end1 <= first.Start {
		second.End = end1
		first.End = r.capacity()
	}
	return
}

// Consume resets the ring to its empty state, returning a set of
// indexes that can be used to construct a copy of the elements that
// were occupied in the ring prior to resetting.
//
// See also Inspect.
func (r Ring) Consume() (first Range, second Range) {
	defer r.reset()
	return r.Inspect()
}

// PushN bulk-pushes count indexes onto the end of the Ring and
// returns ranges covering the indexes that were pushed.
//
// If the Ring can not accommodate all elements before filling up,
// PushN reserves nothing and returns ErrFull; the ranges returned in
// this case are meaningless and have zero length.
func (r Ring) PushN(count uint) (first, second Range, err error) {
	if count == 0 {
		return
	}
	first.Start = r.end()

	first.Start, first.End, err = r.pushN(count)
	if err != nil {
		return
	}
	if first.End <= first.Start && count > 0 {
		second.End = first.End
		first.End = r.capacity()
	}
	return
}

// ShiftN bulk-"read"s count indexes from the start of the Ring and
// returns ranges covering the indexes that were removed.
//
// If the Ring holds only fewer elements as requested, ShiftN reads
// nothing and returns ErrFull; the ranges returned in this case are
// meaningless and have zero length.
func (r Ring) ShiftN(count uint) (first, second Range, err error) {
	if count == 0 {
		return
	}
	first.Start, first.End, err = r.shiftN(count)
	if err != nil {
		return
	}
	if first.End <= first.Start && count > 0 {
		second.End = first.End
		first.End = r.capacity()
	}
	return
}

// Scanner implements iterating over the elements in a Ring without
// removing them. It represents a snapshot of the Ring at the time it
// was created. A scanner can go in either LIFO (oldest element first)
// or FIFO (newest element first) direction.
//
// A Scanner does not update its Ring's range validity when .Next is
// called. Adding or reading elements from the Ring while a Scanner is
// active can mean invalidated indexes will be returned from the
// Scanner.
type Scanner struct {
	cur    uint
	ranges []Range
	fifo   bool
}

// ScanLIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in LIFO (oldest to newest) direction.
func ScanLIFO(ring Ring) *Scanner {
	first, second := ring.Inspect()
	return &Scanner{first.Start, []Range{first, second}, false}
}

// ScanFIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in FIFO (newest to oldest) direction.
func ScanFIFO(ring Ring) *Scanner {
	first, second := ring.Inspect()
	return &Scanner{second.End, []Range{second, first}, true}
}

// Next advances the Scanner to the next available element. If no next
// element is available, it returns false.
func (s *Scanner) Next() bool {
	rg := &s.ranges[0]
	if rg.Empty() {
		s.ranges = s.ranges[1:]
		if len(s.ranges) == 0 {
			return false
		}
		rg = &s.ranges[0]
		if rg.Empty() {
			return false
		}
	}
	if s.fifo {
		s.cur = rg.End - 1
		rg.End--
	} else {
		s.cur = rg.Start
		rg.Start++
	}
	return true
}

// Value returns the index value of the Scanner's current position.
func (s *Scanner) Value() uint {
	return s.cur
}
