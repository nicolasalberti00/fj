package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

// Copy copies text to the system clipboard by using utilities that are present on each platform:
// - pbcopy for MacOS
// - clip for Windows
// - xclip for Linux
// This part could be adjusted in the config in a next release to let the user choose which program to use.
func Copy(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = bytes.NewBufferString(text)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not copy to clipboard: %w", err)
	}

	return nil
}
