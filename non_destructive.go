package o

// All returns all indexes occupied in the ring buffer in order (from
// oldest to youngest). It does not modify the ring.
func All(ring Ring) []uint {
	r := make([]uint, ring.Size())
	elt := ring.start()
	for i := range r {
		r[i] = elt
		elt = ring.mask(elt + 1)
	}
	return r
}

// Rev returns all indexes occupied in the ring buffer, in reverse
// order (from youngest to oldest). It does not modify the ring.
func Rev(ring Ring) []uint {
	r := make([]uint, ring.Size())
	elt := ring.start()
	for i := range r {
		r[len(r)-i-1] = elt
		elt = ring.mask(elt + 1)
	}
	return r
}

// Start1 returns the index of the first occupied entry in the ring
// buffer, to aid in iterating over all indexes in the ring.
func Start1(ring Ring) uint {
	return ring.start()
}

// End1 returns the end index of the first loop when iterating over
// all occupied indexes in the ring buffer. See Start1.
func End1(ring Ring) uint {
	cap := ring.capacity()
	start := ring.start()
	size := ring.Size()
	if start+size > cap {
		return cap
	} else {
		return start + size
	}
}

// End2 returns the end index of the second loop when iterating over
// all occupied indexes in the ring buffer. See Start1.
func End2(ring Ring) uint {
	cap := ring.capacity()
	start := ring.start()
	size := ring.Size()
	if start+size > cap {
		return start
	} else {
		return 0
	}
}
