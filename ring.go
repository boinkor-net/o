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

	// shiftN "reads" n continuous indexes from the ring and
	// returns the first and last (masked) index. If n is larger
	// than the ring's Size, returns zeroes and ErrEmpty.
	shiftN(n uint) (start uint, end uint, err error)
}

// Capacity returns the number of continuous indexes that can be
// represented on the ring.
func (r Ring) Capacity() uint {
	return r.capacity()
}

// Empty returns whether the ring has zero elements that are readable
// on it.
func (r Ring) Empty() bool {
	return r.empty()
}

// ForcePush forces a new element onto the ring, discarding the oldest
// element if the ring is full. It returns the index of the inserted
// element.
//
// Using ForcePush to insert into the Ring means the Ring will lose
// data that has not been consumed yet. This is fine under some
// circumstances, but can have disastrous consequences for code that
// expects to read consistent data. It is generally safer to use .Push
// and handle ErrFull explicitly.
func (r Ring) ForcePush() uint {
	if r.full() {
		_, _ = r.Shift()
	}
	i, _ := r.Push()
	return i
}

// Full returns true if the Ring has occupied all possible index
// values.
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

// NewRing returns a new Ring data structure with the given
// capacity.
//
// If cap is a power of 2, returns a Ring optimized for
// bitwise manipulation of the indexes.
//
// If cap is 0, returns a Ring that does not perform any operations
// and only returns errors.
//
// Otherwise, the returned data structure uses general modulo division
// for its integer adjustments, and is a lot slower than the
// power-of-2 variant.
func NewRing(cap uint) Ring {
	if cap == 0 {
		return Ring{zeroRing{}}
	}
	if bits.OnesCount(cap) == 1 {
		return Ring{&maskRing{cap: cap}}
	}
	return Ring{&basicRing{cap: cap}}
}

// Slice represents a type, usually a collection, that has
// length. This is inspired by (but kept intentionally smaller than)
// sort.Interface.
type Slice interface {
	// Len returns the length of a slice.
	Len() int
}

// NewRingForSlice creates a Ring that fits a slice. The slice's type
// must implement o.Slice (which is also satisfied if the type
// implements sort.Interface).
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
	start, _, err := r.shiftN(1)
	return start, err
}

// Size returns the number of elements in the ring buffer.
func (r Ring) Size() uint {
	return r.size()
}
