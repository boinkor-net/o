package o

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

// RingAccountant specifies the methods that a ring buffer accounting
// data structure must have.
type RingAccountant interface {
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

	// Size returns the capacity the ring buffer.
	Size() uint

	// Mask adjusts an index value to fit the ring buffer. It is
	// used by various internal functions.
	Mask(uint) uint
}
