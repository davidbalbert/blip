// basic switch

var x int
switch x {
    case 1:
        // ...
    case 2, 3, 4:
        // ...
    default:
        // ...
}

var x ?int
switch x {
    // looks like a type declaration. That's nice.
    case x int:
        // ...
    case nil:
        // ...
}

type val union {
    int
    float
}

// similar to ?int, but only supports types, not values.
var x val
switch x {
    case i int:
        // ...
    case f float:
        // ...
}

type val untion {
    n int
    d float
}

// Don't love this
var x val
switch x {
    case n:
        // ...
    case d:
        // ...
}

// similar to ?int, but not exactly the same: you can't do a `case x int` and unwrap things that way.
type optint enum {
    some (int)
    none
}

var x optint
switch x {
    case some (i):
        // ...
    case none:
        // ...
}



type Op enum {
    add
    inc
}

var op Op
switch op {
    case add:
        // ...
    case inc:
        // ...
}

type Insn enum int64 {
    add
    inc
}

type Insn enum {
    add a, b int, d bool
    inc
}

var insn Insn
switch insn {
    case add a, b, d:
        // ...
    case inc(i):
        // ...
}


// or if using this
type Insn enum {
    case (int, int) add
    case (int) add
}

// then maybe

var insn Insn
switch insn {
    case (i, j) add:
        // ...
    case (i) inc:
        // ...
}
