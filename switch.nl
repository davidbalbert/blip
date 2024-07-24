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
    // looks like a type declaration
    case x int:
        // ...
    case nil:
        // ...
}

type Ab enum {
    case a(int)
    case b(string)
}

var ab Ab
switch ab {
    case .a(i):
        // ...
    case .b(s):
        // ...
}

// Are error unions special, or can we have unnamed unions in general?

type hmm union {
    case int
    case error
}

type !int = hmm // the compiler does this automatically

var h hmm
switch h {
    // rhymes with `var x int`, which I like
    case x int:
        // ...
    case err error:
        // ...
}

// What about C unions?
type hmm untagged union {
    var i int
    var s string
}

// could also be spelled `extern union`.
