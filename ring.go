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

// Defines the functions implementations of the accountancy algorithms
// need to provide.
type ringBackend interface {
	shift() (uint, error)

	full() bool

	empty() bool

	size() uint

	mask(uint) uint

	// start returns the index of first element that can be read.
	start() uint

	// end returns the index of the last element that can be read.
	end() uint

	capacity() uint

	// reset adjusts the difference between the read and write
	// points of the ring back to 0.
	reset()

	// pushN accounts for n new elements in the ring and returns
	// the indexes of the first and last element. If not all
	// elements can be inserted, does not push them and returns
	// only ErrNotFound.
	pushN(n uint) (start uint, end uint, err error)
}

// Capacity returns the number of continuous indexes that can be
// represented on the ring. IOW, it returns the highest possible
// index+1.
func (r Ring) Capacity() uint {
	return r.capacity()
}

// Empty returns whether the ring has zero element in it.
func (r Ring) Empty() bool {
	return r.empty()
}

// ForcePush forces a new element onto the ring, discarding the oldest
// element if the ring is full. It returns the index of the inserted
// element.
func (r Ring) ForcePush() uint {
	if r.full() {
		_, _ = r.Shift()
	}
	i, _ := r.Push()
	return i
}

// Full returns whether the ring buffer is at capacity.
func (r Ring) Full() bool {
	return r.full()
}

// Mask adjusts an index value (which potentially exceeds the ring
// buffer's Capacity) to fit the ring buffer and returns the adjusted
// value.
//
// This method is probably most useful in tests, or when doing
// low-level things not supported by o.Ring yet. If you find yourself
// relying on this in code, please file a bug.
func (r Ring) Mask(i uint) uint {
	return r.mask(i)
}

// Returns a new Ring data structure with the given capacity. If cap
// is a power of 2, returns a data structure that is optimized for
// modulo-2 accesses. Otherwise, the returned data structure uses
// general modulo division for its integer math.
func NewRing(cap uint) Ring {
	if cap == 0 {
		return Ring{zeroRing{}}
	}
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

// Push lets a writer account for a new element in the ring,
// and returns that element's index.
//
// Returns ErrFull if the ring is filled to capacity.
func (r Ring) Push() (uint, error) {
	start, _, err := r.pushN(1)
	return start, err
}

// Shift lets a reader account for removing an element from
// the ring for reading, returning that element's index.
//
// Returns ErrEmpty if the ring has no elements to read.
func (r Ring) Shift() (uint, error) {
	return r.shift()
}

// Size returns the number of elements in the ring buffer.
func (r Ring) Size() uint {
	return r.size()
}
