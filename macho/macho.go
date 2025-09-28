package macho

import (
	"encoding/binary"
)

const (
	MH_MAGIC_64           = 0xfeedfacf
	MH_OBJECT             = 1
	CPU_TYPE_ARM64        = 0x0100000c
	CPU_SUBTYPE_ARM64_ALL = 0

	LC_SEGMENT_64    = 0x19
	LC_SYMTAB        = 0x2
	LC_BUILD_VERSION = 0x32

	PLATFORM_MACOS = 1

	S_REGULAR                = 0x0
	S_ATTR_PURE_INSTRUCTIONS = 0x80000000
	S_ATTR_SOME_INSTRUCTIONS = 0x400
)

type MachOGenerator struct {
	Text []byte
}

func (m *MachOGenerator) Generate() []byte {
	// Calculate offsets
	textOffset := uint32(232) // header + commands (including build version)
	symOffset := textOffset + uint32(len(m.Text))
	strOffset := symOffset + 16 // one symbol * 16 bytes
	strSize := uint32(7)        // "_main\0\0\0" padded to 8-byte boundary

	result := make([]byte, 0, 1024)

	// Mach-O header (32 bytes)
	header := make([]byte, 32)
	binary.LittleEndian.PutUint32(header[0:], MH_MAGIC_64)           // magic
	binary.LittleEndian.PutUint32(header[4:], CPU_TYPE_ARM64)        // cputype
	binary.LittleEndian.PutUint32(header[8:], CPU_SUBTYPE_ARM64_ALL) // cpusubtype
	binary.LittleEndian.PutUint32(header[12:], MH_OBJECT)            // filetype
	binary.LittleEndian.PutUint32(header[16:], 3)                    // ncmds (segment + symtab + build_version)
	binary.LittleEndian.PutUint32(header[20:], 200)                  // sizeofcmds (152 + 24 + 24)
	binary.LittleEndian.PutUint32(header[24:], 0)                    // flags
	binary.LittleEndian.PutUint32(header[28:], 0)                    // reserved
	result = append(result, header...)

	// LC_SEGMENT_64 command (152 bytes = 72 + 80 for section)
	segCmd := make([]byte, 152)
	binary.LittleEndian.PutUint32(segCmd[0:], LC_SEGMENT_64)                     // cmd
	binary.LittleEndian.PutUint32(segCmd[4:], 152)                               // cmdsize
	copy(segCmd[8:24], []byte("__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")) // segname
	binary.LittleEndian.PutUint64(segCmd[24:], 0)                                // vmaddr
	binary.LittleEndian.PutUint64(segCmd[32:], uint64(len(m.Text)))              // vmsize
	binary.LittleEndian.PutUint64(segCmd[40:], uint64(textOffset))               // fileoff
	binary.LittleEndian.PutUint64(segCmd[48:], uint64(len(m.Text)))              // filesize
	binary.LittleEndian.PutUint32(segCmd[56:], 5)                                // maxprot (VM_PROT_READ | VM_PROT_EXECUTE)
	binary.LittleEndian.PutUint32(segCmd[60:], 5)                                // initprot
	binary.LittleEndian.PutUint32(segCmd[64:], 1)                                // nsects
	binary.LittleEndian.PutUint32(segCmd[68:], 0)                                // flags
	// Section header for __text (80 bytes) - embedded in segment command
	copy(segCmd[72:88], []byte("__text\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))                            // sectname
	copy(segCmd[88:104], []byte("__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))                           // segname
	binary.LittleEndian.PutUint64(segCmd[104:], 0)                                                           // addr
	binary.LittleEndian.PutUint64(segCmd[112:], uint64(len(m.Text)))                                         // size
	binary.LittleEndian.PutUint32(segCmd[120:], textOffset)                                                  // offset
	binary.LittleEndian.PutUint32(segCmd[124:], 2)                                                           // align (4-byte alignment)
	binary.LittleEndian.PutUint32(segCmd[128:], 0)                                                           // reloff
	binary.LittleEndian.PutUint32(segCmd[132:], 0)                                                           // nreloc
	binary.LittleEndian.PutUint32(segCmd[136:], S_REGULAR|S_ATTR_PURE_INSTRUCTIONS|S_ATTR_SOME_INSTRUCTIONS) // flags
	binary.LittleEndian.PutUint32(segCmd[140:], 0)                                                           // reserved1
	binary.LittleEndian.PutUint32(segCmd[144:], 0)                                                           // reserved2
	binary.LittleEndian.PutUint32(segCmd[148:], 0)                                                           // reserved3
	result = append(result, segCmd...)

	// LC_SYMTAB command (24 bytes)
	symCmd := make([]byte, 24)
	binary.LittleEndian.PutUint32(symCmd[0:], LC_SYMTAB)  // cmd
	binary.LittleEndian.PutUint32(symCmd[4:], 24)         // cmdsize
	binary.LittleEndian.PutUint32(symCmd[8:], symOffset)  // symoff
	binary.LittleEndian.PutUint32(symCmd[12:], 1)         // nsyms
	binary.LittleEndian.PutUint32(symCmd[16:], strOffset) // stroff
	binary.LittleEndian.PutUint32(symCmd[20:], strSize)   // strsize
	result = append(result, symCmd...)

	// LC_BUILD_VERSION command (24 bytes)
	buildCmd := make([]byte, 24)
	binary.LittleEndian.PutUint32(buildCmd[0:], LC_BUILD_VERSION) // cmd
	binary.LittleEndian.PutUint32(buildCmd[4:], 24)               // cmdsize
	binary.LittleEndian.PutUint32(buildCmd[8:], PLATFORM_MACOS)   // platform
	binary.LittleEndian.PutUint32(buildCmd[12:], 0x000B0000)      // minos (macOS 11.0)
	binary.LittleEndian.PutUint32(buildCmd[16:], 0x000B0000)      // sdk (macOS 11.0)
	binary.LittleEndian.PutUint32(buildCmd[20:], 0)               // ntools
	result = append(result, buildCmd...)

	// Pad to fileoff
	for len(result) < int(textOffset) {
		result = append(result, 0)
	}

	// Text data
	result = append(result, m.Text...)

	// Symbol table (16 bytes per symbol)
	symbol := make([]byte, 16)
	binary.LittleEndian.PutUint32(symbol[0:], 1) // n_strx (string table offset)
	symbol[4] = 0x0F                             // n_type (N_SECT | N_EXT)
	symbol[5] = 1                                // n_sect (section 1)
	binary.LittleEndian.PutUint16(symbol[6:], 0) // n_desc
	binary.LittleEndian.PutUint64(symbol[8:], 0) // n_value (address)
	result = append(result, symbol...)

	// String table
	stringTable := []byte("\x00_main\x00\x00\x00") // null + "_main" + padding
	result = append(result, stringTable...)

	return result
}
