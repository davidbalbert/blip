package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var Magic = []byte{0x7f, 'E', 'L', 'F'}

const (
	ELFCLASS64  = 2
	ELFDATA2LSB = 1
	EV_CURRENT  = 1

	ET_REL  = 1
	ET_EXEC = 2

	EM_AARCH64 = 183

	SHT_NULL     = 0
	SHT_PROGBITS = 1
	SHT_SYMTAB   = 2
	SHT_STRTAB   = 3

	SHF_ALLOC     = 0x2
	SHF_EXECINSTR = 0x4

	STB_GLOBAL = 1
	STT_FUNC   = 2

	PT_LOAD = 1

	PF_X = 0x1
	PF_R = 0x4
)

func Generate(text []byte) []byte {
	ehdrSize := uint64(64)
	shdrSize := uint64(64)
	symSize := uint64(24)

	shstrtab := []byte("\x00.text\x00.symtab\x00.strtab\x00.shstrtab\x00")
	strtab := []byte("\x00_start\x00")

	textOffset := ehdrSize
	textSize := uint64(len(text))

	symtabOffset := textOffset + textSize
	symtabSize := symSize

	strtabOffset := symtabOffset + symtabSize
	strtabSize := uint64(len(strtab))

	shstrtabOffset := strtabOffset + strtabSize
	shstrtabSize := uint64(len(shstrtab))

	shdrOffset := (shstrtabOffset + shstrtabSize + 7) &^ 7
	numSections := uint16(5)

	result := make([]byte, 0, shdrOffset+uint64(numSections)*shdrSize)

	// ELF header (64 bytes)
	ehdr := make([]byte, ehdrSize)
	copy(ehdr[0:4], Magic)                                     // e_ident[EI_MAG0..EI_MAG3]: magic number
	ehdr[4] = ELFCLASS64                                       // e_ident[EI_CLASS]: 64-bit
	ehdr[5] = ELFDATA2LSB                                      // e_ident[EI_DATA]: 2's complement, little-endian
	ehdr[6] = EV_CURRENT                                       // e_ident[EI_VERSION]: current version
	ehdr[7] = 0                                                // e_ident[EI_OSABI]: ELFOSABI_NONE
	binary.LittleEndian.PutUint16(ehdr[16:], ET_REL)           // e_type: relocatable object
	binary.LittleEndian.PutUint16(ehdr[18:], EM_AARCH64)       // e_machine: ARM64
	binary.LittleEndian.PutUint32(ehdr[20:], EV_CURRENT)       // e_version: current
	binary.LittleEndian.PutUint64(ehdr[24:], 0)                // e_entry: entry point (none for object)
	binary.LittleEndian.PutUint64(ehdr[32:], 0)                // e_phoff: program header offset (none)
	binary.LittleEndian.PutUint64(ehdr[40:], shdrOffset)       // e_shoff: section header offset
	binary.LittleEndian.PutUint32(ehdr[48:], 0)                // e_flags: processor flags
	binary.LittleEndian.PutUint16(ehdr[52:], uint16(ehdrSize)) // e_ehsize: ELF header size
	binary.LittleEndian.PutUint16(ehdr[54:], 0)                // e_phentsize: program header entry size
	binary.LittleEndian.PutUint16(ehdr[56:], 0)                // e_phnum: number of program headers
	binary.LittleEndian.PutUint16(ehdr[58:], uint16(shdrSize)) // e_shentsize: section header entry size
	binary.LittleEndian.PutUint16(ehdr[60:], numSections)      // e_shnum: number of section headers
	binary.LittleEndian.PutUint16(ehdr[62:], 4)                // e_shstrndx: section name string table index
	result = append(result, ehdr...)

	result = append(result, text...)

	// Symbol table entry (24 bytes)
	sym := make([]byte, symSize)
	binary.LittleEndian.PutUint32(sym[0:], 1)         // st_name: offset into string table
	sym[4] = (STB_GLOBAL << 4) | STT_FUNC             // st_info: global function
	sym[5] = 1                                        // st_other: default visibility
	binary.LittleEndian.PutUint16(sym[6:], 0)         // st_shndx: undefined section
	binary.LittleEndian.PutUint64(sym[8:], 0)         // st_value: symbol value
	binary.LittleEndian.PutUint64(sym[16:], textSize) // st_size: symbol size
	result = append(result, sym...)

	result = append(result, strtab...)
	result = append(result, shstrtab...)

	for uint64(len(result)) < shdrOffset {
		result = append(result, 0)
	}

	// Null section header (required as first entry)
	shdrNull := make([]byte, shdrSize)
	result = append(result, shdrNull...)

	// .text section header
	shdrText := make([]byte, shdrSize)
	binary.LittleEndian.PutUint32(shdrText[0:], 1)                       // sh_name: offset into shstrtab
	binary.LittleEndian.PutUint32(shdrText[4:], SHT_PROGBITS)            // sh_type: program data
	binary.LittleEndian.PutUint64(shdrText[8:], SHF_ALLOC|SHF_EXECINSTR) // sh_flags: alloc + exec
	binary.LittleEndian.PutUint64(shdrText[16:], 0)                      // sh_addr: virtual address
	binary.LittleEndian.PutUint64(shdrText[24:], textOffset)             // sh_offset: file offset
	binary.LittleEndian.PutUint64(shdrText[32:], textSize)               // sh_size: section size
	binary.LittleEndian.PutUint32(shdrText[40:], 0)                      // sh_link: no linked section
	binary.LittleEndian.PutUint32(shdrText[44:], 0)                      // sh_info: no extra info
	binary.LittleEndian.PutUint64(shdrText[48:], 4)                      // sh_addralign: 4-byte alignment
	binary.LittleEndian.PutUint64(shdrText[56:], 0)                      // sh_entsize: no fixed-size entries
	result = append(result, shdrText...)

	// .symtab section header
	shdrSymtab := make([]byte, shdrSize)
	binary.LittleEndian.PutUint32(shdrSymtab[0:], 7)             // sh_name: offset into shstrtab
	binary.LittleEndian.PutUint32(shdrSymtab[4:], SHT_SYMTAB)    // sh_type: symbol table
	binary.LittleEndian.PutUint64(shdrSymtab[8:], 0)             // sh_flags: none
	binary.LittleEndian.PutUint64(shdrSymtab[16:], 0)            // sh_addr: virtual address
	binary.LittleEndian.PutUint64(shdrSymtab[24:], symtabOffset) // sh_offset: file offset
	binary.LittleEndian.PutUint64(shdrSymtab[32:], symtabSize)   // sh_size: section size
	binary.LittleEndian.PutUint32(shdrSymtab[40:], 3)            // sh_link: associated strtab index
	binary.LittleEndian.PutUint32(shdrSymtab[44:], 1)            // sh_info: first non-local symbol
	binary.LittleEndian.PutUint64(shdrSymtab[48:], 8)            // sh_addralign: 8-byte alignment
	binary.LittleEndian.PutUint64(shdrSymtab[56:], symSize)      // sh_entsize: symbol entry size
	result = append(result, shdrSymtab...)

	// .strtab section header
	shdrStrtab := make([]byte, shdrSize)
	binary.LittleEndian.PutUint32(shdrStrtab[0:], 15)            // sh_name: offset into shstrtab
	binary.LittleEndian.PutUint32(shdrStrtab[4:], SHT_STRTAB)    // sh_type: string table
	binary.LittleEndian.PutUint64(shdrStrtab[8:], 0)             // sh_flags: none
	binary.LittleEndian.PutUint64(shdrStrtab[16:], 0)            // sh_addr: virtual address
	binary.LittleEndian.PutUint64(shdrStrtab[24:], strtabOffset) // sh_offset: file offset
	binary.LittleEndian.PutUint64(shdrStrtab[32:], strtabSize)   // sh_size: section size
	binary.LittleEndian.PutUint32(shdrStrtab[40:], 0)            // sh_link: no linked section
	binary.LittleEndian.PutUint32(shdrStrtab[44:], 0)            // sh_info: no extra info
	binary.LittleEndian.PutUint64(shdrStrtab[48:], 1)            // sh_addralign: 1-byte alignment
	binary.LittleEndian.PutUint64(shdrStrtab[56:], 0)            // sh_entsize: no fixed-size entries
	result = append(result, shdrStrtab...)

	// .shstrtab section header
	shdrShstrtab := make([]byte, shdrSize)
	binary.LittleEndian.PutUint32(shdrShstrtab[0:], 23)              // sh_name: offset into shstrtab
	binary.LittleEndian.PutUint32(shdrShstrtab[4:], SHT_STRTAB)      // sh_type: string table
	binary.LittleEndian.PutUint64(shdrShstrtab[8:], 0)               // sh_flags: none
	binary.LittleEndian.PutUint64(shdrShstrtab[16:], 0)              // sh_addr: virtual address
	binary.LittleEndian.PutUint64(shdrShstrtab[24:], shstrtabOffset) // sh_offset: file offset
	binary.LittleEndian.PutUint64(shdrShstrtab[32:], shstrtabSize)   // sh_size: section size
	binary.LittleEndian.PutUint32(shdrShstrtab[40:], 0)              // sh_link: no linked section
	binary.LittleEndian.PutUint32(shdrShstrtab[44:], 0)              // sh_info: no extra info
	binary.LittleEndian.PutUint64(shdrShstrtab[48:], 1)              // sh_addralign: 1-byte alignment
	binary.LittleEndian.PutUint64(shdrShstrtab[56:], 0)              // sh_entsize: no fixed-size entries
	result = append(result, shdrShstrtab...)

	return result
}

