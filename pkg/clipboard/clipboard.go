package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Copy copies text to the system clipboard
func Copy(text string) error {
	switch runtime.GOOS {
	case "darwin":
		return copyOSX(text)
	case "windows":
		return copyWindows(text)
	case "linux":
		return copyLinux(text)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// copyOSX copies text to clipboard on macOS
func copyOSX(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = os.Stdin

	// Create a pipe to write to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Write text to stdin
	if _, err := stdin.Write([]byte(text)); err != nil {
		return err
	}

	// Close stdin
	if err := stdin.Close(); err != nil {
		return err
	}

	// Wait for the command to finish
	return cmd.Wait()
}

// copyWindows copies text to clipboard on Windows
func copyWindows(text string) error {
	cmd := exec.Command("clip")
	cmd.Stdin = os.Stdin

	// Create a pipe to write to stdin
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Write text to stdin
	if _, err := stdin.Write([]byte(text)); err != nil {
		return err
	}

	// Close stdin
	if err := stdin.Close(); err != nil {
		return err
	}

	// Wait for the command to finish
	return cmd.Wait()
}

// copyLinux copies text to clipboard on Linux
func copyLinux(text string) error {
	// Try xclip first
	if hasCommand("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = os.Stdin

		// Create a pipe to write to stdin
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return err
		}

		// Write text to stdin
		if _, err := stdin.Write([]byte(text)); err != nil {
			return err
		}

		// Close stdin
		if err := stdin.Close(); err != nil {
			return err
		}

		// Wait for the command to finish
		return cmd.Wait()
	}

	// Try xsel if xclip is not available
	if hasCommand("xsel") {
		cmd := exec.Command("xsel", "--clipboard", "--input")
		cmd.Stdin = os.Stdin

		// Create a pipe to write to stdin
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return err
		}

		// Write text to stdin
		if _, err := stdin.Write([]byte(text)); err != nil {
			return err
		}

		// Close stdin
		if err := stdin.Close(); err != nil {
			return err
		}

		// Wait for the command to finish
		return cmd.Wait()
	}

	return fmt.Errorf("no clipboard command found (xclip or xsel required)")
}

// hasCommand checks if a command is available
func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
