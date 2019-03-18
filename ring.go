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
type Ring struct {
	ringBackend
}

type ringBackend interface {
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

	// end returns the index of the last element that can be read.
	end() uint

	// capacity returns the number of elements that the ring
	// accounts for.
	capacity() uint

	// reset adjusts the difference between the read and write
	// points of the ring back to 0.
	reset()

	// add accounts for n new elements in the ring. If fewer
	// elements could be accounted for, only accounts for the ones
	// that could fit and returns ErrFull.
	add(n uint) (uint, error)
}

// ForcePush forces a new element onto the ring, discarding the oldest
// element if the ring is full. It returns the index of the inserted
// element.
func (r Ring) ForcePush() uint {
	if r.Full() {
		_, _ = r.Shift()
	}
	i, _ := r.Push()
	return i
}

// Returns a new Ring data structure with the given capacity. If cap
// is a power of 2, returns a data structure that is optimized for
// modulo-2 accesses. Otherwise, the returned data structure uses
// general modulo division for its integer math.
func NewRing(cap uint) Ring {
	if bits.OnesCount(cap) == 1 {
		return Ring{&maskRing{cap: cap}}
	}
	return Ring{&basicRing{cap: cap}}
}

// A type, usually a collection, that has length. This is inspired by
// (but kept intentionally smaller than) sort.Interface.
type Slice interface {
	// Len returns the length of a slice.
	Len() int
}

// NewRingForSlice creates a Ring that fits a slice. The slice's type
// must implement o.Slice (which is satisfied if the type implements
// sort.Interface, also).
//
// It is not advisable to resize the slice after creating a ring for
// it.
func NewRingForSlice(i Slice) Ring {
	return NewRing(uint(i.Len()))
}
