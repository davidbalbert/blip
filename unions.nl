// Are error unions special, or do we have general tagged unions?
type interror union {
    case int
    case error
}

type !int = interror // the compiler does this automatically

var ie interror
switch ie {
    // rhymes with `var x int`, which I like
    case x int:
        // ...
    case err error:
        // ...
}

// What about C unions?
type hmm untagged union {
    i int
    s string
}

// could also be spelled `extern union`.

// alternatively, any union that uses var instead of case could be considered untagged, but I think that's too
// easy to miss.
