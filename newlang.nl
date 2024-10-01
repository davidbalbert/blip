// Primitive types
int8
int16
int32
int64
int128

uint8 (byte)
uint16
uint32
uint64
uint128

float16
float32
float64

// Width of a data register.
int
uint

// Width of a pointer
uintptr

// C types
C.char
C.schar
C.uchar
C.short
C.ushort
C.int
C.uint
C.long
C.ulong

// Bridging to C
//
// Sized ints, uints and floats can be passed to C functions with no cast.
// TODO: What about float16? There are multiple different float16 types in C.
//
// Unsized ints (int, uint), must be cast to a specific C type.
//
// TODO: how are C structs with bitfields imported?

// Go has complex64 and complex128. Not sure if we want those.


// An array. E.g. [5]int. Size is known statically.
[N]T

// A slice. E.g. []int. Length and capacity are stored in the slice.
// TODO:
// - Can slices be nil, like Go?
[]T

// A map. E.g. map[string]int. Alt: dict[string]int.
map[K]V

// A single Unicode code point.
rune (int32)

// A string is an sequence of arbitrary bytes with no assumed encoding. Source files are always UTF-8, and
// thus string literals are as well. Like slices, strings store their length in bytes, which can be queried
// in constant time. Strings are null-terminated, and can be bridged to C without copying. Strings may
// contain user-inserted null bytes. User-inserted null bytes are included in the length, but the null
// terminator is not. If a string with a user-inserted null byte at index i is passed to a C function, the=
// C function will see a string of length i.
//
// TODO:
// - See slices.txt for a bunch of open questions about slices and strings.
// - How to deal with extended grapheme clusters?
string

// Composite types

// Structs
struct {
    x int
    y int
}

// or
struct {
    x, y int
}

// Struct embedding like Go.
// This struct has all the same methods as Conn. It also has an embedded int, though ints don't
// have any methods.
struct {
    Conn
    int
}

// TODO: We'd like to have a way to represent header + data in a struct. In C, you'd use a flexible,
// array member, e.g.
//
// struct string {
//     len int;
//     char data[];
// }
//
// We could do the same in Newlang (here's a possible spelling), but I wonder if there might be a better
// way to express things.
//
type buf struct {
    len int
    data ...byte // Alt: `data byte...`
}

// I'm not sure if you can put this the stack or not. While the size can't be known by looking at the struct
// definition itself, as long as the size of data doesn't change, its equivalent to alloca(3).
//
// Heap allocating this is easy – the size (sizeof(int) + len) is dynamically stored in the heap, and dropping
// is therefore dynamic as well.


// Tagged unions. Basically enum with associated values.
union {
    one int
    two (int, int) // This is not a tuple. Like multiple return values, it's a syntactic construct.
}

// Unions also support types without an explicit tag name. This is conceptually equivalent to something
// like `Conn | File`. Unlike structs, methods are not promoted. Instead, consider interfaces, or if you
// have a union, pattern matching.
//
// Note: if we decide to enable method promotion on unions, we should only do it when the union is composed
// entirely of embedded types, and when the method is present on every type in the union.
union {
    Conn
    File
}

// Untagged union. Equivalent to C union. Alt: untagged union, extern union.
unsafe union {
    i int
    f float64
}


// Enums

// No underlying type. Can't be bridged to C.
enum {
    north
    south
    east
    west
}

// unsized, so no bridging.
enum {
    north int  // 0
    south      // 1
    east       // 2
    west       // 3
}

// sized, bridgable
enum {
    north int8
    south
    east
    west
}

// Imported from C. You can define these, but it's better to use sized types.
//
// TODO:
// - How are enums imported? Requirements on C enums are pretty loose, and I want enums to be more strict.
// - Is it possible to import C enums that are defined in other ways? E.g. with NS_ENUM in Apple's APIs?
enum {
    north C.int
    south
    east
    west
}

// Specifying values
enum {
    north int = 1
    south         // 2
    east = 10
    west          // 11
}

// String values
//
// TODO: Are these bridgable to C? If so, how?
enum {
    north string // "north"
    south        // "south"
    east         // "east"
    west         // "west"
}

