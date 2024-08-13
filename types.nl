// Primitive types
int
int8
int16
int32
int64
int128

uint
uint8 (byte)
uint16
uint32
uint64
uint128

float16
float32
float64

uintptr

rune (int32)

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

// Go has complex64 and complex128. Not sure if we want those.


// An array. E.g. [5]int. Size is known statically.
[N]T

// A slice. E.g. []int. Length and capacity are stored in the slice.
// TODO:
// - Can slices be nil, like Go?
// - Should slices be mutable? I think so.
[]T

// A map. E.g. map[string]int. Alt: dict[string]int.
map[K]V

// A string. No assumed encoding. Source files are always UTF-8, and thus string literals are as well.
// TODO:
// - What should this be? In Go it's an immutable []byte. If you iterate you get runes.
// - How to deal with extended grapheme clusters?
// - What about C interop? Should string be null terminated? If the answer is no, then passing a string
//   to a function that takes a *C.char, needs to be an explicit copy.
string

// Composite types

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

// Tagged unions. Basically enum with associated values.
union {
    one int
    two (int, int) // Do we have tuples or is this just a syntactic construct? Let's start with the latter.
}

// Unions also support types without an explicit tag name. Equivalent to something like `Conn | File`.
//
// Does this promote methods? For simplicity, I think it probably shouldn't. We should stick with interfaces
// for that usecase. If we decide to promote, we can only promote methods that are present on every type in
// in the union. To be consistent with structs, we should only promote methods if all fields of the union
// are embedded.
union {
    Conn
    File
}

// Untagged union. Equivalent to C union. Alt: unsafe union, extern union.
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
// An open question: are these bridgable to C? We have to figure out whether strings are
// bridgable to C, which depends on whether they're null terminated or not.
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



// Type modifiers

// Move only types. Must be consumed or escaped.
type Fd (nocopy int)

// Anonymous underlying types don't use parens
type Fd nocopy struct {
    fd int
    // ...
}

// Types that contain or have underlying nocopy types must also be nocopy.

// automatically nocopy
type File Fd

// Must be explicitly declared nocopy
type File nocopy struct {
    fd Fd
}

// Invalid: Fd is already nocopy:
type File (nocopy Fd)

// Invalid, a nocopy struct must be nocopy:
type File struct {
    fd Fd
}

// Should you be able to use nocopy on variable declarations?
var i (nocopy int) = 5

// Pros:
// - You definitely need to be able to use nocopy on pointer declarations: `var p (nocopy *int)`. This creates
//   a nice symmetry.
// - It also creates a symmetry with type declarations. If type declarations are creating new names* for another
//   type, it stands that anything you can say in a type declaration, you should be able to say in a variable
//   declaration too.
//   *This is a lie. Type declarations create new types, not new names. Does this invalidate the above pro?
//
// Cons:
// - If you have a nocopy type like Fd, then `var f (nocopy Fd)` is redundant. Is it allowed? My hunch is no, but
//   I'm not sure.
// - It might create a false sense of symmetry with pointer variables. Borrowed, counted, weak, etc. can only be
//   used on pointers. Maybe we should have type modifiers and pointer modifiers and only allow type modifiers
//   in type declarations.

// If the answer to the above is yes, is this allowed? It's redundant.
var f (nocopy Fd)


// How are nocopy types printed in an IDE?
var i (nocopy int) // "(nocopy int)"
var f Fd           // "Fd (nocopy)". Or maybe just "Fd"? The former has less symmetry, but is more helpful, so I
                   // think it's better.


// Pointers

// A borrowed pointer. Panic on nil dereference. No action on drop. Cannot outlive the value it
// points to, which means unless we have some sort of support for lifetime dependencies, it can't
// be escaped. Passing to C is unsafe, but allowed. When passed to C, the programmer is responsible
// responsible for ensuring that it doesn't escape.
//
// TODO: How is it passed to C? Options:
// - automatically cast to (unsafe *T)
// - must be manually cast to (unsafe *T)
*int

// A pointer with unknown ownership. The programmer is responsible for managing the memory. Panics
// on nil dereference, but use-after-free is possible. C pointers are imported as unsafe.
(unsafe *int)
     
// An owned pointer. Just like other noncopy types, it must be consumed or escaped. Just like it's
// possible to borrow a nocopy type, it's possible to borrow a nocopy pointer.
//
// TODO: can you make multiple borrows of a nocopy *T? I would like the answer to be yes.
(nocopy *int)

(counted *int)  // A refcounted pointer. Alt: (rc *int), (strong *int). Never nil.
(weak *int)     // A weak pointer. Derived from a refcounted pointer. Becomes nil when the refcounted pointer is freed.

// owned and refcounted pointers have shorthands:
$*int // (nocopy *int)
#*int // (counted *int)


// For clarity, we'll write pointers out in longhand, but in general, shorthand is preferred.

// Just like structs which contain nocopy types, pointers to nocopy types must also be nocopy.
var f Fd
var p1 *Fd = &f          // error: Fd is nocopy so p1 must be nocopy
var p2 (nocopy *Fd) = &f // ok
p3 := &f                 // ok, typeof(p3) is (nocopy *Fd)


// Reference counting

// In pseudocode, a refcounted pointer is a pointer to a struct struct that's stored on the heap. The count
// is updated atomically.
struct {
    refcount int
    cleanup func(T)
    value T
}

// To make a refcounted pointer, use the rc builtin. The type passed to rc is copied into the refcounted struct.
// If it's a pointer, the pointer is copied. If it's a value, the value is copied.
func rc(v T) (counted *T)

// You can also supply a cleanup function that will be called when the refcount reaches 0.
func rc(v T, cleanup func(v T)) (counted *T)

// Refcounted pointers can own nocopy types, and must provide a cleanup function to consume the
// owned value.
func rc(v (nocopy T), deinit func(v (nocopy T))) (counted *T)

// You can also integrate external reference counted types by providing custom retain and release functions.
func rc(v T, retain func(v T), release func(v T)) (counted *T)

// A custom refcounted pointer has a different layout in memory:
struct {
    retain func(T)
    release func(T)
    value T
}

// You can make a weak reference using the weak builtin
func weak(p (counted *T)) (weak *T)

p := = rc(5) // typeof(p) is (counted *int)
w := weak(p)  // typeof(w) is (weak *int)

p := rc(fd, close) // typeof(p) is (counted *Fd)

// TODO: non-escapable types, lifetime dependencies, etc.
type Foo (noescape int)



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
