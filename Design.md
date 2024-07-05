A sketch for a new systems programming language. Rust and Swift are C++ replacements – too big. Zig is a C replacement, but it doesn't have any memory safety. Go is lots of fun, but GC makes it a bit too high level. CSP is great.

Basically, I want to use Go for lower level things, but it's not quite suited to some of those.

I want approximately the amount of memory safety/data race safety you get from Go, but without GC – zero values/no uninitialized variables, nil pointer dereference is well defined and will panics, default concurrency patterns (channels and goroutines) push you towards writing correct code without forcing it on you (there's still shared memory and locks if you want to shoot yourself in the foot). Doesn't have ironclad type level guarantees like Rust and Swift. It should just get you 90% of the way there.

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


Notes:

- C replacement (of course)
- Smaller than Rust
- Fun like Go
	- Fast compile times
	- CSP 
	- Mostly Imperative 
- CSP
	- No growable stacks
	- Statically allocate able stacks if necessary 
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
- Checked pointer deterrence
- Flavors - with upgrading 
	- noalloc (only stack + heap)
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
    - Iterate over code points
    - Iterate over bytes
    - Iterate over code points
    - Iterate over grapheme clusters, but a Char is not a grapheme cluster.
    - UTF-16 support (for )
- Structs
    - Default C layout? Or do you have to opt in?
- Iteration
    - co-routine based iterators ("interior iteration") 
    - Like Ruby: be able to turn a co-routine based iterator into an Enumerable/Iterator ("exterior iteration").
    - What about going back and forth? A cursor? Something else?
    - Swift-like slicing of arbitrary data structures?
- Defer
- Catchable panic?
    - I hope so. Want to be able to test custom data structures
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
- Nil or optional?
- Tagged unions?
- Error handling?
- Generics without fat pointers?
- How close to Go syntax?
- Can you add some array programming to an imperative language in a way that feels nice?
- bigint by default?
- Registers as variables?
    - I.e. better inline assembly?
- What are the target usecases?
    - Not servers. Use Go for that. Anything that can use a GC should use a GC.
    - OSs/firmware
    - Medium to high performance requirements. Doesn't have to be as fast as Rust. But you should have control of performance.
    - Go semantics with perfect C interop.
    - Memory safe C (not memory safe C++)
