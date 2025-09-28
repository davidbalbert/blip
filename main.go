package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davidbalbert/blip/lex"
	"github.com/davidbalbert/blip/macho"
	"github.com/davidbalbert/blip/obj/arm64"
)

func codegen(content []byte) []byte {
	gen := arm64.Generator{}

	lexer := lex.NewLexer(content)
	operation := '+'
	firstValue := true

	for {
		token := lexer.NextToken()
		tokenType := token.Type()

		if tokenType == lex.TokEOF {
			break
		}

		switch tokenType {
		case lex.TokInt:
			start := token.Pos()
			end := start
			for end < uint32(len(content)) && isDigit(content[end]) {
				end++
			}
			value, _ := strconv.Atoi(string(content[start:end]))
			if firstValue {
				gen.MovImm(0, value) // mov x0, #value
				firstValue = false
			} else if operation == '+' {
				gen.AddImm(0, 0, value) // add x0, x0, #value
			} else {
				gen.SubImm(0, 0, value) // sub x0, x0, #value
			}
		case lex.TokAdd:
			operation = '+'
		case lex.TokSub:
			operation = '-'
		case lex.TokInvalid:
		}
	}

	gen.MovImm16(16, 1) // mov x16, #1 (sys_exit)
	gen.SVC(0x80)       // svc #0x80

	return gen.Bytes()
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

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

	machineCode := codegen(content)

	machoGen := macho.NewMachOGenerator()
	machoGen.SetTextData(machineCode)
	objData := machoGen.Generate()

	err = os.WriteFile(objFile, objData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", objFile)
}