enum {
    north string  // "north"
    south         // "south"
    east = "East"
    west          // "west"
}

// TODO:
// - Exhaustiveness checking.
// - Can you create an enum out of an a raw value? If so, what happens when it's not a valid value?
// - Do we have frozen/non-frozen enums like Swift? In Swift, for non-frozen enums, which can have
//   new cases added later, you need to have an "@unknown default" case when pattern matching.

// Bitfields can't be enums because they're not disjoint.

// One option: consts that work just like Go, including textual substitution and iota.
type flag int16
const (
    flag1 flag = 1 << iota
    flag2
    flag3
    all = flag1 | flag2 | flag3
)

// Unlike Go, let's allow reasigning iota.
const (
    flag1 = 1 << iota // 1  (iota == 0)
    flag2             // 2  (iota == 1)

    iota = 6
    flag4 = 1 << iota // 64 (iota == 6)
    flag5             // 128 (iota == 7)

    all = 255
)

// Traditionally, to remove a flag, you have to AND with the inverse. E.g.
f := all & ^flag1 & ^flag2

// Could we have a more readable syntax for this? Binary ! might work, but it could be hard to parse.
f := all ! flag1 ! flag2

// Perhaps |-? I don't love this.
f := all |- flag1 |- flag2


// Type declarations. For now, these work like Go. They create new types, must be explicitly converted, and
// the new type doesn't inherit any methods defined on the underlying type.
type flag int

type point struct {
    x, y int
}

type dir enum {
    up, down, left, right
}

// Equivalent to `type byte = uint8` in Go. Allows for implicit conversion.
alias byte uint8

// Non-copyable types.
//
// Non-copyable types are move-only. They must be escaped or explicitly consumed. Non-copyable types are
// useful for types that manage resources and must be cleaned up later, like a net.Conn that ownes a
// file descriptor.

// If a type embeds a non-copyable type, it is also non-copyable. The type mem.NoCopy is a zero-sized
// non-copyable type. C.f. structs.HostLayout in Go. To make a non-copyable struct, embed mem.NoCopy:
type Fd struct {
    mem.NoCopy
    fd int
}

// Because Fd is non-copyable, Conn is non-copyable too:
type Conn {
    fd Fd
    // ...
}

// An example:
var a Fd = open(...)
b := a
a.Read(...) // error: a was moved to b
b.Read(...) // ok

// Unions are non-copyable if any of their associated values are also non-copyable. You can also make a
// union non-copyable by including mem.NoCopy as an unamed case:
union {
    _ mem.NoCopy
    s string
    n int
}

// Enums are always copyable.
//
// In an IDE, a non-copyable type is displayed like this: "Fd (nocopy)"


// Math

// Arethmetic operators (+, -, *, /, %, etc.) are defined to trap on overflow (TODO: which ones?).
// .+, .-, and .* (more?) wrap on overflow.
//
// TODO: what about a way to explicitly check the overflow flag, e.g. math.Overflow()? If we do that, what does
// that mean in the following cases:
// - Array operations (see below)
// - Compound expressions (e.g. a + b*c)
// - etc.

// Pointers

// There are 5 pointer types: owned, borrowed, unsafe, reference counted, and weak. All pointer types can be nil.
// Dereferencing a nil pointer is defined behavior – it panics. All pointers besides unsafe pointers have temporal
// safety – they prevent use-after-free.

$*int       // owned
*int        // borrowed
!*int       // unsafe
#*int       // reference counted
(weak *int) // weak

// Other possible spellings and pointer types
~*int // weak
%*int // Non-null, if we add it. Would this be borrowed or owned? Would we need both? Can we get away without this?

