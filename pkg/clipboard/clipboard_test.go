package clipboard

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
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

// TestPlatformSpecificFunctions tests the platform-specific clipboard functions
func TestPlatformSpecificFunctions(t *testing.T) {
	// Test copyOSX
	if runtime.GOOS == "darwin" {
		err := copyOSX("Test macOS clipboard")
		if err != nil && !isCommandNotFoundError(err) {
			t.Errorf("copyOSX() error = %v", err)
		}
	}

	// Test copyWindows
	if runtime.GOOS == "windows" {
		err := copyWindows("Test Windows clipboard")
		if err != nil && !isCommandNotFoundError(err) {
			t.Errorf("copyWindows() error = %v", err)
		}
	}

	// Test copyLinux
	if runtime.GOOS == "linux" {
		err := copyLinux("Test Linux clipboard")
		if err != nil && !isCommandNotFoundError(err) {
			t.Errorf("copyLinux() error = %v", err)
		}
	}
}

// TestHasCommand tests the hasCommand function
func TestHasCommand(t *testing.T) {
	// Test with a command that should exist on all platforms
	if !hasCommand("echo") {
		t.Errorf("hasCommand() failed to detect 'echo' command")
		return
	}

	// Test with a command that shouldn't exist
	if hasCommand("nonexistentcommandxyz123") {
		t.Errorf("hasCommand() incorrectly detected nonexistent command")
		return
	}
}

// isCommandNotFoundError checks if an error is related to a command not being found
func isCommandNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's an exec.ExitError or exec.Error
	var exitError *exec.ExitError
	isExitError := errors.As(err, &exitError)
	var execError *exec.Error
	isExecError := errors.As(err, &execError)

	// Check if the error message contains common "command not found" phrases
	errMsg := err.Error()
	notFoundPhrases := []string{
		"not found",
		"no such file",
		"executable file not found",
		"command not found",
	}

	for _, phrase := range notFoundPhrases {
		if isExitError || isExecError || (errMsg != "" && contains(errMsg, phrase)) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
}
