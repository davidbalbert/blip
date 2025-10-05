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
		fmt.Fprintf(os.Stderr, "usage: %s <source-file>\n", os.Args[0])
		os.Exit(2)
	}

	sourceFile := os.Args[1]
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		os.Exit(1)
	}

	tokens := lex.Lex(content)
	nodes := parse.Parse(tokens)
	text := obj.Codegen(nodes, tokens, content)

	objData := macho.Generate(text)

	objFile := sourceFile + ".o"
	if err := os.WriteFile(objFile, objData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
}
