// basic switch

var x int
switch x {
    case 1:
        // ...
    case 2:
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
    case int
    case float
}

// similar to ?int, but only supports types, not values.
var x val
switch x {
    case i int:
        // ...
    case f float:
        // ...
}

// similar to ?int, but not exactly the same: you can't do a `case x int` and unwrap things that way.
type optint enum {
    case some (int)
    case none
}

var x optint
switch x {
    case some(i):
        // ...
    case none:
        // ...
}



type Op enum {
    case add
    case inc
}

var op Op
switch op {
    case add:
        // ...
    case inc:
        // ...
}


type Insn enum {
    case add(int, int)
    case inc(int)
}

var insn Insn
switch insn {
    case add(i, j):
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
