package o

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
