# Blip Compiler Project

## Instructions to the agent (YOU MUST FOLLOW THESE)
- Always run tests after making changes.

## Overview
Blip is an incremental compiler written in Go using data-oriented design principles. The project is being built from the ground up, starting with minimal functionality and expanding incrementally.

## Current Status
- **Language**: Go
- **Design Philosophy**: Data-oriented design (CPU-friendly design, cache efficiency, memory layout optimization, keeping frequently used data together and small)
- **Current Functionality**: Parses arithmetic expressions from `.bl` source files and generates Mach-O (macOS) or ELF (Linux) object files directly with ARM64 machine code

## Architecture

### Code Structure
- **`cmd/compile`** - Thin wrapper that calls `internal/cmd/compile.Main()`
- **`cmd/link`** - Thin wrapper that calls `internal/cmd/link.Main()`
- **`internal/cmd/compile`** - Compiler logic (importable for tests)
- **`internal/cmd/link`** - Linker logic (importable for tests)
- **`lex/`** - Lexer
- **`parse/`** - Parser
- **`obj/`** - Code generation
- **`macho/`** - Mach-O file generation (macOS)
- **`elf/`** - ELF file generation (Linux)
- **`target/`** - Target OS detection
- **`test/integration/`** - End-to-end integration tests

### Pipeline
1. Lex `.bl` source file → token stream (`lex.Lex`)
2. Parse token stream → parse tree (`parse.Parse`)
3. Generate ARM64 machine code from parse tree (`obj.Codegen`)
4. Generate object file → `macho.Generate` (macOS) or `elf.Generate` (Linux)
5. Link into executable → `internal/link.Link` (Mach-O with ad-hoc code signature on macOS, ELF on Linux)

## File Extensions & Workflow
- **Source**: `.bl` files
- **Generated**: `.o` object files
- **Testing**: Integration tests in `test/integration/` test the complete pipeline by executing the test binary as both compiler and linker

## Target Platforms
- **Current**: ARM64 on macOS and Linux
- **Planned**: x86-64 on macOS and Linux

## Environment Variables
- **`BLOS`** - Target OS for cross-compilation: `"macos"` or `"linux"` (defaults to host OS)

## Key Commands
- **Test**: `go test ./...` (all tests are automated)
- **Build Compiler**: `go build ./cmd/compile`
- **Build Linker**: `go build ./cmd/link`
- **Run Compiler**: `go run ./cmd/compile <source-file>`
- **Run Linker**: `go run ./cmd/link <object-file> <output-executable>`

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
```bash
# test.bl
10 + 5 - 3

# Compilation process:
# 1. Lex: "10" "+" "5" "-" "3" → tokens
# 2. Parse: tokens → postorder parse tree
# 3. Codegen: parse tree → ARM64 machine code
# 4. Object file generation → test.bl.o

# Usage:
go run ./cmd/compile test.bl       # Creates test.bl.o
go run ./cmd/link test.bl.o test   # Creates executable
./test; echo $?                    # Outputs: 12
```

## Testing Architecture
Integration tests use the cmd/go pattern:
- Tests import `internal/cmd/compile` and `internal/cmd/link`
- Test binary acts as both compiler and linker via `BLIP_TEST_MODE` env var
- `BLIP_TEST_MODE=compile` → runs `compile.Main()`
- `BLIP_TEST_MODE=link` → runs `link.Main()`
- Tests exec the test binary to exercise complete Main() paths
- Changes to compiler or linker code invalidate test cache automatically

## Cross-Platform Testing (test/crosstest)
Runs ARM64 Linux binaries in a QEMU VM for cross-platform testing on macOS.

### Components
- **`run.sh`** - Host script that launches QEMU, sends binary, captures output
- **`kbuild/`** - Builds minimal Linux kernel with BusyBox initramfs
- **`kbuild/init.sh`** - VM init script that receives and runs the binary

### Protocol
Uses virtio-serial pipes for communication:
- Port 1: stdin, Port 2: stdout, Port 3: stderr, Port 4: control
- Control channel sends: `[8-byte total_size][8-byte program_size][program][8-byte args_size][null-separated args]`
- Returns exit code over control channel

### Usage
```bash
cd test/crosstest
./run.sh <program> [args...]
```

### Notes
- Size fields use leading zeros; init.sh strips them to avoid octal interpretation
- Always use `timeout` when testing run.sh to avoid hangs
- Kernel image built via Docker: `make` in kbuild/
