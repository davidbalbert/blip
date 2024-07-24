// A third take, focusing on single writer, multi reader. This time, *int is a borrowed pointer.

var x $*int // an owned pointer
var x &*int // a borrowed immutable pointer
var x *int  // a borrowed mutable pointer

// optionals
var x ?$*int
var x ?&*int
var x ?*int


// Or maybe, if we don't care about data races, we could do this:
var x $*int // an owned pointer
var x *int  // a borrowed pointer

// How does this relate to normal pointers. I'm really not sure yet.
