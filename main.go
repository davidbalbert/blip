package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davidbalbert/blip/lex"
)

func codegen(content []byte) []string {
	var insns []string
	insns = append(insns, ".global _main")
	insns = append(insns, ".align 2")
	insns = append(insns, "")
	insns = append(insns, "_main:")

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
			value := string(content[start:end])
			if firstValue {
				insns = append(insns, "    mov x0, #"+value)
				firstValue = false
			} else if operation == '+' {
				insns = append(insns, "    add x0, x0, #"+value)
			} else {
				insns = append(insns, "    sub x0, x0, #"+value)
			}
		case lex.TokAdd:
			operation = '+'
		case lex.TokSub:
			operation = '-'
		case lex.TokInvalid:
		}
	}

	insns = append(insns, "    mov x16, #1")
	insns = append(insns, "    svc #0x80")

	return insns
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

	tempDir := "/tmp"
	asmFile := filepath.Join(tempDir, sourceFile+".s")

	insns := codegen(content)
	asm := strings.Join(insns, "\n") + "\n"
	err = os.WriteFile(asmFile, []byte(asm), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", asmFile)
}
