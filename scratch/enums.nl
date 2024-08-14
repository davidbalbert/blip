// this has no specified underlying type, so it's not bridgable to C.
// or perhaps enums are compiler asigned ints by default?
type op enum {
    add
    sub
    mul
    div
}

// If you want to specify a type, you need to do it in the first case. This is also not bridgable to C,
// because we don't know how int and C.int are related. int32 would be bridgable.
type op enum {
    add int = 1
    sub = 5
    mul = 10
    div = 20
}

// if you leave out values, and you're using a numeric type, successive values are incremented

type op enum {
    add int // 0
    sub     // 1
    mul     // 2
    div     // 3
    mod = 128
    pow     // 129
}

// In Go, constants without an assignment expression, repeat the most recently specified expression
// by textual substitution.

const (
    a = 1
    b       // = 1
    c       // = 1
)

// iota is therefore necessary, because it increments in each successive constant declaration

const (
    a = iota // (0)
    b        // = iota (1)
    c        // = iota (2)
)

// We don't need iota in enums, because the following should be a compile error. All enum cases
// must be unique.

type op enum {
    add int = 1
    sub = 1
    mul = 1
    div = 1
}

// For enums that are numbers (or perhaps enums that are ints?), we can default to incrementing, like specified above:

type op enum {
    add int = 5
    sub    // 6
    mul    // 7
    div    // 8
}

// For strings, we can default to the name of the case:

type op enum {
    add string // "add"
    sub        // "sub"
    mul = "Mul"
    div        // "div"
}

// Perhaps that's it for enums without associated values. They can either have no underlying type, ints, or strings.
// Alternatively, they could always have a default underlying type of int.

// How about using them?

type op enum {
    add int
    sub
    mul
    div
}

// variables
x := op.add
var x op = op.add // maybe op:add? 
var x op = add

// given this function
func printInt(x int)

// do we need to cast, or do we have implicit conversions?
printInt((int)x) // explicit
printInt(int(x)) // explicit, alt syntax
printInt(x)      // implicit

// if we don't have a variable, we always have to qualfiy with `op`:

printInt(int(op.add)) // explicit
printInt(op.add)      // implicit

// otoh, with
func printOp(x op)

// we can always use the unqualified name
printOp(add)

// definitely have to cast here:
func printInt32(x int32)

printInt32(int32(x))
printInt32((int32)x) // alt syntax
printInt32(x) // error


// Could these get exported to C? At minimum, we'd have to export everything:

type Op enum {
    Add int64
    Sub
    Mul
    Div
}

func Perform(op Op)

// And we could export them as int constants.

typedef Op int64_t; // remember, this doesn't actually create a new type in C, just a new name. But its still
                    // useful for documentation purposes.

Op OpAdd;
Op OpSub;
Op OpMul;
Op OpDiv;

// But what would the signature of Perform be? I think the only possible declaration is
void Perform(Op op);

// But this means that Perform could receive values that it doesn't expect. Two choices:
// 1. All `switches` over enums are well defined if they receive an unexpected value. There's an implicit
//    `default` case that panics.
// 2. The C version of Perform is a wrapper around the Newlang version that Panics.
// 3. We can't export any functions that take enums as arguments.

// I think the first option is the best. It's useful to export enums to C, and undefined behavior is bad.


// Associated values

type Insn enum {
    case add(int, int)
    case inc(int)
}

// or maybe
type Insn enum {
    case (int, int) add
    case (int) add
}