func GenerateExecutable(textData []byte) []byte {
	baseAddr := uint64(0x400000)
	ehdrSize := uint64(64)
	phdrSize := uint64(56)

	textOffset := (ehdrSize + phdrSize + 0xfff) &^ 0xfff
	textSize := uint64(len(textData))
	entryAddr := baseAddr + textOffset

	fileSize := textOffset + textSize

	result := make([]byte, 0, fileSize)

	// ELF header (64 bytes)
	ehdr := make([]byte, ehdrSize)
	copy(ehdr[0:4], Magic)                                     // e_ident[EI_MAG0..EI_MAG3]: magic number
	ehdr[4] = ELFCLASS64                                       // e_ident[EI_CLASS]: 64-bit
	ehdr[5] = ELFDATA2LSB                                      // e_ident[EI_DATA]: little-endian
	ehdr[6] = EV_CURRENT                                       // e_ident[EI_VERSION]: current version
	ehdr[7] = 0                                                // e_ident[EI_OSABI]: ELFOSABI_NONE
	binary.LittleEndian.PutUint16(ehdr[16:], ET_EXEC)          // e_type: executable
	binary.LittleEndian.PutUint16(ehdr[18:], EM_AARCH64)       // e_machine: ARM64
	binary.LittleEndian.PutUint32(ehdr[20:], EV_CURRENT)       // e_version: current
	binary.LittleEndian.PutUint64(ehdr[24:], entryAddr)        // e_entry: entry point address
	binary.LittleEndian.PutUint64(ehdr[32:], ehdrSize)         // e_phoff: program header offset
	binary.LittleEndian.PutUint64(ehdr[40:], 0)                // e_shoff: section header offset (none)
	binary.LittleEndian.PutUint32(ehdr[48:], 0)                // e_flags: processor flags
	binary.LittleEndian.PutUint16(ehdr[52:], uint16(ehdrSize)) // e_ehsize: ELF header size
	binary.LittleEndian.PutUint16(ehdr[54:], uint16(phdrSize)) // e_phentsize: program header entry size
	binary.LittleEndian.PutUint16(ehdr[56:], 1)                // e_phnum: number of program headers
	binary.LittleEndian.PutUint16(ehdr[58:], 0)                // e_shentsize: section header entry size
	binary.LittleEndian.PutUint16(ehdr[60:], 0)                // e_shnum: number of section headers
	binary.LittleEndian.PutUint16(ehdr[62:], 0)                // e_shstrndx: section name string table index
	result = append(result, ehdr...)

	// Program header (56 bytes)
	phdr := make([]byte, phdrSize)
	binary.LittleEndian.PutUint32(phdr[0:], PT_LOAD)              // p_type: loadable segment
	binary.LittleEndian.PutUint32(phdr[4:], PF_R|PF_X)            // p_flags: read + execute
	binary.LittleEndian.PutUint64(phdr[8:], textOffset)           // p_offset: file offset
	binary.LittleEndian.PutUint64(phdr[16:], baseAddr+textOffset) // p_vaddr: virtual address
	binary.LittleEndian.PutUint64(phdr[24:], baseAddr+textOffset) // p_paddr: physical address
	binary.LittleEndian.PutUint64(phdr[32:], textSize)            // p_filesz: size in file
	binary.LittleEndian.PutUint64(phdr[40:], textSize)            // p_memsz: size in memory
	binary.LittleEndian.PutUint64(phdr[48:], 0x1000)              // p_align: page alignment
	result = append(result, phdr...)

	for uint64(len(result)) < textOffset {
		result = append(result, 0)
	}

	result = append(result, textData...)

	return result
}

