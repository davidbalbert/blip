# macho - Mach-O Object File Generation

## Overview
The macho package generates Mach-O object files directly from machine code, bypassing the need for external assemblers.

## Architecture

### API

#### Key Functions
- `Generate(text []byte) []byte` - Generate Mach-O object file from machine code
- `GenerateExecutable(textData []byte, identifier string) []byte` - Generate Mach-O executable

#### Mach-O Structure (Object Files)
- **Header**: 64-bit Mach-O magic, ARM64 CPU type, object file type
- **Load Commands**:
  - `LC_SEGMENT_64`: __TEXT segment with __text section
  - `LC_SYMTAB`: Symbol table with _main symbol
  - `LC_BUILD_VERSION`: Platform version info for macOS
- **Data**: Raw machine code, symbol table, string table

#### Example Usage
```go
// Generate object file
objData := macho.Generate(machineCode)
os.WriteFile("output.o", objData, 0644)

// Generate executable
exeData := macho.GenerateExecutable(machineCode, "myprogram")
os.WriteFile("myprogram", exeData, 0755)
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
