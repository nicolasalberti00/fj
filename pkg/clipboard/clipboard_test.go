package clipboard

import (
	"errors"
	"os/exec"
	"runtime"
	"testing"
)

func TestCopy(t *testing.T) {
	// Skip test if not running on a supported platform
	if runtime.GOOS != "darwin" && runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		t.Skip("Skipping test on unsupported platform")
	}

	// Test with a simple string
	testStr := "Test clipboard content"
	err := Copy(testStr)

	// On CI environments, clipboard commands might not be available
	// so we'll check if the error is related to command not found
	if err != nil {
		if isCommandNotFoundError(err) {
			t.Skip("Clipboard command not available, skipping test")
		} else {
			t.Errorf("Copy() error = %v", err)
			return
		}
	}

	// Ideally, we would verify the clipboard content here,
	// but that's challenging in an automated test environment
	// For now, we'll just check that the function didn't error
}

// isCommandNotFoundError checks if an error is related to a command not being found
func isCommandNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, exec.ErrNotFound) {
		return true
	}

	return false
}
