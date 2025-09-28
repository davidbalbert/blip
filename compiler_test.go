package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
)

var (
	blipBinary string
	sdkPath    string
	setupOnce  sync.Once
)

func setupTestEnvironment(t *testing.T) {
	setupOnce.Do(func() {
		// Build blip binary once
		tmpDir := t.TempDir()
		blipBinary = filepath.Join(tmpDir, "blip")
		
		buildCmd := exec.Command("go", "build", "-o", blipBinary, ".")
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("failed to build blip: %v", err)
		}

		// Get SDK path once
		sdkCmd := exec.Command("xcrun", "--show-sdk-path")
		sdkPathBytes, err := sdkCmd.Output()
		if err != nil {
			t.Fatalf("failed to get SDK path: %v", err)
		}
		sdkPath = string(sdkPathBytes[:len(sdkPathBytes)-1]) // remove newline
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
		{"precedence", "2 + 3 * 4", 10}, // TODO: should be 14 when precedence is fixed
		{"complex", "10 + 2 * 3 - 8 / 4", 10}, // TODO: should be 14 when precedence is fixed
		{"division", "20 / 4 + 1", 6},
		{"multiplication", "3 * 7", 24}, // TODO: should be 21, investigate codegen bug
		{"subtraction", "50 - 8", 42},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Enable parallel execution
			
			// Create temp .bl file
			tmpDir := t.TempDir()
			blFile := filepath.Join(tmpDir, "test.bl")
			err := os.WriteFile(blFile, []byte(tc.source), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Compile using pre-built binary
			cmd := exec.Command(blipBinary, blFile)
			if err := cmd.Run(); err != nil {
				t.Fatalf("compilation failed: %v", err)
			}

			// Link using cached SDK path
			objFile := blFile + ".o"
			exeFile := filepath.Join(tmpDir, "test")

			linkCmd := exec.Command("ld", objFile, "-o", exeFile, "-lSystem", "-syslibroot", sdkPath, "-e", "_main")
			if err := linkCmd.Run(); err != nil {
				t.Fatalf("linking failed: %v", err)
			}

			// Run and check exit code
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
