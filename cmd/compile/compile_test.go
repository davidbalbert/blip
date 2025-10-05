package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var (
	blipCompileBinary string
	blipLinkBinary    string
	setupOnce         sync.Once
)

func setupTestEnvironment(t *testing.T) {
	setupOnce.Do(func() {
		tmpDir := t.TempDir()
		blipCompileBinary = filepath.Join(tmpDir, "blip-compile")
		blipLinkBinary = filepath.Join(tmpDir, "blip-link")

		buildCompileCmd := exec.Command("go", "build", "-o", blipCompileBinary, ".")
		if err := buildCompileCmd.Run(); err != nil {
			t.Fatalf("failed to build blip compile binary: %v", err)
		}

		buildLinkCmd := exec.Command("go", "build", "-o", blipLinkBinary, "../link")
		if err := buildLinkCmd.Run(); err != nil {
			t.Fatalf("failed to build blip link binary: %v", err)
		}
	})
}

func TestCompiler(t *testing.T) {
	setupTestEnvironment(t)

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
			err := os.WriteFile(blFile, []byte(tc.source), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			compileCmd := exec.Command(blipCompileBinary, blFile)
			if err := compileCmd.Run(); err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			objFile := blFile + ".o"
			exeFile := filepath.Join(tmpDir, "test")

			linkCmd := exec.Command(blipLinkBinary, objFile, exeFile)
			if err := linkCmd.Run(); err != nil {
				t.Fatalf("linking failed: %v", err)
			}
			runCmd := exec.Command(exeFile)
			err = runCmd.Run()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode := exitErr.ExitCode()
					if exitCode != tc.expected {
						t.Errorf("got exit code %d, want %d", exitCode, tc.expected)
					}
				} else {
					t.Fatalf("execution failed: %v", err)
				}
			} else {
				// Command succeeded (exit code 0)
				if tc.expected != 0 {
					t.Errorf("got exit code 0, want %d", tc.expected)
				}
			}
		})
	}
}
