package o

type maskRing struct {
	cap, read, write uint
}

func (r *maskRing) mask(val uint) uint {
	return val & (r.cap - 1)
}

func (r *maskRing) start() uint {
	return r.mask(r.read)
}

func (r *maskRing) reset() {
	r.read = r.write
}

func (r *maskRing) capacity() uint {
	return r.cap
}

func (r *maskRing) end() uint {
	return r.mask(r.write)
}

func (r *maskRing) pushN(n uint) (uint, uint, error) {
	start := r.write
	if n > r.cap-r.size() {
		i := r.mask(start)
		return i, i, ErrFull
	}
	r.write += n
	return r.mask(start), r.mask(r.write), nil
}

func (r *maskRing) full() bool {
	return r.size() == r.cap
}

func (r *maskRing) empty() bool {
	return r.read == r.write
}

func (r *maskRing) shift() (uint, error) {
	if r.empty() {
		return 0, ErrEmpty
	}
	i := r.read
	r.read++
	return r.mask(i), nil
}

func (r *maskRing) size() uint {
	return r.write - r.read
}

var _ ringBackend = &maskRing{}
