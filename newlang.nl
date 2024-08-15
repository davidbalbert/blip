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

// A string. No assumed encoding. Source files are always UTF-8, and thus string literals are as well. Strings
// are null-terminated, and can be bridged to C without copying. Like slices, strings store their length and
// capacity, which can be queried in constant time.
// TODO:
// - What should this be? In Go it's an immutable []byte. If you iterate you get runes.
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

// You can also use type modifiers in variable declarations.
var i (nocopy int) = 5

var f (nocopy Fd) // error: Fd is already nocopy

// How are nocopy types printed in an IDE?
var i (nocopy int) // "(nocopy int)"
var f Fd           // "Fd (nocopy)"

// An example:
var x int = 5
y := x         // x is copied here
x++
print(x, y)    // 6, 5

var z (nocopy int) = 5
w := z         // z is moved to w
z++            // error: z was moved to w
print(z, w)    // error: z was moved to w

// Pointers

// There are 5 pointer types: borrowed, owned, unsafe, reference counted, and weak. All pointers except
// reference counted pointers can be nil. Dereferencing a nil pointer is defined behavior – it panics.
// All pointers besides unsafe pointers prevent have temporal safety – they prevent use-after-free.

// A borrowed pointer. Panic on nil dereference. No action on drop. Cannot outlive the value it
// points to. Passing to C is unsafe, but allowed. When passed to C, the programmer is responsible
// for ensuring that if the pointer is escaped, it doesn't outlive its referent.
//
// TODO: How is it passed to C? Options:
// - automatically cast to (unsafe *T)
// - must be manually cast to (unsafe *T)
*int

// A pointer that owns the memory it points to. A piece of memory can have exactly one owner. Just like
// other nocopy types, performing an assignment moves, rather than copies the pointer. This ensures that
// the "single owner" invariant holds. Unlike other nocopy types, you are allowed to allow a nocopy pointer
// to reach the end of its scope without consuming or escaping it. At the end of its scope, the pointer
// is automatically dropped (consumed), and the memory is freed.
//
// TODO: explain this better.
// Owned pointers to nocopy types (nocopy *(nocopy T)) still needs to be consumed or escaped.
(nocopy *int)

// A pointer with unknown ownership. The programmer is responsible for managing the memory. Panics
// on nil dereference, but use-after-free is possible. C pointers are imported as unsafe.
(unsafe *int)
     
(counted *int)  // A refcounted pointer. Alt: (rc *int), (strong *int). Never nil.
(weak *int)     // A weak pointer. Derived from a refcounted pointer. Becomes nil when the refcounted pointer is freed.

// owned and refcounted pointers have shorthands:
$*int // (nocopy *int)
#*int // (counted *int)

// For clarity, we'll write pointers out in longhand, but in general, shorthand is preferred.


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


// Pointer conversions

// A stack value can be borrowed multiple times
var x int
var p1 *int = &x          // borrowed pointer.
p2 := &x                  // multiple borrows of a stack allocated value is ok. typeof(p2) is *int

// A local variable can be moved into an owned pointer, in which case the original variable is consumed. Escape
// analysis is performed. If the variable escapes (e.g. is returned) or is consumed by a function, it will be
// heap allocated.
var x int
var p (nocopy *int) = &x // may be stack or heap allocated
print(x)                 // error: x was moved to p. It doesn't matter that x is copyable.

// Another syntax for the above
var x int
p := (nocopy *int)(&x)

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
var p (nocopy *int) = &5 
p := &5
var p *int = &5 // error: you cannot borrow a literal

// To force heap allocation, use make.
p := make(int) // typeof(p) is (nocopy *int). The int is on the heap.


// You can borrow an owned pointer
var x int
var p1 (nocopy *int) = &x
var p2 *int = p1 // implicit borrow. typeof(p2) is *int


// On the other hand, if a variable has been borrowed, it can't be moved into an owned pointer.
var x int
p1 := &x                   // typeof p1 is *int
var p2 (nocopy *int) = &x  // error: p1 borrows x, so p2 can't own it.









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