func ExtractTextSection(data []byte) ([]byte, error) {
	if len(data) < 64 {
		return nil, fmt.Errorf("file too short for ELF header")
	}

	if !bytes.Equal(data[0:4], Magic) {
		return nil, fmt.Errorf("not an ELF file")
	}

	shoff := binary.LittleEndian.Uint64(data[40:48])
	shentsize := binary.LittleEndian.Uint16(data[58:60])
	shnum := binary.LittleEndian.Uint16(data[60:62])
	shstrndx := binary.LittleEndian.Uint16(data[62:64])

	if shoff == 0 || shnum == 0 {
		return nil, fmt.Errorf("no section headers")
	}

	shstrOff := shoff + uint64(shstrndx)*uint64(shentsize)
	if shstrOff+uint64(shentsize) > uint64(len(data)) {
		return nil, fmt.Errorf("shstrtab section header out of bounds")
	}
	shstrFileOff := binary.LittleEndian.Uint64(data[shstrOff+24 : shstrOff+32])
	shstrSize := binary.LittleEndian.Uint64(data[shstrOff+32 : shstrOff+40])
	if shstrFileOff+shstrSize > uint64(len(data)) {
		return nil, fmt.Errorf("shstrtab content out of bounds")
	}
	shstrtab := data[shstrFileOff : shstrFileOff+shstrSize]

	for i := uint16(0); i < shnum; i++ {
		off := shoff + uint64(i)*uint64(shentsize)
		if off+uint64(shentsize) > uint64(len(data)) {
			return nil, fmt.Errorf("section header out of bounds")
		}

		nameIdx := binary.LittleEndian.Uint32(data[off : off+4])
		if uint64(nameIdx) >= uint64(len(shstrtab)) {
			continue
		}

		var name []byte
		for j := nameIdx; j < uint32(len(shstrtab)) && shstrtab[j] != 0; j++ {
			name = append(name, shstrtab[j])
		}

		if string(name) == ".text" {
			sectOff := binary.LittleEndian.Uint64(data[off+24 : off+32])
			sectSize := binary.LittleEndian.Uint64(data[off+32 : off+40])
			if sectOff+sectSize > uint64(len(data)) {
				return nil, fmt.Errorf(".text section extends beyond file")
			}
			return data[sectOff : sectOff+sectSize], nil
		}
	}

	return nil, fmt.Errorf("no .text section found")
}
