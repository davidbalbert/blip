# macho - Mach-O Object File Generation

## Overview
The macho package generates Mach-O object files directly from machine code, bypassing the need for external assemblers.

## Architecture

### MachOGenerator

#### Data Structures
```go
type MachOGenerator struct {
    textData []byte  // Raw machine code bytes
}
```

#### Key Functions
- `NewMachOGenerator()` - Create new generator
- `SetTextData(data []byte)` - Set machine code to embed
- `Generate() []byte` - Generate complete Mach-O object file

#### Mach-O Structure
- **Header**: 64-bit Mach-O magic, ARM64 CPU type, object file type
- **Load Commands**:
  - `LC_SEGMENT_64`: __TEXT segment with __text section
  - `LC_SYMTAB`: Symbol table with _main symbol
  - `LC_BUILD_VERSION`: Platform version info for macOS
- **Data**: Raw machine code, symbol table, string table

#### Example Usage
```go
machoGen := macho.NewMachOGenerator()
machoGen.SetTextData(machineCode)
objData := machoGen.Generate()
os.WriteFile("output.o", objData, 0644)
```

## Technical Details
- Generates minimal but complete Mach-O object files
- Hardcoded offsets and sizes for simplicity
- Includes proper platform version information
- Creates exportable _main symbol for linker
- Direct binary generation without external tools

## Platform Support
- **Current**: ARM64 on macOS with platform version 11.0
- **Format**: 64-bit Mach-O object files
- **Linking**: Compatible with system linker (`ld`)
