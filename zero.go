package o

// Implements the ring algorithms for rings of size 0. These are
// special because we really would like to avoid division by zero.
type zeroRing struct{}

func (z zeroRing) shiftN(uint) (uint, uint, error) {
	return 0, 0, ErrEmpty
}

func (z zeroRing) full() bool {
	return true
}

func (z zeroRing) empty() bool {
	return false
}

func (z zeroRing) size() uint {
	return 0
}

func (z zeroRing) mask(uint) uint {
	return 0
}

func (z zeroRing) start() uint {
	return 0
}

func (z zeroRing) end() uint {
	return 0
}

func (z zeroRing) capacity() uint {
	return 0
}

func (z zeroRing) reset() {}

func (z zeroRing) pushN(n uint) (start uint, end uint, err error) {
	return 0, 0, ErrFull
}

var _ ringBackend = zeroRing{}
