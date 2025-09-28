package main

import (
	"fmt"
	"os"

	"github.com/davidbalbert/blip/compile"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <source-file>\n", os.Args[0])
		os.Exit(2)
	}

	sourceFile := os.Args[1]
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}

	objData := compile.Compile(content)

	objFile := sourceFile + ".o"
	if err := os.WriteFile(objFile, objData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
}
