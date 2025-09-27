use io, os

func square(x: Int?) -> Int? {
    if let x {
        return x*x
    }
    return nil
}


func square(x: Int?) -> Int? {
    guard let x else {
        return nil
    }
    return x * x
}

const cat1 = {path : byte[:]
    var buf : byte[32*1024]
    var f, n

    f = std.open(path, std.Oread) ? else |e|
        std.die("could not read {}", e)

    while true
        match std.read(f, buf[:])
        | `std.Ok n:    std.write(std.out, buf[:n])
        | `std.Err e:   std.die("could not read: {}", e)
        ;;
    ;;
    std.close(f)
}

func cat2(path: String) {
    guard let fd = open(path) else {
        return
    }

    while let s = read(fd) {
        print(s)
    }
}

struct Array<T> {
    class Storage {
        // ...
    }
    var storage: Storage

    mutating func append(x: T) {
        if !isKnownUniquelyReferenced(&self) {
            storage = copyStorage()
        }
    }
}

struct Foo: ~Copyable {
}

type Fd nocopy struct {
}

type Fd moved struct {
}

type Node counted struct {
}

nocopy struct Fd {
}

counted struct Node {
}

// THIS IS GOOD!
func free(x : !*int) {
    p = (unsafe *void) x;
}

template<typename T>
T&& move(T& t) { return t; }

class X {
    int v;
    X(X& other ) { v = other.v; }
    X(X&& other) { v = other.v; other.v = 0;}
}


func bad(path string) Fd {
    // AN INTERESTING SPELLING OF LIFETIMES
    fd1 := open(path) in .
    fd2 := open("/dev/null") in .

    if rand() % 2 == 0 {
        close(fd2)
        return makeBufferedReader(fd1)
    } else {
        close(fd1)
        return fd2
    }
}

func makeBufferedReader(fd Fd) BufferedReader {
    return BufferedReader{fd}
}

func cat1(path : string)
{
    std::unique_ptr<File> f = os.open(path)
    X x = std::move(y);

    for c : f.bychunk(io.BufSize) {
        os.Out.write(c)
    } else {
        os.Err.print("error reading %s: %r\n", path)
    }


}

func main() {

}
