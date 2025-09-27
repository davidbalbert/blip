package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Data-oriented design: focus on data and transformations
// Program represents the minimal AST - just an integer for now
type Program struct {
	ExitCode int
}

// Assembly represents generated assembly instructions
type Assembly struct {
	Instructions []string
}

// parseProgram reads source file and extracts the integer
func parseProgram(filename string) (Program, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Program{}, err
	}

	// For now, the entire program is just one integer
	text := strings.TrimSpace(string(content))
	exitCode, err := strconv.Atoi(text)
	if err != nil {
		return Program{}, err
	}

	return Program{ExitCode: exitCode}, nil
}

// generateARM64 transforms program data into ARM64 assembly data
func generateARM64(program Program) Assembly {
	var instructions []string

	// ARM64 assembly to exit with the given code
	instructions = append(instructions, ".global _main")
	instructions = append(instructions, ".align 2")
	instructions = append(instructions, "")
	instructions = append(instructions, "_main:")
	instructions = append(instructions, "    mov x0, #"+strconv.Itoa(program.ExitCode))
	instructions = append(instructions, "    mov x16, #1") // sys_exit
	instructions = append(instructions, "    svc #0x80")   // system call

	return Assembly{Instructions: instructions}
}

// writeAssembly outputs assembly data to file
func writeAssembly(filename string, assembly Assembly) error {
	content := strings.Join(assembly.Instructions, "\n") + "\n"
	return os.WriteFile(filename, []byte(content), 0644)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source-file>\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := os.Args[1]

	// Read and parse source
	program, err := parseProgram(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Create temp directory for compilation
	tempDir, err := os.MkdirTemp("", "blip-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Temp dir error: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Generate assembly in temp directory
	baseName := strings.TrimSuffix(filepath.Base(sourceFile), ".bl")
	asmFile := filepath.Join(tempDir, baseName+".s")
	
	asm := generateARM64(program)
	err = writeAssembly(asmFile, asm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", asmFile)
}
