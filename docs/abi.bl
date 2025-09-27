// Prior art:
//
// In Go, a package is a collection of source files. The Go compiler builds a single .o file for
// each package.
//
// TODO: Should we use this model, or should we have a one-to-one mapping from source files to .o
// files like C does? The former is convenient: a non-exported type or function can be referenced
// from any files in the package. In C, a static function can only be used in the file that
// declares it. That's a worse behavior, but it might be useful to have a one-to-one mapping for
// perfect C interop.

// For now, let's ignore C interop and use Go's definition of "package".

// Exported functions and types start with a capital letter
func Foo() {}
type Bar struct {}

// Unexported functions and types start with a lowercase letter
func foo() {}
type bar struct {}

// Exported types are exported both to other Blip packages as well as to C via a generated header.
// From this, it follows that any exported type or function must use the C ABI. Internal types which
// are only visible from within a single package, including closures, are free to use whatever calling
// convention they see fit. Because a package is always compiled as a unit, different source files in
// the same pacakge are guaranteed to agree on an ABI.
//
// TODO:
// 1. What about types and functions that can't be represented in C?
// 2. Can we find a nice way to make methods callable from C?
