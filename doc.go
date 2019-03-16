// Package o provides a ring-buffer accounting abstraction that allows
// you to build your own ring buffers without having to dynamically
// cast values from interface{} back and forth.
//
// Implementing your own
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
// Behavior when full
//
// The ring buffer accountants defined in this package all return
// errors if they're full or empty. To simply overwrite the oldest
// element, use function ForcePush.
//
// Non-destructive operations on Rings
//
// If your code needs to only inspect the contents of the ring instead
// of shifting them out, you can use a stateful iterator in either
// LIFO or FIFO direction available to do this conveniently:
//
//     o.ScanLIFO(ring) and
//     o.ScanFIFO(ring)
//
// See Scanner for defails and usage examples.
//
// Thread Safety
//
// None of the data structures provided here are safe from data
// races. To use them in thread-safe ring buffer implementations,
// users must protect both the accounting operations and your backing
// buffer writes with a Mutex.
//
// Credit
//
// The ring buffer accounting techniques in this package and were
// translated into go from a post on the blog of Juho Snellman,
// https://www.snellman.net/blog/archive/2016-12-13-ring-buffers/.
package o
