// allocating

type Allocator interface {
    Allocate(size uint) *void | error
}

type Freer interface {
    Free(ptr *void)
}

type AllocateFreer interface {
    Allocator
    Freer
}

func make(T) !*T                    // allocates a T and returns a pointer to its zero value. Panics if allocation fails.
x := &T{}                           // Same as make(T). typeof(x) is !*T

// Make has support for custom allocators. In this form, it returns nil on failure. It must keep metadata about which
// allocator was used so the runtime knows how to free the pointer.
func make(T, mem.AllocateFreer) *T

// For non-panicing allocation with the same semantics as make(T), you can use the DefaultAllocator.
x := make(T, mem.DefaultAllocator)


// If it's important to be able to see the error on failure, we could have this signature instead.
func make(T, mem.AllocateFreer) !*T | error