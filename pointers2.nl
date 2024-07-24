// Another sketch, attempting a cleaner relationship between optionals and pointers.

var x int   // An int. Bridgable to C.
var x ?int  // An optional int. Not bridgable by default. Maybe we could have generated types
            // that get exposed to C, e.g optint.

// I don't love this because I really want the default pointer type to be optional but spelled *int.
// That said, doing that adds a lot of inconsistency.

var x *int  // a non-optional pointer to an int. Can be passed to C (devolves into an optional pointer),
            // but can't be nil. Can never be returned by a C function unless we build a system to annotate
            // C functions with their pointer nullability.
var x ?*int // An optional pointer to an int. Can be nil. Equivalent to *int in C, and can be bridged.
            // Panics on nil dereference. Can be pattern matched on.

var x #*int         // A refcounted pointer to an int. Cannot be nil. Cannot be passed to C without an unsafe cast.
var x ?#*int        // An optional refcounted pointer to an int. Can be nil. Cannot be passed to C without an unsafe cast.
                    // I do think this is probably necessary, but I'm not 100% sure.
var x (weak ?*int)  // An optional weak pointer to an int. There are no non-optional weak pointers. Cannot be passed to C
                    // without an unsafe cast. I'd rather spell this as (weak *int), but that would add similar inconsistencies.


var x $*int // An owned pointer. Move only. Cannot be nil. Cannot be passed to C without an unsafe cast.
            // Equivalent to (nocopy *int). I think I like !*int better, but we might want to reserve that for
            // functions that return an error (in the same way Zig does), and $ is kinda cute for ownership.
var x ?$*int // An optional owned pointer. Is this necessary?

// Still not sure if we want to disallow data races. If we do, then we need equivalents to Rust's &mut int and &int.
var x &*int         // Immutable reference.
var x (mut &*int)   // Mutable reference.

// Mutable by default:
var x &int          // Mutable reference.
var x (const &int)  // Immutable reference.

// Other alternatives:
var x (borrow *int) // Alts: (loaned *int), (borrowed *int), (const *int)
var x (inout *int)  // Alts: (mutable *int), (mut *int)
