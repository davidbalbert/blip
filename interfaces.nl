// Let's start with a sketch of Go's interfaces.

alias any interface{}
func Println(...any) {}

// Any type can implement the empty interface. int and *int are different types, but both implement any.
x := 5
var y any = x // interface{}(int)

p := &x
var q any = p // interface{}(*int)

Println(y, q) // ok

// You can't overload functions. Methods are just functions with a special name (Type.Method), and a special place to declare
// the receiver, which is just the first argument to the function.
//
// This means, for a given method, we have to pick whether the receiver is a pointer or a value. We can't have both.
type File struct{}
func (f File) Read([]byte) (int, error) { return 0, nil }
func (f *File) Read([]byte) (int, error) { return 0, nil } // error: File.Read is already declared above.

// This means that while File and *File can implement any, only one of them can implement Reader (or any other non-empty interface).

type Reader interface {
    Read([]byte) (int, error)
}

// How to adapt this to our various pointer types?

// A method can be defined with any receiver type, including *T, (nocopy *T), (counted *T) and (weak *T). It may only be declared for one of them.
//
// A method defined on *T is a borrow. All pointer types can be implicitly converted into a borrow.

type File struct{}
func (f *File) Read([]byte) (int, error) { return 0, nil }
var f *File; f.Read(nil)        // ok
var f $*File; f.Read(nil)       // ok
var f #*File; f.Read(nil)       // ok
var f (weak *File); f.Read(nil) // ok

// Value receivers are also implicitly converted into borrows.
var f File; f.Read(nil) // ok

// You can define a method with a receiver of type (counted *T) or (weak *T), but it's not that useful.


// The big question is how to handle consuming receivers. First let's redefine File to be nocopy. This isn't strictly necessary to have a consuming
// Close method – making a owned pointer to a variable will move the variable even if it's copyable. But making File nocopy is what will force us
// to call f.Close() in a defer.
type File (nocopy struct{})

// Then we'll declare Close as a consuming method by using a (nocopy *File) receiver.
func (f (nocopy *File)) Close() error { return nil }

// The question is, how do we define the Closer interface? Here's a start:
type Closer nocopy interface {
    Close() error
}

// But what about defining ReadCloser?
type ReadCloser interface {
    Reader
    Closer
}

// Normally something that contains a nocopy thing must also be nocopy. But that's clearly not right here. Read is borrowing, but Close is consuming.
// So we need syntax to attach it to the method:
type Closer interface {
    nocopy Close() error
}

type Closer interface {
    eat Close() error
}

type Closer interface {
    consume Close() error
}

type Closer interface {
    move Close() error
}

type Closer interface {
    Close() error nocopy
}

// 
