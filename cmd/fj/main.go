package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fj/pkg/clipboard"
	"fj/pkg/config"
	"fj/pkg/formatter"
)

const (
	version = "0.1.0"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		_, _ = fmt.Fprintf(os.Stderr, "Using default configuration.\n")
		cfg = config.DefaultConfig()
	}

	// Parse command line flags
	cmdConfig := parseFlags(cfg)

	// Process input
	inputData, err := getInput(cmdConfig.TrustAllURLs)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while getting input: %v\n", err)
		os.Exit(1)
	}

	// Format JSON
	opts := formatter.Options{
		IndentSpaces: cmdConfig.IndentSpaces,
		SortKeys:     cmdConfig.SortKeys,
	}

	formattedJSON, err := formatter.Format(inputData, opts)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)

		// Try auto-correction if formatting fails
		_, _ = fmt.Fprintf(os.Stderr, "Attempting to auto-correct JSON...\n")
		correctedJSON, corrErr := formatter.AutoCorrect(inputData)
		if corrErr != nil {
			fmt.Fprintf(os.Stderr, "Auto-correction failed: %v\n", corrErr)
			os.Exit(1)
		}

		// Try formatting again with corrected JSON
		formattedJSON, err = formatter.Format(correctedJSON, opts)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error formatting corrected JSON: %v\n", err)
			os.Exit(1)
		}

		_, _ = fmt.Fprintf(os.Stderr, "Auto-correction successful!\n")
	}

	// Output formatted JSON
	fmt.Println(string(formattedJSON))

	// Copy to clipboard if requested
	if cmdConfig.CopyToClipboard {
		if err := clipboard.Copy(string(formattedJSON)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Println("Copied to clipboard!")
		}
	}

	// Save to file if requested
	if cmdConfig.OutputDir != "" {
		outputPath := generateOutputPath(cmdConfig.OutputDir)
		if err := saveToFile(formattedJSON, outputPath); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to save to file: %v\n", err)
		} else {
			fmt.Printf("Saved to %s\n", outputPath)
		}
	}
}

// parseFlags parses command line flags and returns a Config
func parseFlags(defaultCfg config.Config) config.Config {
	// Define flags
	indentPtr := flag.Int("indent", defaultCfg.IndentSpaces, "Number of spaces for indentation")
	sortPtr := flag.Bool("sort", defaultCfg.SortKeys, "Sort object keys")
	silentPtr := flag.Bool("silent", defaultCfg.SilentMode, "Silent mode")
	clipboardPtr := flag.Bool("clipboard", defaultCfg.CopyToClipboard, "Copy result to clipboard")
	saveDirPtr := flag.Bool("save-to-dir", defaultCfg.SaveToDir, "Save to directory")
	outputDirPtr := flag.String("outdir", defaultCfg.OutputDir, "Output directory for saved files")
	trustPtr := flag.Bool("trust-all", defaultCfg.TrustAllURLs, "Trust all URLs without prompting")
	versionPtr := flag.Bool("version", false, "Show version information")
	helpPtr := flag.Bool("help", false, "Show help information")
	saveConfigPtr := flag.Bool("save-config", false, "Save current flags as default configuration")

	// Parse flags
	flag.Parse()

	// Show version and exit if requested
	if *versionPtr {
		fmt.Printf("fj version %s\n", version)
		os.Exit(0)
	}

	// Show help and exit if requested
	if *helpPtr {
		showHelp()
		os.Exit(0)
	}

	// Create config from flags
	cfg := config.Config{
		IndentSpaces:    *indentPtr,
		SortKeys:        *sortPtr,
		SilentMode:      *silentPtr,
		CopyToClipboard: *clipboardPtr,
		SaveToDir:       *saveDirPtr,
		OutputDir:       *outputDirPtr,
		TrustAllURLs:    *trustPtr,
		MaxMemoryMB:     defaultCfg.MaxMemoryMB,
		MaxProcessors:   defaultCfg.MaxProcessors,
		LogToFile:       defaultCfg.LogToFile,
		LogFilePath:     defaultCfg.LogFilePath,
	}

	// Save config if requested
	if *saveConfigPtr {
		if err := config.SaveConfig(cfg); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to save configuration: %v\n", err)
		} else {
			fmt.Println("Configuration saved successfully!")
		}
	}

	return cfg
}

// getInput reads JSON input from URL, stdin or file
func getInput(trustAllURLs bool) ([]byte, error) {
	args := flag.Args()

	// No args, so we check if it's from terminal or is from a pipe
	if len(args) <= 0 {
		// Check type of file from stdin
		file, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat stdin: %v", err)
		}
		if (file.Mode() & os.ModeCharDevice) != 0 {
			return io.ReadAll(os.Stdin)
		}
		return nil, errors.New("no input file specified in pipe")
	}

	// We have args, so we can treat the first one
	input := strings.TrimSpace(args[0])

	// 1. URL Handling
	inputURL, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("input is not a valid URL")
	}
	if inputURL != nil {
		// Security prompt for URLs unless trust-all is enabled
		if !trustAllURLs {
			fmt.Printf("Do you trust the URL: %s? [y/n] ", input)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				return nil, fmt.Errorf("failed to read input from URL: %v", err)
			}

			if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
				return nil, fmt.Errorf("URL access denied by user")
			}
		}

		return readFromURL(input)
	}

	// 2. We try to read a file
	inputFile, err := os.ReadFile(input)
	// If no err, we got a file
	if err == nil {
		return inputFile, nil
	}
	// 3. We have an error while reading the file, so we treat it as a raw JSON string
	if !json.Valid([]byte(input)) {
		return nil, errors.New("invalid JSON input")
	}
	return []byte(input), nil
}

// readFromURL fetches JSON from a URL
func readFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// generateOutputPath generates a file path for saving output
func generateOutputPath(outputDir string) string {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to create output directory: %v\n", err)
		outputDir = "."
	}

	// Generate filename based on current time
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("json_%s.json", timestamp)

	return filepath.Join(outputDir, filename)
}

// saveToFile saves data to a file
func saveToFile(data []byte, path string) error {
	return os.WriteFile(path, data, 0644)
}

// showHelp displays help information
func showHelp() {
	helpText := `fj - JSON formatter utility

Usage:
  fj [options] [file|url]

Options:
  -indent int       Number of spaces for indentation (default 2)
  -sort             Sort object keys
  -clipboard        Copy result to clipboard (default true)
  -outdir string    Output directory for saved files
  -trust-all        Trust all URLs without prompting
  -save-config      Save current flags as default configuration
  -version          Show version information
  -help             Show this help information

Examples:
  fj file.json                  Format JSON from file
  fj https://example.com/data   Format JSON from URL
  cat file.json | fj            Format JSON from stdin
  fj -indent 4 file.json        Format with 4-space indentation
  fj -sort file.json            Format with sorted keys

Configuration:
  fj uses a configuration file stored in:
  - Windows: %APPDATA%\fj\config.json
  - macOS:   ~/Library/Application Support/fj/config.json
  - Linux:   ~/.config/fj/config.json
`
	fmt.Print(helpText)
}
