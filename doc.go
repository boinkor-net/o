// Package o provides a ring-buffer accounting abstraction that allows
// you to build your own ring buffers without having to dynamically
// cast values from interface{} back and forth.
//
// # Implementing your own
//
// The trick to having a type-safe ring buffer is simple: You define a
// new structure that contains the accounting interface (defined here)
// and a buffer of the appropriate capacity on it. The accounting
// interface gives your code the indexes that it needs to update the
// ring at, and your code can go its merry way.
//
// For an example, see the ring buffer backed ReadWriter defined in
// package ringio.
//
// # Behavior when full
//
// The ring buffer accountants defined in this package all return
// errors if they're full or empty. To simply overwrite the oldest
// element, use function ForcePush.
//
// # Non-destructive operations on Rings
//
// If your code needs to only inspect the contents of the ring instead
// of shifting them out, you can use the Inspect method (which returns
// ranges, see the next section) or a stateful iterator in either LIFO
// or FIFO direction available to do this conveniently:
//
//	o.ScanLIFO(ring) and
//	o.ScanFIFO(ring)
//
// See Scanner for defails and usage examples.
//
// # Ranges across a Ring
//
// A Ring assumes that all indexes between the first occupied index
// (the "read" end) and the last occupied index (the "write" end) are
// continuously occupied. Since rings wrap around to zero, that
// doesn't mean however, that each continuous index is greater than
// the index before it
//
// To make it easier to deal with these index ranges, every operation
// that deals with ranges (e.g. Inspect, Consume, PushN) will return
// two Range objects, by convention named first and second.
//
// These ranges cover the two parts of the buffer. Assume a ring
// buffer like this, x marking occupied elements and _ marking
// unoccupied ones:
//
//	  0   1   2   3   4   5   6   7 (Capacity)
//	+---+---+---+---+---+---+---+
//	| x | _ | _ | x | x | x | x |
//	+---+---+---+---+---+---+---+
//	      ^       ^ Read end points here
//	      |
//	      +- Write end points here
//
// The way o.Ranges represents a range between the Read and the Write
// end is:
//
//	first  = o.Range{Start: 3, End: 7}
//	second = o.Range{Start: 0, End: 1}
//
// Note that the End index of a range is the first index greater than
// the one that's occupied. This allows using these Range ends as
// points in a slice expression without modification.
//
// # Thread Safety
//
// None of the data structures provided here are safe from data
// races. To use them in thread-safe ring buffer implementations,
// users must protect both the accounting operations and backing
// buffer writes with a Mutex.
//
// # Credit
//
// The ring buffer accounting techniques in this package and were
// translated into go from a post on the blog of Juho Snellman,
// https://www.snellman.net/blog/archive/2016-12-13-ring-buffers/.
package o
