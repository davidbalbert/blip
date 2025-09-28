package main

import (
	"fmt"
	"os"

	"github.com/davidbalbert/blip/lex"
	"github.com/davidbalbert/blip/macho"
	"github.com/davidbalbert/blip/obj"
	"github.com/davidbalbert/blip/parse"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source-file>\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := os.Args[1]

	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		os.Exit(1)
	}

	objFile := sourceFile + ".o"
	tokens := lex.Lex(content)
	nodes := parse.Parse(tokens)
	insns := obj.Codegen(nodes, tokens, content)

	machoGen := macho.MachOGenerator{Text: insns}
	objData := machoGen.Generate()

	err = os.WriteFile(objFile, objData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", objFile)
}
