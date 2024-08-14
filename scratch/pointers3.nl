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

// How does this relate to normal pointers? I'm really not sure yet.

// Owned pointers are nocopy. I.e. $*int is the same as (nocopy *int). Perhaps C functions that return
// pointers, should always return owned pointers. Perhaps if we have rules that any nocopy value has
// to either be consumed, then we can guarantee that owned pointers are always freed. E.g.
//
// TODO: Obviously we need to be able to escape owned pointers as well. Figure out how to deal with that.

// If these C declarations
void *malloc(size_t size);
void free(void *ptr);

// Were imported imported as
func malloc(size C.size_t) $*void
func free(ptr $*void)

// Then we'd be guaranteed that any allocation would be freed.

import "C"

func bad() {
    p := C.malloc(10)
    // ...
    // error: p must be consumed
}

func good1() {
    p := C.malloc(10)
    // ...
    C.free(p);
}

func good2() {
    p := C.malloc(10)
    defer C.free(p)
    // ...
}

// I think it might be a reasonable default to have C functions always return owned pointers. But it's
// definitely not reasonable to assume that all pointer arguments are owned. In tons of calls, the pointer
// is borrowed instead of owned. So what if we imported the same functions like this:

func malloc(size C.size_t) $*void
func free(ptr *void)

// Putting aside the question of what *void means – perhaps in this world we need (unsafe *void), we could
// do this:

func good() {
    p := malloc(10)
    defer C.free(eat p)
    // ...
}

// Alt spellings of eat: move, consume, give, donate.

