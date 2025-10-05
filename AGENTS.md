# Blip Compiler Project

## Overview
Blip is an incremental compiler written in Go using data-oriented design principles. The project is being built from the ground up, starting with minimal functionality and expanding incrementally.

## Current Status
- **Language**: Go
- **Design Philosophy**: Data-oriented design (CPU-friendly design, cache efficiency, memory layout optimization, keeping frequently used data together and small)
- **Current Functionality**: Parses arithmetic expressions from `.bl` source files and generates Mach-O object files directly with ARM64 machine code

## Architecture

### Pipeline
1. Lex `.bl` source file → token stream (`lex.Lex`)
2. Parse token stream → parse tree (`parse.Parse`)
3. Generate ARM64 machine code from parse tree (`codegen`)
4. Generate Mach-O object file → `macho.MachOGenerator`
5. Link into an ARM64 Mach-O executable with embedded ad-hoc code signature → `cmd/link`

## File Extensions & Workflow
- **Source**: `.bl` files
- **Generated**: `.o` object files (during testing)
- **Testing**: All compilation and execution testing is automated via Go tests

## Target Platforms
- **Current**: ARM64 on macOS
- **Planned**: x86-64 on macOS and Linux, ARM64 on Linux

## Key Commands
- **Test**: `go test ./...` (all tests are automated)
- **Development**: Do NOT build or run the compiler directly. All functionality should be tested through automated Go tests.

## Technical Notes
- Generates object files directly without assembly intermediate step
- Uses macOS system calls (sys_exit)
- Clean data transformations: source text → tokens → parse tree → machine code → object file
- Data-oriented approach: compact structs, cache-friendly memory layout, CPU-sympathetic design

## Code Style
- **Comments**: Leave no unnecessary comments. Only write comments where the code would be unclear without them. Comments should never repeat what the code does - they should explain context that's necessary to understand the code, and only if it's actually necessary.
- **Formatting**: Always auto-format code after making changes using the format_file tool.

## Future Plans
- Expand language syntax beyond arithmetic expressions
- Support multiple target architectures
- Maintain incremental, test-driven development approach

## Example
```
# test.bl
10 + 5 - 3

# Compilation process:
# 1. Lex: "10" "+" "5" "-" "3" → tokens
# 2. Parse: tokens → postorder parse tree
# 3. Codegen: parse tree → ARM64 machine code
# 4. Object file generation → test.bl.o

# Usage:
go run ./cmd/compile test.bl
go run ./cmd/link test.bl.o test
./test; echo $?  # outputs: 12
```
