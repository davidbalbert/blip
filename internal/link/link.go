package link

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidbalbert/blip/macho"
)

const (
	MH_MAGIC_64   = 0xfeedfacf
	LC_SEGMENT_64 = 0x19
)

func Main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <object-file> <output-executable>\n", os.Args[0])
		os.Exit(1)
	}

	objFile := os.Args[1]
	outFile := os.Args[2]

	objData, err := os.ReadFile(objFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading object file: %v\n", err)
		os.Exit(1)
	}

	executable, err := Link(objFile, outFile, objData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outFile, executable, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing executable: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("linked %s -> %s\n", objFile, outFile)
}

// Link links an object file into an executable.
// objFile is the path to the object file.
// outFile is the path for the output executable.
// objData is the content of the object file.
// Returns the executable data or an error.
func Link(objFile, outFile string, objData []byte) ([]byte, error) {
	textData, _, err := extractTextSection(objData)
	if err != nil {
		return nil, fmt.Errorf("error reading object file: %w", err)
	}

	executable := macho.GenerateExecutable(textData, filepath.Base(outFile))
	return executable, nil
}

func extractTextSection(data []byte) ([]byte, uint64, error) {
	if len(data) < 32 {
		return nil, 0, fmt.Errorf("file too short")
	}

	magic := binary.LittleEndian.Uint32(data[0:4])
	if magic != MH_MAGIC_64 {
		return nil, 0, fmt.Errorf("not a Mach-O 64-bit file")
	}

	ncmds := binary.LittleEndian.Uint32(data[16:20])
	offset := uint32(32) // start after header

	for i := uint32(0); i < ncmds; i++ {
		if offset+8 > uint32(len(data)) {
			return nil, 0, fmt.Errorf("truncated load command")
		}

		cmd := binary.LittleEndian.Uint32(data[offset : offset+4])
		cmdsize := binary.LittleEndian.Uint32(data[offset+4 : offset+8])

		if cmd == LC_SEGMENT_64 {
			if offset+72 > uint32(len(data)) {
				return nil, 0, fmt.Errorf("truncated segment command")
			}

			segname := string(data[offset+8 : offset+24])
			if segname[:6] == "__TEXT" {
				nsects := binary.LittleEndian.Uint32(data[offset+64 : offset+68])
				if nsects > 0 {
					sectOffset := offset + 72 // After segment command header
					sectSize := binary.LittleEndian.Uint64(data[sectOffset+40 : sectOffset+48])
					sectFileoff := binary.LittleEndian.Uint32(data[sectOffset+48 : sectOffset+52])

					if uint64(sectFileoff)+sectSize > uint64(len(data)) {
						return nil, 0, fmt.Errorf("text section extends beyond file")
					}

					textData := data[sectFileoff : uint64(sectFileoff)+sectSize]
					return textData, 0x100000000, nil // Base address for executables
				}
			}
		}

		offset += cmdsize
	}

	return nil, 0, fmt.Errorf("no __TEXT segment found")
}
