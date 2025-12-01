package link

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davidbalbert/blip/elf"
	"github.com/davidbalbert/blip/macho"
	"github.com/davidbalbert/blip/platform"
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

	var executable []byte
	switch platform.Target() {
	case platform.MacOS:
		textData, err := macho.ExtractTextSection(objData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading object file: %v\n", err)
			os.Exit(1)
		}
		executable = macho.GenerateExecutable(textData, filepath.Base(outFile))
	case platform.Linux:
		textData, err := elf.ExtractTextSection(objData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading object file: %v\n", err)
			os.Exit(1)
		}
		executable = elf.GenerateExecutable(textData)
	}

	err = os.WriteFile(outFile, executable, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing executable: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("linked %s -> %s\n", objFile, outFile)
}
