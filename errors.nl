// errors

// One option is Go style
func foo() (int, error)

// But let's explore. Can we do something like a result type, but not use generics?

func foo() !int

func optional() {
    var x ?int = try? foo()
}

func defaultValue() {
    var x int = try? foo() ?? 0
}

// Would this imply that !int has to be a type everywhere, or can we get away with having
// ! just appear in the return type?
func directReturn() !int {
    return foo() // alt: return try! foo()
}

func shortCircuit() {
    var x int = try foo()
    // only run if foo errors, return otherwise
}

// What about when we have to return a value with no error?
func shortCircuit2() int {
    var x int = try foo() else -1
    // or
    var x int = try foo() else return -1
    // or
    var x int = try foo() else {
        return -1
    }
}

// What about when we need to get a handle on the error?
func withError() {
    var x int = try foo() recover err in {
        // handle error.
        // must return (like guard in Swift).
        return
    }
}

// I don't love the above. Is there something simpler?

func withError() {
    if try x = foo() {
        // typeof(x) == int
    } else recover err {
        // typeof(err) == error
    }
}

// Other spellings:
if x = try foo() {
    // ...
} else err = recover {
    // ...
}

try x = foo() {
    // ...
} catch err {
    // ...
}

// err could be a default name
try x = foo() {
    // ...
} catch {
    // err is in scope
}

// I don't love any of these. Go's error handling is so simple. If we go for a "result" type, it should
// still feel easy and breezy like Go. It would be great if we could just get rid of try/catch. Those terms
// feel heavy. recover isn't so bad though.

func recovery() {
    var x int = try foo() recover {
        // err is in scope
        // you must return
        return
    }
}

// The above is a bit odd because try/recover looks like an expression, but most other things aren't going to
// be expressions.

interface error {
    func Error() string
}

// I like this syntax for switch/match.
switch foo() {
    case x int:
        // ...
    case err error:
        // ...
}

