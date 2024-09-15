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
