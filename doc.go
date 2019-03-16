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
// of shifting them out, there are two functions available to do this
// conveniently:
//
//     - All(ring) returns the occupied indexes from oldest to youngest
//     - Rev(ring) returns the occupied indexes from youngest to oldest
//
// If instead you prefer not to allocate an array the size of your
// buffer, you can also write two for loops with the functions
// Start1(), End1() and End2():
//
//     for i := o.Start1(ring); i<o.End1(ring); i++ {
//         // process the first batch of elements
//     }
//     for i := uint(0); i < o.End2(ring); i++ {
//         // process the second batch of elements
//     }
//
// The first for loop will iterate over at most [start; len(array)),
// and the second will iterate over at most [0; end).
//
// Thread Safety
//
// None of the data structures provided here are safe from data
// races. To use them in thread-safe ring buffer implementations,
// users must protect both the accounting operations and your backing
// buffer writes with a Mutex.
//
package o