// Owned pointers
//
// Owned pointers own the memory they point to, and are responsible for freeing that memory, and disposing of the
// underlying object when they go out of scope. A piece of memory can have exactly one owner. Owned pointers are
// move-only. This ensures that the "single owner" invariant holds. Owned pointers can be passed as arguments to
// functions or returned. Because they are move-only, doing either transfers ownership of the referenced memory.
// An owned pointer's referenced memory is always heap allocated.
//
// When an owned pointer goes out of scope, the memory it refers to is freed. Owned pointers to nocopy types must
// be explicitly dropped. Owned pointers to copiable types are dropped implicitly.
//
// To pass an owned pointer to a C function, you must explicitly cast it to an unsafe pointer. This makes the C
// function responsible for freeing its memory. This cast is considered a move, and the owned pointer can no longer
// be refered to after the cast.
//
// Owned pointers can be set to nil, or to a new address. If the pointer was non-nil previously, the memory it pointed
// to is freed. You cannot set an owned pointer to nil or a new address if it is being borrowed. You can't move an
// owned pointer if it's being borrowed either.
//
// Owned pointers are fat pointers – they need to know how to free their memory:

struct {
    p !*T
    free func(!*T)
}

// free is derived from the mem.Allocator that allocated p. Alternatively, we could store a mem.Allocator value, but
// then the owned pointer would be 3 words instead of 2 (interface values themselves are fat pointers).


// Borrowed pointers
//
// Borrowed pointers can be created from owned pointers, reference counted pointers as well as non-pointers, and can't
// outlive the value they point to. No action is performed on drop, and they never have to be dropped explicitly. When
// passed to C, the the programmer is responsible for ensuring that if the pointer is escaped, it doesn't outlive
// its referent.
//
// TODO: it should be possible to reassign a borrowed pointer to a new address. Specifically, this is useful for
// things like iterating over a linked list (see this Rust example from https://rust-lang.github.io/rfcs/2094-nll.html).
// Specifically note the `list = n` reassignment:
//
// struct List<T> {
//     value: T,
//     next: Option<Box<List<T>>>,
// }
//
// fn to_refs<T>(mut list: &mut List<T>) -> Vec<&mut T> {
//     let mut result = vec![];
//     loop {
//         result.push(&mut list.value);
//         if let Some(n) = list.next.as_mut() {
//             list = n;
//         } else {
//             return result;
//         }
//     }
// }
//
// TODO: what about setting a borrowed pointer to nil or initializing one to nil (or creating a zero-valued struct that
// contains a borrowed pointer)? Nil has a static lifetime, so I think this should be ok? But what about when a struct
// has a non-exported borrow? Possible solution: all zero valued structs have a static lifetime. The only way to assign
// to a non-exported borrow or to create an instance with a non-zero value is through a function. And annotations on
// function signatures should be able to restrict the lifetime of returned struct (or the struct that's mutated).

// Unsafe pointers
//
// Unsafe pointers have unknown ownership. The programmer is responsible for freeing the memory when appropriate.
// They still panic on nil dereference, but use-after-free is possible. All C pointers are imported as unsafe.
// Alt: (unsafe *T), (raw *T), (unowned *T).


// Reference counting

// A reference counted pointer owns its memory and frees it when its reference count reaches 0. Like owned pointers,
// a reference counted pointer cannot be set to nil or to point at another reference counted value while it is
// being borrowed.

// In pseudocode, a refcounted pointer is a fat pointer. One of the fields points to the value, and the other to the
// metadata, which includes the reference count. The count is updated atomically.

struct {
    meta !*struct {
        refcount int
        cleanup func($*T)
    }
    value !*T
}

// To make a refcounted pointer, use the rc builtin. The metadata is allocated on the heap and the bytes of v are
// copied in. You can also pass in an owned pointer. In these forms, T must be copyable.
func rc(v T) #*T
func rc(p $*T) #*T

// You can also supply a cleanup function that will be called when the refcount reaches 0. In this form, T may
// be non-copyable.
func rc(v T, cleanup func(p $*T)) #*T
func rc(p $*T, cleanup func(p $*T)) #*T

// TODO:
// - it feels a bit odd that cleanup takes an owned pointer in both cases. I believe this is correct, but it would
//   somehow feel nicer if the first form took the value itself.
// - Maybe we shouldn't have the non-pointer taking form at all. E.g. if you want to create a refcounted pointer to
//   a value, you need to take the address of that value:

rc1 := rc(&5)

