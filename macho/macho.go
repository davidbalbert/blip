package macho

import (
	"crypto/sha256"
	"encoding/binary"
)

const (
	MH_MAGIC_64           = 0xfeedfacf
	MH_OBJECT             = 1
	MH_EXECUTE            = 2
	CPU_TYPE_ARM64        = 0x0100000c
	CPU_SUBTYPE_ARM64_ALL = 0

	LC_SEGMENT_64     = 0x19
	LC_SYMTAB         = 0x2
	LC_DYSYMTAB       = 0xb
	LC_DYLD_INFO_ONLY = 0x80000022
	LC_LOAD_DYLINKER  = 0xe
	LC_LOAD_DYLIB     = 0xc
	LC_UUID           = 0x1b
	LC_BUILD_VERSION  = 0x32
	LC_MAIN           = 0x80000028
	LC_CODE_SIGNATURE = 0x1d

	PLATFORM_MACOS = 1

	S_REGULAR                = 0x0
	S_ATTR_PURE_INSTRUCTIONS = 0x80000000
	S_ATTR_SOME_INSTRUCTIONS = 0x400

	VM_PROT_READ    = 1
	VM_PROT_WRITE   = 2
	VM_PROT_EXECUTE = 4
)

// Create a Mach-O object file from machine code
func Generate(text []byte) []byte {
	// Calculate offsets
	textOffset := uint32(232) // header + commands (including build version)
	symOffset := textOffset + uint32(len(text))
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
	binary.LittleEndian.PutUint64(segCmd[32:], uint64(len(text)))                // vmsize
	binary.LittleEndian.PutUint64(segCmd[40:], uint64(textOffset))               // fileoff
	binary.LittleEndian.PutUint64(segCmd[48:], uint64(len(text)))                // filesize
	binary.LittleEndian.PutUint32(segCmd[56:], 5)                                // maxprot (VM_PROT_READ | VM_PROT_EXECUTE)
	binary.LittleEndian.PutUint32(segCmd[60:], 5)                                // initprot
	binary.LittleEndian.PutUint32(segCmd[64:], 1)                                // nsects
	binary.LittleEndian.PutUint32(segCmd[68:], 0)                                // flags
	// Section header for __text (80 bytes) - embedded in segment command
	copy(segCmd[72:88], []byte("__text\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))                            // sectname
	copy(segCmd[88:104], []byte("__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))                           // segname
	binary.LittleEndian.PutUint64(segCmd[104:], 0)                                                           // addr
	binary.LittleEndian.PutUint64(segCmd[112:], uint64(len(text)))                                           // size
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
	result = append(result, text...)

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

// Create Mach-O executable from machine code
func GenerateExecutable(textData []byte, identifier string) []byte {
	baseAddr := uint64(0x100000000)
	textSize := uint64(len(textData))
	pageSize := uint64(0x4000) // 16KB pages on ARM64

	// Calculate offsets
	pagezeroCmd := uint32(72)     // LC_SEGMENT_64 for __PAGEZERO
	textSegmentCmd := uint32(152) // LC_SEGMENT_64 + section for __TEXT
	linkeditCmd := uint32(72)     // LC_SEGMENT_64 for __LINKEDIT
	dyldInfoCmdSize := uint32(48)
	symtabCmdSize := uint32(24)
	dysymtabCmdSize := uint32(80)
	uuidCmdSize := uint32(24)
	buildCmdSize := uint32(24)
	mainCmdSize := uint32(24)
	dylinkerCmdSize := uint32(32)
	dylibPath := "/usr/lib/libSystem.B.dylib"
	dylibNameLen := len(dylibPath) + 1
	dylibCmdSize := uint32((24 + dylibNameLen + 7) &^ 7)
	codeSigCmdSize := uint32(16)
	totalCmdsSize := pagezeroCmd + textSegmentCmd + linkeditCmd + dyldInfoCmdSize + symtabCmdSize + dysymtabCmdSize + uuidCmdSize + buildCmdSize + mainCmdSize + dylinkerCmdSize + dylibCmdSize + codeSigCmdSize

	textOffset := uint32(0x2f8) // Standard entry offset used by ld-generated binaries
	textVMAddr := baseAddr + uint64(textOffset)

	result := make([]byte, 0, int(textOffset)+len(textData))

	// Mach-O header
	header := make([]byte, 32)
	binary.LittleEndian.PutUint32(header[0:], MH_MAGIC_64)
	binary.LittleEndian.PutUint32(header[4:], CPU_TYPE_ARM64)
	binary.LittleEndian.PutUint32(header[8:], 0) // cpusubtype
	binary.LittleEndian.PutUint32(header[12:], MH_EXECUTE)
	binary.LittleEndian.PutUint32(header[16:], 12)
	binary.LittleEndian.PutUint32(header[20:], totalCmdsSize)
	binary.LittleEndian.PutUint32(header[24:], 0x200085)
	binary.LittleEndian.PutUint32(header[28:], 0) // reserved
	result = append(result, header...)

	// LC_SEGMENT_64 for __PAGEZERO
	pagezeroSegCmd := make([]byte, 72)
	binary.LittleEndian.PutUint32(pagezeroSegCmd[0:], LC_SEGMENT_64)
	binary.LittleEndian.PutUint32(pagezeroSegCmd[4:], pagezeroCmd)
	copy(pagezeroSegCmd[8:24], []byte("__PAGEZERO\x00\x00\x00\x00\x00\x00"))
	binary.LittleEndian.PutUint64(pagezeroSegCmd[24:], 0)        // vmaddr
	binary.LittleEndian.PutUint64(pagezeroSegCmd[32:], baseAddr) // vmsize
	binary.LittleEndian.PutUint64(pagezeroSegCmd[40:], 0)        // fileoff
	binary.LittleEndian.PutUint64(pagezeroSegCmd[48:], 0)        // filesize
	binary.LittleEndian.PutUint32(pagezeroSegCmd[56:], 0)        // maxprot
	binary.LittleEndian.PutUint32(pagezeroSegCmd[60:], 0)        // initprot
	binary.LittleEndian.PutUint32(pagezeroSegCmd[64:], 0)        // nsects
	binary.LittleEndian.PutUint32(pagezeroSegCmd[68:], 0)        // flags
	result = append(result, pagezeroSegCmd...)

	// LC_SEGMENT_64 for __TEXT with __text section
	textSegCmd := make([]byte, 152)
	binary.LittleEndian.PutUint32(textSegCmd[0:], LC_SEGMENT_64)
	binary.LittleEndian.PutUint32(textSegCmd[4:], textSegmentCmd)
	copy(textSegCmd[8:24], []byte("__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	binary.LittleEndian.PutUint64(textSegCmd[24:], baseAddr)                     // vmaddr
	binary.LittleEndian.PutUint64(textSegCmd[32:], pageSize)                     // vmsize (page aligned)
	binary.LittleEndian.PutUint64(textSegCmd[40:], 0)                            // fileoff
	binary.LittleEndian.PutUint64(textSegCmd[48:], pageSize)                     // filesize (page aligned)
	binary.LittleEndian.PutUint32(textSegCmd[56:], VM_PROT_READ|VM_PROT_EXECUTE) // maxprot
	binary.LittleEndian.PutUint32(textSegCmd[60:], VM_PROT_READ|VM_PROT_EXECUTE) // initprot
	binary.LittleEndian.PutUint32(textSegCmd[64:], 1)                            // nsects
	binary.LittleEndian.PutUint32(textSegCmd[68:], 0)                            // flags
	// __text section (80 bytes)
	copy(textSegCmd[72:88], []byte("__text\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	copy(textSegCmd[88:104], []byte("__TEXT\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	binary.LittleEndian.PutUint64(textSegCmd[104:], textVMAddr)                                                  // addr
	binary.LittleEndian.PutUint64(textSegCmd[112:], textSize)                                                    // size
	binary.LittleEndian.PutUint32(textSegCmd[120:], textOffset)                                                  // offset
	binary.LittleEndian.PutUint32(textSegCmd[124:], 2)                                                           // align (4-byte alignment)
	binary.LittleEndian.PutUint32(textSegCmd[128:], 0)                                                           // reloff
	binary.LittleEndian.PutUint32(textSegCmd[132:], 0)                                                           // nreloc
	binary.LittleEndian.PutUint32(textSegCmd[136:], S_REGULAR|S_ATTR_PURE_INSTRUCTIONS|S_ATTR_SOME_INSTRUCTIONS) // flags
	binary.LittleEndian.PutUint32(textSegCmd[140:], 0)                                                           // reserved1
	binary.LittleEndian.PutUint32(textSegCmd[144:], 0)                                                           // reserved2
	binary.LittleEndian.PutUint32(textSegCmd[148:], 0)                                                           // reserved3
	result = append(result, textSegCmd...)

	linkeditOffset := len(result)
	linkeditSeg := make([]byte, linkeditCmd)
	binary.LittleEndian.PutUint32(linkeditSeg[0:], LC_SEGMENT_64)
	binary.LittleEndian.PutUint32(linkeditSeg[4:], linkeditCmd)
	copy(linkeditSeg[8:24], []byte("__LINKEDIT\x00\x00\x00\x00\x00\x00"))
	binary.LittleEndian.PutUint64(linkeditSeg[24:], baseAddr+pageSize) // vmaddr
	binary.LittleEndian.PutUint64(linkeditSeg[32:], pageSize)          // vmsize
	binary.LittleEndian.PutUint64(linkeditSeg[40:], uint64(pageSize))  // fileoff
	// filesize filled after signature generation
	binary.LittleEndian.PutUint32(linkeditSeg[56:], VM_PROT_READ)
	binary.LittleEndian.PutUint32(linkeditSeg[60:], VM_PROT_READ)
	result = append(result, linkeditSeg...)

	dyldInfoCmd := make([]byte, dyldInfoCmdSize)
	binary.LittleEndian.PutUint32(dyldInfoCmd[0:], LC_DYLD_INFO_ONLY)
	binary.LittleEndian.PutUint32(dyldInfoCmd[4:], dyldInfoCmdSize)
	result = append(result, dyldInfoCmd...)

	symtabCmd := make([]byte, symtabCmdSize)
	binary.LittleEndian.PutUint32(symtabCmd[0:], LC_SYMTAB)
	binary.LittleEndian.PutUint32(symtabCmd[4:], symtabCmdSize)
	result = append(result, symtabCmd...)

	dysymtabCmd := make([]byte, dysymtabCmdSize)
	binary.LittleEndian.PutUint32(dysymtabCmd[0:], LC_DYSYMTAB)
	binary.LittleEndian.PutUint32(dysymtabCmd[4:], dysymtabCmdSize)
	result = append(result, dysymtabCmd...)

	uuidCmd := make([]byte, uuidCmdSize)
	binary.LittleEndian.PutUint32(uuidCmd[0:], LC_UUID)
	binary.LittleEndian.PutUint32(uuidCmd[4:], uuidCmdSize)
	uuidBytes := sha256.Sum256(textData)
	copy(uuidCmd[8:], uuidBytes[:16])
	result = append(result, uuidCmd...)

	buildCmd := make([]byte, buildCmdSize)
	binary.LittleEndian.PutUint32(buildCmd[0:], LC_BUILD_VERSION)
	binary.LittleEndian.PutUint32(buildCmd[4:], buildCmdSize)
	binary.LittleEndian.PutUint32(buildCmd[8:], PLATFORM_MACOS)
	binary.LittleEndian.PutUint32(buildCmd[12:], 0x000B0000) // minos
	binary.LittleEndian.PutUint32(buildCmd[16:], 0x000B0000) // sdk
	binary.LittleEndian.PutUint32(buildCmd[20:], 0)          // ntools
	result = append(result, buildCmd...)

	mainCmd := make([]byte, mainCmdSize)
	binary.LittleEndian.PutUint32(mainCmd[0:], LC_MAIN)
	binary.LittleEndian.PutUint32(mainCmd[4:], mainCmdSize)
	binary.LittleEndian.PutUint64(mainCmd[8:], uint64(textOffset)) // entryoff (relative to __TEXT)
	binary.LittleEndian.PutUint64(mainCmd[16:], 0)
	result = append(result, mainCmd...)

	dylinkerCmd := make([]byte, dylinkerCmdSize)
	binary.LittleEndian.PutUint32(dylinkerCmd[0:], LC_LOAD_DYLINKER)
	binary.LittleEndian.PutUint32(dylinkerCmd[4:], dylinkerCmdSize)
	binary.LittleEndian.PutUint32(dylinkerCmd[8:], 12) // offset to dylinker path
	copy(dylinkerCmd[12:], []byte("/usr/lib/dyld\x00\x00\x00\x00\x00\x00\x00"))
	result = append(result, dylinkerCmd...)

	dylibCmd := make([]byte, dylibCmdSize)
	binary.LittleEndian.PutUint32(dylibCmd[0:], LC_LOAD_DYLIB)
	binary.LittleEndian.PutUint32(dylibCmd[4:], dylibCmdSize)
	binary.LittleEndian.PutUint32(dylibCmd[8:], 24)          // name offset
	binary.LittleEndian.PutUint32(dylibCmd[12:], 2)          // timestamp
	binary.LittleEndian.PutUint32(dylibCmd[16:], 0x054C0000) // current version 1356.0.0
	binary.LittleEndian.PutUint32(dylibCmd[20:], 0x00010000) // compatibility version 1.0.0
	copy(dylibCmd[24:], dylibPath)
	result = append(result, dylibCmd...)

	codeSigCmdOffset := len(result)
	codeSigCmd := make([]byte, 16)
	binary.LittleEndian.PutUint32(codeSigCmd[0:], LC_CODE_SIGNATURE)
	binary.LittleEndian.PutUint32(codeSigCmd[4:], codeSigCmdSize)
	// dataoff/datasize filled after signature generation
	result = append(result, codeSigCmd...)

	// Pad to the start of the __text section
	for len(result) < int(textOffset) {
		result = append(result, 0)
	}

	// Text data
	result = append(result, textData...)

	for len(result) < int(pageSize) {
		result = append(result, 0)
	}

	sigOffset := len(result)
	sigLen := len(buildCodeSignature(result, identifier, textSize))
	binary.LittleEndian.PutUint32(result[codeSigCmdOffset+8:], uint32(sigOffset))
	binary.LittleEndian.PutUint32(result[codeSigCmdOffset+12:], uint32(sigLen))
	binary.LittleEndian.PutUint64(result[linkeditOffset+48:], uint64(sigLen))

	sigData := buildCodeSignature(result, identifier, textSize)
	result = append(result, sigData...)

	return result
}

func buildCodeSignature(code []byte, identifier string, execSize uint64) []byte {
	const (
		csMagicEmbeddedSignature = 0xfade0cc0
		csMagicCodeDirectory     = 0xfade0c02
		csSlotCodeDirectory      = 0
		codeDirVersion           = 0x20400
		codeDirFlags             = 0x20002
		hashTypeSHA256           = 2
		hashSize                 = 32
		pageSizeLog2             = 12
		codeDirHeaderSize        = 88
	)

	if identifier == "" {
		identifier = "blip"
	}

	pageSize := 1 << pageSizeLog2
	codeLimit := len(code)
	nCodeSlots := 0
	if codeLimit > 0 {
		nCodeSlots = (codeLimit + pageSize - 1) / pageSize
	}

	identBytes := []byte(identifier)
	identLen := len(identBytes) + 1 // include null terminator
	hashOffset := codeDirHeaderSize + identLen
	if rem := hashOffset & 3; rem != 0 {
		hashOffset += 4 - rem
	}
	totalHashes := nCodeSlots * hashSize
	codeDirLen := hashOffset + totalHashes
	codeDir := make([]byte, codeDirLen)

	binary.BigEndian.PutUint32(codeDir[0:], csMagicCodeDirectory)
	binary.BigEndian.PutUint32(codeDir[4:], uint32(codeDirLen))
	binary.BigEndian.PutUint32(codeDir[8:], codeDirVersion)
	binary.BigEndian.PutUint32(codeDir[12:], codeDirFlags)
	binary.BigEndian.PutUint32(codeDir[16:], uint32(hashOffset))
	binary.BigEndian.PutUint32(codeDir[20:], uint32(codeDirHeaderSize))
	binary.BigEndian.PutUint32(codeDir[24:], 0) // nSpecialSlots
	binary.BigEndian.PutUint32(codeDir[28:], uint32(nCodeSlots))
	binary.BigEndian.PutUint32(codeDir[32:], uint32(codeLimit))
	codeDir[36] = hashSize
	codeDir[37] = hashTypeSHA256
	codeDir[38] = 0
	codeDir[39] = pageSizeLog2
	binary.BigEndian.PutUint32(codeDir[40:], 0) // spare2
	binary.BigEndian.PutUint32(codeDir[44:], 0) // scatterOffset
	binary.BigEndian.PutUint32(codeDir[48:], 0) // teamIDOffset
	binary.BigEndian.PutUint32(codeDir[52:], 0) // spare3
	binary.BigEndian.PutUint64(codeDir[56:], uint64(codeLimit))
	binary.BigEndian.PutUint64(codeDir[64:], 0) // execSegBase
	binary.BigEndian.PutUint64(codeDir[72:], execSize)
	binary.BigEndian.PutUint64(codeDir[80:], 1) // execSegFlags (main binary)

	copy(codeDir[codeDirHeaderSize:], identBytes)
	// trailing byte already zero

	if nCodeSlots > 0 {
		hashRegion := codeDir[hashOffset:]
		for i := 0; i < nCodeSlots; i++ {
			start := i * pageSize
			end := start + pageSize
			if end > codeLimit {
				end = codeLimit
			}
			sum := sha256.Sum256(code[start:end])
			copy(hashRegion[i*hashSize:(i+1)*hashSize], sum[:])
		}
	}

	superBlob := make([]byte, 12+8+len(codeDir))
	binary.BigEndian.PutUint32(superBlob[0:], csMagicEmbeddedSignature)
	binary.BigEndian.PutUint32(superBlob[4:], uint32(len(superBlob)))
	binary.BigEndian.PutUint32(superBlob[8:], 1)
	binary.BigEndian.PutUint32(superBlob[12:], csSlotCodeDirectory)
	binary.BigEndian.PutUint32(superBlob[16:], 20)
	copy(superBlob[20:], codeDir)

	if rem := len(superBlob) % 16; rem != 0 {
		padding := make([]byte, 16-rem)
		superBlob = append(superBlob, padding...)
		binary.BigEndian.PutUint32(superBlob[4:], uint32(len(superBlob)))
	}

	return superBlob
}
