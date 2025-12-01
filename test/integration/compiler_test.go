package integration_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/davidbalbert/blip/internal/cmd/compile"
	"github.com/davidbalbert/blip/internal/cmd/link"
	"github.com/davidbalbert/blip/platform"
)

var testBinary string
var repoRoot string

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

	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot = filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))

	os.Exit(m.Run())
}

func run(exePath string, targetOS platform.OS) (exitCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hostOS := platform.Host()
	var cmd *exec.Cmd
	if targetOS == hostOS {
		cmd = exec.CommandContext(ctx, exePath)
	} else if targetOS == platform.Linux && hostOS == platform.MacOS {
		cmd = exec.CommandContext(ctx, "./run.sh", exePath)
		cmd.Dir = filepath.Join(repoRoot, "test", "crosstest")
	} else {
		return 0, fmt.Errorf("unsupported target OS %q on host %q", targetOS, hostOS)
	}

	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return 0, fmt.Errorf("timed out")
	}
	if err == nil {
		return 0, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), nil
	}
	return 0, fmt.Errorf("unexpected error running command: %v\n%s", err, output)
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
		for _, targetOS := range platform.OSs {
			t.Run(tc.name+"/"+string(targetOS), func(t *testing.T) {
				tmpDir := t.TempDir()
				blFile := filepath.Join(tmpDir, "test.bl")

				if err := os.WriteFile(blFile, []byte(tc.source), 0644); err != nil {
					t.Fatalf("failed to write source file: %v", err)
				}

				compileCmd := exec.Command(testBinary, blFile)
				compileCmd.Env = append(os.Environ(),
					"BLIP_TEST_MODE=compile",
					"BLOS="+string(targetOS),
				)
				if output, err := compileCmd.CombinedOutput(); err != nil {
					t.Fatalf("compilation failed: %v\n%s", err, output)
				}

				objFile := blFile + ".o"
				exeFile := filepath.Join(tmpDir, "test")

				linkCmd := exec.Command(testBinary, objFile, exeFile)
				linkCmd.Env = append(os.Environ(),
					"BLIP_TEST_MODE=link",
					"BLOS="+string(targetOS),
				)
				if output, err := linkCmd.CombinedOutput(); err != nil {
					t.Fatalf("linking failed: %v\n%s", err, output)
				}

				exitCode, err := run(exeFile, targetOS)
				if err != nil {
					t.Fatalf("failed to run binary: %v", err)
				}
				if exitCode != tc.expected {
					t.Errorf("got %d, want %d", exitCode, tc.expected)
				}
			})
		}
	}
}
