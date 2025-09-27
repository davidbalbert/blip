# Blip Compiler Project

## Overview
Blip is an incremental compiler written in Go using data-oriented design principles. The project is being built from the ground up, starting with minimal functionality and expanding incrementally.

## Current Status
- **Language**: Go
- **Design Philosophy**: Data-oriented design (CPU-friendly design, cache efficiency, memory layout optimization, keeping frequently used data together and small)
- **Current Functionality**: Parses single integer from `.bl` source files and generates ARM64 assembly that exits with that value

## Architecture

### Data Structures
```go
type Program struct {
    ExitCode int  // Currently just a single integer
}

type Assembly struct {
    Instructions []string  // Generated assembly instructions
}
```

### Pipeline
1. Parse `.bl` source file → `Program` struct
2. Generate ARM64 assembly → `Assembly` struct  
3. Write assembly to temporary directory (avoids Go assembler conflicts)
4. External toolchain: `as` (assembler) → `ld` (linker) → executable

## File Extensions & Workflow
- **Source**: `.bl` files
- **Generated**: `.s` assembly files (in temp directories)
- **Build chain**: `test.bl` → `test.s` → `test.o` → `test` (executable)

## Target Platforms
- **Current**: ARM64 on macOS
- **Planned**: x86-64 on macOS and Linux, ARM64 on Linux

## Key Commands
- **Compile**: `go run . <source.bl>` or `go build -o blip && ./blip <source.bl>`
- **Full build**: `as file.s -o file.o && ld file.o -o executable -lSystem -syslibroot $(xcrun --show-sdk-path) -e _main`

## Technical Notes
- Uses temporary directories for assembly generation to avoid conflicts with Go's built-in assembler
- Assembly generated uses macOS system calls (sys_exit)
- Single-file design (everything in `main.go`) for now
- Clean data transformations: source text → Program → Assembly → file output
- Data-oriented approach: compact structs, cache-friendly memory layout, CPU-sympathetic design

## Code Style
- **Comments**: Leave no unnecessary comments. Only write comments where the code would be unclear without them. Comments should never repeat what the code does - they should explain context that's necessary to understand the code, and only if it's actually necessary.

## Future Plans
- Expand language syntax beyond single integers
- Support multiple target architectures
- Eventually generate object files directly (bypass external assembler)
- Maintain incremental, test-driven development approach

## Example
```
# test.bl
42

# Generated assembly (test.s)
.global _main
.align 2

_main:
    mov x0, #42
    mov x16, #1
    svc #0x80
```
