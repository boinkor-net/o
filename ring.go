package o

import "math/bits"

type fullErr uint

func (e fullErr) Error() string {
	return "inserting into a full ring"
}

type emptyErr uint

func (e emptyErr) Error() string {
	return "reading from an empty ring"
}

// ErrEmpty indicates a removal operation on an empty ring.
const ErrEmpty emptyErr = iota

// ErrFull indicates an addition operation on a full ring.
const ErrFull fullErr = iota

// Ring provides accounting functions for ring buffers.
type Ring interface {
	// Push lets a writer account for a new element in the ring,
	// and returns that element's index.
	//
	// Returns ErrFull if the ring is filled to capacity.
	Push() (uint, error)

	// Shift lets a reader account for removing an element from
	// the ring for reading, returning that element's index.
	//
	// Returns ErrEmpty if the ring has no elements to read.
	Shift() (uint, error)

	// Full returns whether the ring buffer is at capacity.
	Full() bool

	// Empty returns whether the ring has zero element in it.
	Empty() bool

	// Size returns the number of elements in the ring buffer.
	Size() uint

	// Mask adjusts an index value to fit the ring buffer.
	Mask(uint) uint

	// start returns the index of first element that can be read.
	start() uint

	// capacity returns the number of elements that the ring
	// accounts for.
	capacity() uint
}

// ForcePush forces a new element onto the ring, discarding the oldest
// element if the ring is full.
func ForcePush(r Ring) uint {
	if r.Full() {
		_, _ = r.Shift()
	}
	i, _ := r.Push()
	return i
}

// Returns a new Ring data structure. If cap is a power of 2, returns
// a data structure that is optimized for modulo-2
// accesses. Otherwise, the returned data structure uses general
// modulo division for its integer math.
func NewRing(cap uint) Ring {
	if bits.OnesCount(cap) == 1 {
		return &maskRing{cap: cap}
	}
	return &basicRing{cap: cap}
}
