package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/davidbalbert/blip/internal/cmd/compile"
	"github.com/davidbalbert/blip/internal/cmd/link"
)

var testBinary string

func TestMain(m *testing.M) {
	switch os.Getenv("BLIP_TEST_MODE") {
	case "compile":
		compile.Main()
		os.Exit(0)
	case "link":
		link.Main()
		os.Exit(0)
	}

	var err error
	testBinary, err = os.Executable()
	if err != nil {
		panic(fmt.Sprintf("failed to get test binary path: %v", err))
	}

	os.Exit(m.Run())
}

func run(cmd *exec.Cmd) int {
	err := cmd.Run()
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	panic(fmt.Sprintf("unexpected error running command: %v", err))
}

func TestCompiler(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected int
	}{
		{"simple_number", "42", 42},
		{"precedence", "2 + 3 * 4", 14},
		{"complex", "10 + 2 * 3 - 8 / 4", 14},
		{"division", "20 / 4 + 1", 6},
		{"multiplication", "3 * 7", 21},
		{"subtraction", "50 - 8", 42},
		{"parens_basic", "(2 + 3) * 4", 20},
		{"parens_change_precedence", "2 * (3 + 4)", 14},
		{"nested_parens", "((2 + 3) * 4) - 5", 15},
		{"parens_division", "20 / (2 + 2)", 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			blFile := filepath.Join(tmpDir, "test.bl")

			if err := os.WriteFile(blFile, []byte(tc.source), 0644); err != nil {
				t.Fatalf("failed to write source file: %v", err)
			}

			compileCmd := exec.Command(testBinary, blFile)
			compileCmd.Env = append(os.Environ(), "BLIP_TEST_MODE=compile")
			if err := compileCmd.Run(); err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			objFile := blFile + ".o"
			exeFile := filepath.Join(tmpDir, "test")

			linkCmd := exec.Command(testBinary, objFile, exeFile)
			linkCmd.Env = append(os.Environ(), "BLIP_TEST_MODE=link")
			if err := linkCmd.Run(); err != nil {
				t.Fatalf("linking failed: %v", err)
			}

			runCmd := exec.Command(exeFile)
			exitCode := run(runCmd)
			if exitCode != tc.expected {
				t.Errorf("got %d, want %d", exitCode, tc.expected)
			}
		})
	}
}