p1 := make(int32, anAllocator)
rc2 := rc(p1)

// You can also integrate external reference counted types by providing custom retain and release functions.
//
// IsRetained should be true if the value or pointer is given to you with a +1 refcount.
func rc(v T, retain func(v T), release func(v T), isRetained bool) #*T
func rc(p !*T, retain func(p !*T), release func(p !*T), isRetained bool) #*T

// TODO:
// - passing in isRetained feels a bit ugly. Is there a better way?
// - The idea behind having both T and !*T is that the former could be a struct that contains a reference counted pointer.
//   That makes sense, but it raises a question – will the type of the value always mirror the type of the parameter passed
//   to retain and release?


// A custom refcounted pointer has one of the following layouts in memory (again, pseudocode):

struct {
    meta !*struct {
        retain func(T)
        release func(T)
        refcount int
    }
    value T
}

struct {
    meta !*struct {
        retain func(!*T)
        release func(!*T)
        refcount int
    }
    value !*T
}

// TODO: in the first case, we don't really have a fat pointer because we're storing the value directly. We're doing that because
// we assume the value is some opaque type that wraps a pointer. This means we have to know the size of the value at compile time.
// But I believe we'll always know this.


// You can make a weak reference using the weak builtin
//
// TODO: how to make weak references to custom refcounted pointers? We need to know when the refcount reaches 0
// so we can nil out the weak pointers. One option: weak pointers could be nil'd out as soon as the number of
// #*T pointers (i.e. the number of strong references from Newlang code) drops to 0. This is a bit ugly though:
// if some C code is holding a strong reference, the weak pointer will be nil'd out even though the object is
// still alive.
func weak(p #*T) (weak *T)

p := = rc(5) // typeof(p) is #*int
w := weak(p)  // typeof(w) is (weak *int)

p := rc(fd, close) // typeof(p) is #*Fd


// Pointer conversions

// A variable can be borrowed multiple times.
var x int
p1 := &x           // typeof(p1) == *int
var p2 *int = &x  // type can be explicitly specified

// A local variable can be moved into an owned pointer, in which case the original variable is consumed. Escape
// analysis is performed. If the variable escapes (e.g. is returned) or is consumed by a function, it will be
// heap allocated.
var x int
var p $*int = &x  // may be stack or heap allocated
print(x)          // error: x was moved to p. It doesn't matter that x is copyable.

// Another syntax for the above
var x int
p := (*int)(&x)

// TOOD: I don't love that the type of &x can be either $*int or *int depending on its context. But I do want
// the following examples to work, and they seem to require it. This could be solved by having another operator
// in addition to '&'.

func foo() $*int {
    var x int
    return &x // x is heap allocated
}

func bar($*int) { ... }
func foo() {
    var x int
    bar(&x) // x is heap allocated
}


// Taking the address of a literal works the same way. If it escapes, it's heap allocated, otherwise it's
// stack allocated. Either way, the pointer is owned.
var p $*int = &5
p := &5

// You can't borrow a literal value
var p *int = &5 // error: you cannot borrow a literal

// To force heap allocation, use make.
p := make(int) // typeof(p) is $*int. The int is on the heap.


// You can borrow an owned pointer.
var x int
p1 := ($*int)(&x)
var p2 *int = p1

// On the other hand, if a variable has been borrowed, it can't be moved into an owned pointer.
var x int
p1 := &x           // typeof p1 is *int
var p2 $*int = &x  // error: p1 borrows x, so p2 can't own it.


// Passing pointers to C functions.

// void inc(int *p);
func inc(p !*int)

// To call inc with a borrowed pointer, do the following.
var x int
p := &x    // typeof(p) is *int
inc(p.!)

// This also works with refcounted pointers. The reference count is not incremented.
p := rc(5) // typeof(p) is #*int
inc(p.!)

// And weak pointers:
p := rc(5)
w := weak(p)
inc(w.!)

// In all of these cases, you must ensure that the unsafe pointer doesn't outlive the value it points to. The
// easiest way to do this is to make sure inc doesn't escape the pointer. If call a function that does escape
// the pointer, you must ensure that the pointer's referent lives at least as long as the unsafe pointer does.

// Owned pointers can also be passed to C functions, but the conversion causes ownership to be transferred. Unless
// you free the memory later, it will be leaked. Freeing the memory can be tricky – it depends on what allocator
// allocated the memory.
p := &5
inc(p.!) // $*int -> !*inc
// p has leaked

// An unsafe pointer of any type can be used where an unsafe void pointer is expected.
// void printptr(void *p);
func printptr(p !*void)

var x int
p := (&x).!      // typeof(p) is !*int
printptr(p)      // !*int -> !*void
printptr((&x).!) // *int -> !*void


// !*void can be converted to a typed unsafe pointer using p.(!*T).
//
// void call(void (*f)(void*), void *arg);
func call(f func(!*void), arg !*void)

func printint(p !*void) {
    // TODO: can we make this panic if p doesn't point to an int?
    // See liballocs: https://github.com/stephenrkell/liballocs
    ip := p.(!*int) // !*void -> !*int
    fmt.Println(*ip)
}

x := 5
call(printint, (&x).!) // *int -> !*void


// It is possible to safely escape a reference counted pointer using the unsafe.Retain and unsafe.Release
// builtins. The former converts its argument to an unsafe pointer and increments the reference count, the
// latter does the opposite conversion with a decrement.
//
// TODO: could Release be made to panic if p isn't reference counted?
func Retain(p #*T) !*T
func Release(p !*T) #*T

// This example assumes the presence of a single threaded event loop.
//
// void callAfter(void (*f)(void*), void *arg, double seconds);
func callAfter(f func(!*void), arg !*void, seconds float64)

func printint(p !*void) {
    ip := p.(!*int)
    fmt.Println(*ip)
    unsafe.Release(ip)
}

p := rc(5)
callAfter(printint, unsafe.Retain(p), 1.0)


// Receiving pointers from C functions.

// Consider the following:
//
// struct Node;
// char *get_name(struct Node *node);
func get_name(node !*C.Node) !*C.char

// Node owns its name. In other words, the !*C.char returned by get_name is a borrow by convention. Callers
// of get_name must ensure that the name doesn't outlive the node. While it's possible to use the unsafe
// pointer directly, you can also convert it to a borrow by specifying the borrow's lifetime.

n := &C.Node{}           // typeof(n) is $*C.Node
s1 := get_name(n.!)      // typeof(s1) is !*C.char
s2 := get_name(n.!) in n // typeof(s2) is *C.char. Alt: .(in n) or in! n.

// TODO: does this mean you can't call a method on an unsafe pointer? This would be bad.

// You can't bind the lifetime of a borrow to an unsafe pointer. We don't know how long it will live:
n := C.make_node()     // typeof(n) is !*C.Node
s1 := get_name(n)      // typeof(s1) is !*C.char
s2 := get_name(n) in n // error: cannot bind lifetime of borrow to an unsafe pointer


// Consider the following:
//
// char *copy_name(struct Node *node);
func copy_name(node !*C.Node) !*C.char

// The semantics of copy_name is that the caller owns the returned string and is responsible for freeing it.
// You can convert the unsafe pointer to an owned pointer using unsafe.Claim:
func Claim(p !*T, alloc mem.Allocator) $*T

n := &C.Node{}
s1 := copy_name(n.!)                   // typeof(s1) is !*C.char
s2 := unsafe.Claim(s1, mem.CAllocator) // typeof(s2) is $*C.char.

// mem.CAllocator is a built-in allocator that uses C.malloc and C.free.
//
// TODO: It's a bit wierd for mem to depend on C, but I can't think of anything better. Some options:
// C.Allocator is probably the clearest, but I think it's bad to have a type in the C that's not actually
// a C type. unsafe.CAllocator might also work. It goes with unsafe.Claim, but nothing about it is
// actually unsafe.

// You can also free the string manually:
n := &C.Node{}
s := copy_name(n.!) // typeof(s1) is !*C.char.
defer C.free(s)

// If you don't do this, s will leak.


// Borrowed pointers and lifetime dependence

// The goal is to prevent use-after-free and double-free. In other words, to provide temporal safety.
//
// To start, all memory has a single owner. This owner is a variable that's responsible for freeing the
// memory at the appropriate time.
//
// Local variables and owned pointers own their memory and free it at the end of the variable or pointer's
// lexical scope, by manipulation of the stack pointer or by using the appropriate allocator respectively.
//
// A reference counted pointer also owns its memory, and frees it when its reference count reaches 0.
//
// A borrowed pointer doesn't own the memory it points to. Instead, it is a temporary reference to memory
// owned by someone else. In order to prevent use-after-free and double-free bugs, a borrowed pointer cannot
// outlive the memory it points to. When an owned or reference counted pointer is being borrowed, it must
// not be set to nil or made to point to a different address.
//
// The compiler guarantees that a borrowed pointer cannot outlive its referent:

// This is fine. P is alive until the end of the function, and doesn't outlive x.
func foo() {
    var x int = ...
    p := &x // typeof(p) is *int. ok: p doesn't outlive x.
}

func foo() {
    var p *int
    {
        var x int = ...
        p = &x // error: p outlives x
    }
}


func bar(p *int) {
    // ...
}

func foo() {
    var x int = ...
    bar(&x) // ok: the borrow of x lives until bar returns.
}

// Under some circumstances, borrows can be returned from functions and methods.

// The lifetime of res is the smaller of the lifetimes of a and b.
func max(a, b *int) res *int {
    if a > b {
        return a
    }
    return b
}

// You can explicitly specify lifetimes if the result depends on only some of the arguments.
// In this case, the lifetime of res is allowed to outlive the lifetime of b, but not a.
func add(a, b *int) *int in a {
    *a += b
    return a
}

// This works even if the borrows are of different types.
func foo(t *T, u *U, cond bool) *int {
    if cond {
        return &t.x
    }
    return &u.x
}

// The rules are the same for methods. The receiver is treated as another argument.
func (t *T) foo(u *U, cond bool) *int {
    if cond {
        return &t.x
    }
    return &u.x
}

func foo(t *T, u *U, cond bool) res *int in t {
    if cond {
        return t.x += u.y
    }
    return &t.x
}

// What about re-assigning borrows?
func swap(a, b *int) {
    a, b = b, a // error: a and b must have the same lifetime to reassign.
}

// Why is the above the case? If the above were allowed, you could do this:
func caller() {
    var x int
    px := &x
    {
        var y int
        py := &y
        swap(px, py)
    }
    print(*px) // error: use-after-free because after the swap px points to y, which has been freed.
}

// To fix this, we can explicitly specify lifetimes. In the above example, the call to swap would shorten
// the lifetime of px to the end of the inner block.
//
// TODO: this is hard to read, can we fix it? Maybe `func swap(a, b *int in .)` or something like that?
func swap(a in b, b in a *int) {
    a, b = b, a
}

// Of course, you can swap values rather than addresses without worrying about lifetimes:
func swap(a, b, *int) {
    *a, *b = *b, *a
}

// Should this be allowed directly, or do you have to specify `*int in package`?
func alwaysNil() *int {
    return nil
}

// What about closures?
//
// Options:
// - perform escape analysis to determine that x should be heap allocated, and deallocated when the closure is dropped.
// - error: the closure outlives x and therefore cannot capture its value directly.
func caller() func() *int{
    var x int

    return func() *int {
        return &x
    }
}

// A similar example:
// - perform escape analysis to determine that x should be heap allocated and is owned by the closure.
// - error: the closure outlives x, and therefore cannot capture px.
func caller() func() *int {
    var x int
    px := &x

    return func() *int {
        return px
    }
}


// What is the lifetime of p? Nil's lifetime is static. Does `p = &x` shorten p's lifetime to that of x? It seems like
// the logical thing to do, but I don't know how expensive that would be.
func foo() {
    var x int
    var p *int = nil
    p = &x
    print(*p)
}



// A composite type with a borrowed pointer is considered borrowed. A composite type containing another composite type
// is also borrowed and obeys the same rules as above.

// Consider this struct
type pair struct {
    x, y *int
}

// What happens when we do this? One possibility is to capture the lifetimes as they go in and out of the struct.
func foo() {
    var x int
    var z1 *int
    var z2 *int
    {
        var y int

        p := pair{x: &x, y: &y}
        z1 = p.x
        z2 = p.y // error: p.y outlives y
    }
}

// Another option is to assume that the lifetime of any references in the struct are the same as the lifetime of the struct itself.
func foo() {
    var x int
    var z1 *int
    var z2 *int
    {
        var y int

        p := pair{x: &x, y: &y}
        z1 = p.x // error: z1 outlives p
        z2 = p.y // error: z2 outlives p
    }
}

// If we added a "." lifetime to refer to the lifetime of the composite type itself, then the original definition of pair would be
// equivalent to the following:
type pair struct {
    x, y *int in .
}



// TODO: Iterators and loops
//
// For custom data structures, I want Go-style coroutine-based internal-iteration by default. But we're going
// to be able to convert those into external iterators. The question is how to make sure the underlying data
// structure is retained by the iterator. The answer is that the iterator should borrow the data structure
// and that should keep it alive. It's a good test for the lifetime system. It has to be at least that expressive.
//
// Should borrows be non-nil or nilable? Thoughts:
//
// - In Go, you can call a method on a nil pointer. It's sometimes useful, but how useful?
// - If we want to have a data structure that borrows another one (which we do), that field will be a
//   borrowed pointer. If we want to have zero-value initialization by default (which we do), that field will
//   have to be able to be nil. That's the zero value for any pointer type.
//
// So what does this mean for lifetimes? An owned pointer can be nil. If you borrow it, the borrow will be
// nil too.
//
// One problem with the above: what's the best way to make a doubly linked list? One option: use #*List for next
// and (weak *List) for prev. That's what Swift would do (modulo whatever's going on in indirect enums). Could you
// have
//
// Owned pointers are freed when they go out of scope. Their memory must also be freed when set to nil or overwritten
// with another value. This means it has to be an error to make an owned pointer point to a new location or nil while
// it is being borrowed. Ditto for a borrowed pointer – we can't nil out a borrowed pointer or update the location it.
// points to. But we can update its fields.
//

// Is it possible to mutate while iterating? Here's a stardard, immutable for loop:
for n := range []int{10, 20, 30} {
    // n is the value
}

for i, n := range []int{10, 20, 30} {
    // i is the index, n is the value
}

// Ways to mutate in Go
nums := []int{10, 20, 30}
for i, _ := range nums {
    if nums[i] == 20 {
        nums[i] == 50
    }
}

// If we could borrow the value, then we wouldn't have to use index above. That on its own doesn't seem valuable enough,
// but it would be valuable on a custom data structure, e.g. a linked list.
//
// Maybe something like this? Is this useful enough?
nums := []int{10, 20, 30}
for &n := range nums {
    if *n == 20 {
        *n = 50
    }
}

// For custom data structures using range-over-func, it might be as simple as having the yield signatures be
func() bool
func(*V) bool
func(K, *V) bool

// If you define a yield function for one of the pointer variants, it can still work with the loop variable `n` instead of `&n`.
// It would just be an implicit dereference.
//
// I'm not sure if there are lifetime issues here. What if the values have borrows in them?
//
// This doesn't support inserting or removing. That's too much of a lift. Use a data structure specific Cursor (or maybe Iterator)
// and custom methods.
//
// All that said, maybe this is too hard and too much. Go doesn't have mutable iterators, and after all, you can do everything
// with functions.


// Functions

func inc(i int) int {
    return int + 1
}

func inc(p *int) {
    *p++
}

// Functions can return multiple values. This is not a tuple.

func divmod(a, b int) (int, int) {
    return a / b, a % b
}




// Error handling

// Errors are values. They can be returned from functions and passed around like any other value.
func canError() error {
    return anErr
}

// Functions that can return both a value and an error are indicated by an error union: `int | error`.
// Error unions are not the same as normal tagged unions. They have special ergonomics for error handling.
func div(a, b int) int | error {
    if b == 0 {
        return errDivideByZero
    }
    return a / b
}

// Functions that return multiple values along with an error are spelled like this:
func divmod(a, b int) (int, int) | error {
    if b == 0 {
        return errDivideByZero
    }
    return a / b, a % b
}

// At the call site, error values look like an additional return value. But the error must be checked
// before using the other return values.
func bad1() {
    q, r, err := divmod(5, 2)
    print(q, r) // error, err must be checked
}

func bad2() {
    q, r, err := divmod(5, 0)
    if err != nil {
        // no return
    }
    print(q, r) // error, err must be checked before using q and r.
}

func good() {
    q, r, err := divmod(5, 2)
    if err != nil {
        return
    }
    print(q, r)
}

// If we don't touch q and r, we don't need to check the error, but we'll get a warning for
// unused variables.
func warning() {
    q, r, err := divmod(5, 2) // warning: q, r, and err are unused.
}


// For functions that just return an error, we're not forced to check the error, but we'll get
// warnings if we ignore the return value, or if we assign it to a variable that's never used.
func canError() error {
    return anErr
}

func warns1() {
    canError() // warning: return value ignored
    // ...
}

func warns2() {
    err := canError() // warning: err is unused
    // ...
}

// explicitly discarding the error makes the warning go away
func good() {
    _ = canError()
    // ...
}


// Error recovery (similar to Zig's errdefer). Consider a function that opens two files.
func bad(p1, p2 string) (Fd, Fd) | error {
    fd1, err := open(p1)
    if err != nil {
        return err
    }

    fd2, err := open(p2)
    if err != nil {
        return err // error: fd1 must be consumed
    }

    return fd1, fd2
}

func good1(p1, p2 string) (Fd, Fd) | error {
    fd1, err := open(p1)
    if err != nil {
        return err
    }

    fd2, err := open(p2)
    if err != nil {
        close(fd1) // no error: fd1 is consumed
        return err
    }

    return fd1, fd2
}

func good2(p1, p2 string) (Fd, Fd) | error {
    fd1, err := open(p1)
    if err != nil {
        return err
    }

    // Recover is like defer, but it only runs on branches where the function returns an error.
    recover close(fd1)

    fd2, err := open(p2)
    if err != nil {
        return err // no error: fd1 is consumed by recover.
    }

    return fd1, fd2
}


// Errors don't have to be returned in error unions. They are values like any other.
func div(a, b, float64) (float64, error) {
    if b == 0 {
        return Inf, errDivideByZero
    }
    return a / b, nil
}

func good() {
    q, err := div(5, 0) // warning: err is unused
    print(q)
}


// `or` can be used to provide a default value when a function returns an error.
//
// Alts: `catch`, `handle`, `else`.

func divmod(a, b int) (int, int) | error { ... }
func withDefault() (int, int) {
    // When using `or`, we declare 2 variables for the return values, not 3.
    q, r := divmod(5, 0) or 0, 0
    return q, r
}

// Ideas for array programming, SIMD, etc.

// Arrays (of the same size? what about slices?) be operated on like J or APL, only using wrapping operators.
[3]int{1, 2, 3} .+ [3]int{4, 5, 6} // [3]int{5, 7, 9}

// You can also operate on arrays with scalars. The scalar can be on either side of the operator.
[3]int{1, 2, 3} .+ 1 // [3]int{2, 3, 4}

// Other possible operators:
.= // broadcast assignment

// This gives some simple high-level ways to express SIMD operations. For more control, there should be a
// "simd" package

// Vector types
//
// Option 1: just use arrays. There's no need for separate vector types.
// Option 2: some possible type names:
int8x16
int16x8
int32x4
int64x2
float4
double2
// etc.

// What about matrix types?
int8x16x16
int16x8x8
int32x4x4
int64x2x2
float4x4
double2x2

// float4x4 and double2x2 are good, but the others are starting to get awkward.

// Option 3: some sort of simd type
simd[16]int8
simd[8]int16
simd[4]int32
// etc.
