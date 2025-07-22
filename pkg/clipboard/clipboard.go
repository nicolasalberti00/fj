package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Copy copies text to the system clipboard by using utilities that are present on each platform:
// - pbcopy for MacOS
// - clip for Windows
// - xclip for Linux
// This part could be adjusted in the config in a next release to let the user choose which program to use.
func Copy(text string) error {

	var copyProgram string

	switch runtime.GOOS {
	case "darwin":
		copyProgram = "pbcopy"
	case "windows":
		copyProgram = "clip"
	case "linux":
		copyProgram = "xclip"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd := exec.Command(copyProgram, text)
	cmd.Stdin = strings.NewReader(text)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not copy to clipboard: %w", err)
	}

	return nil
}
