package o

// BasicRing contains the accounting data for a ring buffer or other
// data structure of arbitrary length. It uses three variables (insert
// index, length of buffer, ring capacity) to keep track of the
// state.
//
// The index wrap-around operation is implemented with modulo division.
type basicRing struct {
	cap, read, length uint
}

func (r *basicRing) mask(val uint) uint {
	return val % r.cap
}

func (r *basicRing) start() uint {
	return r.read
}

func (r *basicRing) end() uint {
	return r.mask(r.read + r.length)
}

func (r *basicRing) capacity() uint {
	return r.cap
}

func (r *basicRing) reset() {
	r.length = 0
}

func (r *basicRing) pushN(n uint) (uint, uint, error) {
	start := r.length
	if n > r.cap-r.length {
		idx := r.mask(r.read + start)
		return idx, idx, ErrFull
	}
	r.length += n
	return r.mask(r.read + start), r.mask(r.read + r.length), nil
}

func (r *basicRing) shiftN(n uint) (uint, uint, error) {
	start := r.read
	if n > r.size() {
		return start, start, ErrEmpty
	}
	r.length -= n
	//i := r.read
	r.read = r.mask(r.read + n)
	return start, r.read, nil
}

func (r *basicRing) full() bool {
	return r.cap == r.length
}

func (r *basicRing) empty() bool {
	return r.length == 0
}

func (r *basicRing) size() uint {
	return r.length
}

var _ ringBackend = &basicRing{}
