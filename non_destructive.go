package o

// All non-destructively returns all indexes, in order, for the given
// ring accountant.
func All(ring RingAccountant) []uint {
	r := make([]uint, ring.Size())
	elt := ring.start()
	for i := range r {
		r[i] = elt
		elt = ring.mask(elt + 1)
	}
	return r
}

// All non-destructively returns all indexes, in reverse order, for
// the given ring accountant.
func Rev(ring RingAccountant) []uint {
	r := make([]uint, ring.Size())
	elt := ring.start()
	for i := range r {
		r[len(r)-i-1] = elt
		elt = ring.mask(elt + 1)
	}
	return r
}

func Start1(ring RingAccountant) uint {
	return ring.start()
}

func End1(ring RingAccountant) uint {
	cap := ring.capacity()
	start := ring.start()
	size := ring.Size()
	if start+size > cap {
		return cap
	} else {
		return start + size
	}
}

func End2(ring RingAccountant) uint {
	cap := ring.capacity()
	start := ring.start()
	size := ring.Size()
	if start+size > cap {
		return start
	} else {
		return 0
	}
}
