package o

// Range is a normalized set of numbers representing continuous range
// of indexes that are occupied in the Ring. Start must always be <=
// End.
//
// They can be used in go slice bounds like so:
//
//	[range.Start:range.End]
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
	if r.End <= r.Start {
		return 0
	}
	return r.End - r.Start
}

// Inspect returns a set of indexes that represent the bounds of the
// elements occupied in the ring.
//
// # Returned indexes
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

	first.End = end1
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
// nothing and returns ErrEmpty; the ranges returned in this case are
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

// indexes is a representation of walking on a Range. It has a first
// index and a last index that is valid, and a direction in which to
// traverse the Range.
//
// Note that "first" and "last" are both inclusive here: In contrast
// to a Range(and Ring)'s idea of Start / End (which is a half-open
// interval), this holds the first and last valid index in the range.
//
// An indexes struct tightens as you traverse it, and it always
// tightens from the beginning: The direction field gets added to
// first, and once first == last, the index traversal is complete and
// no more indexes will be produced.
type indexes struct {
	first, last uint
	direction   int
}

func (r Range) toFIFOTraversal() *indexes {
	if r.Empty() {
		return nil
	}
	return &indexes{
		direction: +1,
		first:     r.Start,
		last:      r.End - 1,
	}
}

func (r Range) toLIFOTraversal() *indexes {
	if r.Empty() {
		return nil
	}
	return &indexes{
		direction: -1,
		first:     r.End - 1,
		last:      r.Start,
	}
}

func (i *indexes) hasNext() bool {
	return i != nil && int(i.first)*i.direction <= int(i.last)*i.direction
}

func (i *indexes) next() uint {
	cur := i.first
	i.first += uint(i.direction)
	return cur
}

// Scanner implements iterating over the elements in a Ring
// without removing them. It represents a snapshot of the Ring at the
// time it was created.
//
// A Scanner does not update its Ring's range validity when .Next is
// called. Adding or reading elements from the Ring while a Scanner is
// active can mean invalidated indexes will be returned from the
// Scanner.
type Scanner struct {
	pos     *uint
	current *indexes
	second  *indexes
}

// ScanFIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in FIFO (newest to oldest) direction.
func ScanFIFO(ring Ring) *Scanner {
	first, second := ring.Inspect()
	return &Scanner{
		current: first.toFIFOTraversal(),
		second:  second.toFIFOTraversal(),
	}
}

// ScanLIFO returns a Scanner for the given Ring that iterates over
// the occupied indexes in LIFO (newest to oldest, think of a stack)
// direction.
func ScanLIFO(ring Ring) *Scanner {
	first, second := ring.Inspect()
	return &Scanner{
		current: second.toLIFOTraversal(),
		second:  first.toLIFOTraversal(),
	}
}

// Next advances the Scanner in the traversal direction (forward in
// FIFO direction, backward in LIFO), returning a boolean indicating
// whether there *is* a next position in the ring.
//
// It is safe to call Next after reaching the last position - in that
// case, it will always return a negative result.
func (s *Scanner) Next() bool {
	s.pos = nil
	ok := s.current.hasNext()
	if ok {
		pos := s.current.next()
		s.pos = &pos
		return true
	}
	// We've exhausted one pool of indexes, pick the next one (or
	// none, there may be no next one):
	s.current = s.second
	ok = s.current.hasNext()
	if !ok {
		return false
	}
	pos := s.current.next()
	s.pos = &pos
	return true
}

// Value returns the next position in the traversal of a Ring's
// occupied positions, in the given order, after Next() returned a
// positive value.
//
// If Value is called before the first call to Next, or after Next
// returned a result indicating there are no more positions, Value
// panics.
func (s *Scanner) Value() uint {
	if s.pos == nil {
		panic("Value called when we know about no valid positions.")
	}
	return *s.pos
}
