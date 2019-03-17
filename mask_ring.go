package o

type maskRing struct {
	cap, read, write uint
}

func (r *maskRing) Mask(val uint) uint {
	return val & (r.cap - 1)
}

func (r *maskRing) start() uint {
	return r.Mask(r.read)
}

func (r *maskRing) reset() {
	r.read = r.write
}

func (r *maskRing) capacity() uint {
	return r.cap
}

func (r *maskRing) end() uint {
	return r.Mask(r.write)
}

func (r *maskRing) add(n uint) (uint, error) {
	space := r.cap - r.Size()
	if n > space {
		r.write += space
		return space, ErrFull
	}
	r.write += n
	return n, nil
}

func (r *maskRing) Full() bool {
	return r.Size() == r.cap
}

func (r *maskRing) Empty() bool {
	return r.read == r.write
}

func (r *maskRing) Push() (uint, error) {
	if r.Full() {
		return 0, ErrFull
	}
	i := r.write
	r.write++

	return r.Mask(i), nil
}

func (r *maskRing) Shift() (uint, error) {
	if r.Empty() {
		return 0, ErrEmpty
	}
	i := r.read
	r.read++
	return r.Mask(i), nil
}

func (r *maskRing) Size() uint {
	return r.write - r.read
}

var _ Ring = &maskRing{}
