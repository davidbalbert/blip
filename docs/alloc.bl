// allocating

type Allocator interface {
    Allocate(size uint) $*void | error
    Free(ptr $*void)
}

func make(T) $*T                    // allocates a T and returns a pointer to its zero value. Panics if allocation fails.
x := &T{}                           // Same as make(T). typeof(x) is $*T.
                                    // TODO: if x doesn't escape, should T be stack allocated?

// Make has support for custom allocators. In this form, it returns nil on failure. It must keep metadata about which
// allocator was used so the runtime knows how to free the pointer.
func make(T, mem.Allocator) $*T | error

// For non-panicing allocation with the same semantics as make(T), you can use the DefaultAllocator.
x, err := make(T, mem.DefaultAllocator)


// TODO: there's a problem with this interface. There aren't a ton of reasons to use custom allocators. But one of them
// is a slab allocator. But slab allocators, by definition, can only allocate a specific size, so this interface doesn't
// make any sense. Figure out what to do about this.


// TODO: Nested ownership. Consider a slab allocator (for which we've solved the above problem):

type SlabAllocator struct {
    elsz uint
    data $*T
}

func (a *SlabAllocator) Allocate() $*T {
    // ...
}

let a = &SlabAllocator{elsz: 8} // typeof(sa) is $*SlabAllocator
let p = make(T, a)              // typeof(p) is $*T

// p is an owned pointer into a region of memory that's owned by a. This needs to be ok. Memory allocation is
// hierarchical, and as long as each owned pointer is freed by the same allocator that allocated it, everything
// should be fine. But I'm sure there are gnarley bits here.


// A thought: SlabAllocator owns a resource. It's a region of memory. For most things that own a resource we use
// mem.NoCopy and then explicitly clean them up, but so far, memory is a special case. Freeing of owned pointers
// (and values) is automatic.
//
// Is a struct that contains an owned pointer inherently move-only? I think it is. Is it inherently mem.NoCopy?
// not right now, because we don't need to explicitly clean it up.
//
// Here's the tension: memory is a resource just like a file descriptor. So a struct that owns some memory
// needs to be unique so the memory never has more than one owner. Which makes it look a lot like mem.NoCopy.
// But it's not quite mem.NoCopy because we don't want to have to explicitly free each owned pointer. It should
// just know how to do it.
//
// This isn't bad necessarily. It's just a bump in the carpet – not totally smooth. It bears extra scruitiny.
