package main

import "os"

type sourceFile struct {
	imports      []string
	declarations []declaration
}

type declarationType int

const (
	dtFunction declarationType = iota
	dtVariable
)

type declaration struct {
	kind declarationType
	val  literal
}

type literalType int

const (
	ltInt literalType = iota
	ltString
)

type literal struct {
	kind literalType
	val  any
}

func parse(contents string) sourceFile {
	// parse the source code into an AST
	// and print the AST to stdout
}

func main() {
	// read each file from the command line, open it, and parse the source
	// code into an AST

	for _, filename := range os.Args[1:] {
		contents, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		parse(string(contents))
	}
}
