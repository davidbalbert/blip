# obj - Machine Code Generation

## Overview
The obj package contains target-specific machine code generators. Currently supports ARM64 architecture.

## Architecture

### ARM64 Generator (`obj/arm64`)

#### Data Structures
```go
type Generator struct {
    Instructions []Insn  // Generated machine code instructions
}

type Insn uint32  // ARM64 instruction encoding
```

#### Key Functions
- `MovImm(reg, imm)` - MOV register, immediate value
- `AddImm(rd, rn, imm)` - ADD register, register, immediate  
- `SubImm(rd, rn, imm)` - SUB register, register, immediate
- `MovImm16(reg, imm)` - MOV register, immediate (for syscalls)
- `SVC(imm)` - Supervisor call (system call)
- `Bytes()` - Convert instructions to raw machine code bytes

#### Instruction Encoding
- Uses hardcoded ARM64 instruction encodings
- Little-endian byte order
- 32-bit instructions (4 bytes each)
- Direct binary generation without assembler

#### Example Usage
```go
gen := arm64.Generator{}
gen.MovImm(0, 42)        // mov x0, #42
gen.MovImm16(16, 1)      // mov x16, #1 (sys_exit)
gen.SVC(0x80)            // svc #0x80
machineCode := gen.Bytes()
```

## Target Architectures
- **Current**: ARM64 on macOS
- **Planned**: x86-64, ARM64 on Linux

## Technical Notes
- CPU-friendly design with compact data structures
- Cache-efficient memory layout
- No external assembler dependencies
- Direct machine code generation for maximum control
