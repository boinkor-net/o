# o - go ring buffers for arbitrary types without `interface{}`

This package provides the data structures that you need in order to
implement an efficient ring buffer in go. In contrast to other ring
buffer packages (and the `Ring` package in the go stdlib which really
should not count as a ring buffer), this package has the following
nice properties:

* It provides the minimum functionality and maximum flexibility
  necessary for your own ring buffer structure.
* It allows multiple modes of usage for different ring buffer usage
  scenarios.
* It does not require casting from `interface{}`.

## Minimum functionality - what do you get?

This package handles the grody integer math in ring buffers (it's not
suuuper grody, but it's not easy to get right on the first try. Let me
help!)

That's it. You are expected to use the `o.Ring` interface provided by
this package in your own structure, with a buffer that you allocate,
and you're supposed to put things onto the right index in that buffer
(with `o.Ring` doing the grody integer math).

You get two buffer data structures: One that works for all kinds of
capacities, and one that is optimized for powers of two.

## Maximum flexibility & multiple usage modes

The default usage mode for `o.Ring` is to `.Push` and `.Shift` for
LIFO operations, similar to queues and typical log buffers. You can
find an example in the `ringio` package implemented here. These
functions return errors if you push onto a full ring, or if you shift
from an empty ring.

You can also use `o.ForcePush` to insert a new element regardless of
whether the ring is full, overwriting the element that's there.

And then, if you do not want to shift out elements to read them, you
can use `o.Ring.All` to get an array containing all the occupied
indexes in the ring (for oldest to newest), and iterate over them -
it's your data structure! You get to go hog wild. `o.Ring.Rev` exists
too, in case you wanted to go from newest to oldest.

If you just can't afford to allocate another ring buffer-sized array
of `uint`, you can also use the index helpers `o.Start1`, `o.End1`,
and `o.End2` to get bounds for two `for` loops. Like I said, go _hog_
wild.

## Why do this at all?

Depending on where you intend to use a ring buffer, it might be
difficult to reason about whether what you get out is what you
expect. The error handling code for that sometimes gets grody, but
really - that isn't the reason why I did this.

Mostly, I did it as a semi-joke that could be useful. Now that I've
written this, I'm no longer sure it is a joke. People might acually
want to use this and feel good about using it, and now I'm terrified
because I think this might actually be a reasonable thing to use,
under some circumstances.
