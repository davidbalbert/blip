# Newlang

A sketch for a new systems programming language to replace C.

- Fun like Go
- More memory safe than Zig
- Smaller than Rust
- Perfect C interop

The most complete sketch is in [newlang.nl](https://github.com/davidbalbert/newlang/blob/main/newlang.nl).

"Newlang" is a placeholder.

## Details and stray thoughts

Name idea: people sometimes use chain mail gloves when cutting things on a mandolin. It lets you use the sharp tool more safely. This is what I want for the language. Maybe there's a name that evokes that.

I wish Snap were available. It evokes a bit of danger (breaking something), power, speed, and ease ("it's a snap!"). Cinch is ok ("it's a cinch"), but I don't love it.

Rust and Swift are like C++ – too big and not much fun. Zig is smaller, but doesn't have enough memory safety. Go is lots of fun, but GC isn't right for all software. CSP is great.

Basically, I want to use Go for lower level things, but it's not quite suited to some of those.

I want approximately the amount of memory safety you get from Go, but without GC – zero values/no uninitialized variables, nil pointer dereference is well defined and will panic, default concurrency patterns (channels and goroutines) push you towards writing correct code without forcing it on you (there's still shared memory and locks if you want to shoot yourself in the foot). Doesn't have ironclad type level guarantees like Rust and Swift. It should just get you 90% of the way there.

Performance is an open question. One option: just as much speed and control as Rust, just with fewer ideas and abstractions. Maybe this is possible given the relaxed memory safety ideas above? Another option: 90% of the performance guarantees as Rust. We may already be doing this by using zero values.

Don't shoot for zero-cost abstractions. We want nice abstractions, and nice abstractions sometimes have costs. If you need more performance, you should be able to program ergonomically without the abstraction.

Test cases:
- An OS kernel
    - Can we start goroutines with statically allocated stacks?
    - Can interrupt handlers write events to a buffered channel?
        - Can we guarantee we never deadlock? Do we need dynamic allocations for this? Channel likely needs an infinite buffer.
    - Device drivers.
- UI toolkit
    - System events delivered through a channel.
        - What happens to the responder chain?
    - Background loading and saving of documents.
    - Main thread restrictions?


Notes (much of this is out of date):

- C replacement (of course)
    - https://www.humprog.org/~stephen/research/papers/kell17some-preprint.pdf
- Smaller than Rust
- Fun like Go
	- Fast compile times
	- CSP
	- Mostly Imperative
- CSP
	- No growable stacks
	- Statically allocate-able stacks if necessary
	- Set stack size (global? per stack?)
	- Interrupt handler??
- No GC
- More memory safe than Zig
	- Multi-reader single writer?
	- Definitely no explicit lifetimes
- Closures
	- Convert to and from func+context
- Fast to compile
- Good C interop
    - Easy to import and call C functions
    - Easy to export a C header
- Generics of some sort
    - associated types?
- Zero values
- One method interfaces and structural conformance
- Checked pointer dereference
- Flavors - with upgrading
	- noalloc (only stack + global)
	- nomultiplex
	- full
- Special casing ok
	- Fixed-size arrays
	- Slices
	- Dictionaries
- Easy cross compile
	- All archs
	- Easy C interop (like `zig cc`)
- Good string support
    - A bit more high level than Go.
    - Iterate over bytes
    - Iterate over code points
    - Iterate over grapheme clusters, but a Char is not a grapheme cluster.
    - UTF-16 support (for interfacing with UI frameworks)
- Structs
    - Default C layout? Or do you have to opt in?
- Iteration
    - co-routine based iterators ("interior iteration")
    - Like Ruby or Go: be able to turn a co-routine based iterator into an Enumerable/Iterator ("exterior iteration").
    - What about going back and forth? A cursor? Something else?
    - Swift-like slicing of arbitrary data structures?
- Defer
- Catchable panic?
    - I hope so. Want to be able to test custom data structures
- Errors not exceptions
    - return instead of throw
- No function/method overloading
- Type inference of local variables only
- Guaranteed cleanup?
    - Go uses finalizers, which rely on GC
    - Swift uses deinit on refcounted classes and non-copyable types (move-only types)
    - C++: RAII
        - Easy for stack allocated local variables. How does it work for heap allocations?
    - An option: just support defer, and don't guarantee anything.
    - **nocopy types must be past to a function that consumes them.**
- Reference counting
    - Maybe normal structs don't have deinit, but you can assign a deinit (multiple deinits?) to a refcounted pointer.
    - maybe nocopy types have a deinit?
        - At least nocopy structs/enums should have a deinit.
        - Is it weird if nocopy structs/enums can have a deinit but nocopy/owned pointers can't?
    - What about owned pointers? Should they have a deinit too?
        - **No! "owned" is equivalent to "nocopy", so all they need is to be passed to a function that consumes them.**
            - This would allow us to use defer rather than RAII and still be guaranteed that we clean things up.
        - Could a refcounted pointer (#*Foo) could contain a nocopy type (e.g. Fd) without being nocopy itself? That would
          allow us to force #*Foo to have a deinit. But types that contain nocopy types must also be nocopy.
            - **You can have a refcounted pointer to a nocopy type!**
- Bootstrapping
    - Ship compiler compiled to VM bytecode (like Zig)
    - But don't use Web Assembly - it doesn't support stack switching yet. Instead, make a custom tiny vm (tvm) with a reference implementation in C.
    - Ideas: https://tinlizzie.org/VPRIPapers/tr2015004_cuneiform.pdf
    - https://github.com/WebAssembly/stack-switching
- Declarative register files?
    - https://github.com/apple/swift-mmio
- Parsers?
    - PEG?
    - Ohm?
- Custom allocators?
    - Zig's approach, while principled feels really heavy. Can we have the same flexibility without having to pass the allocator in? Do we even need this?
- Array programming?
    - GPU compute? Probably getting out of scope.
- Nil or optional?
- Tagged unions?
- Error handling?
- Generics without fat pointers?
- How close to Go syntax?
- bigint by default?
- Registers as variables?
    - I.e. better inline assembly?
- What are the target usecases?
    - Not servers. Use Go for that. Anything that can use a GC should use a GC.
    - OSs/firmware
    - Medium to high performance requirements. Doesn't have to be as fast as Rust. But you should have control of performance.
    - Go semantics with perfect C interop.
    - Memory safe C (not memory safe C++)
- What operations/algorithms are common enough in these usecases to deserve a special case language construct?


### Interesting resources

- https://www.humprog.org/~stephen/research/papers/kell17some-preprint.pdf
- https://www.ralfj.de/blog/2018/07/24/pointers-and-bytes.html
- https://www.cis.upenn.edu/~stevez/papers/KHM+15.pdf
- https://github.com/NICUP14/MiniLang
