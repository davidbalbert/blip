// allocating

type Allocator interface {
    Allocate(size uint) ?*void
}

type Freer interface {
    Free(ptr *void)
}

type AllocateFreer interface {
    Allocator
    Freer
}

func make(T) *T                     // allocates a T and returns a pointer to its zero value. Is this returning
                                    // a non-optional pointer or an optional pointer?
x := &T{}                           // Same as make(T)
func make(T, mem.AllocateFreer) *T  // support for custom allocators. Same semantics as make(T). Must keep metadata
                                    // about which allocator was used so that we can free the pointer.

// By default, these panic when allocation fails. Do we need a version where they don't panic? Some ideas:
var x ?*T = make(T)                 // specifying the type makes it non-panicing.
func make(T) (p ?*T, ok bool)       // alt: if we have some sort of result type, or other way to return an error, make
                                    // could return that instead.

// In general, Newlang won't have function overloading. But like Go, I'm ok to overload a few special built-in functions.
// make is already special because it takes a type as an argument, which normal functions can't do.
//
// I think I'm ok with overloading built-ins based on their number of arguments, e.g. a variant of make that returns
// (?*T, bool), but I'm not sure about overloading based just on the type of an argument (e.g. returning ?*T vs *T).
