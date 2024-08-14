use io, os

func main() {
    for i, arg in os.Args {
        io.Print(arg)

        if i < len(os.Args) - 1 {
            io.Print(" ")
        }
    }
    io.Println()
}