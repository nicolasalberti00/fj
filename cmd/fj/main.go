package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
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
	clipboardPtr := flag.Bool("clipboard", defaultCfg.CopyToClipboard, "Copy result to clipboard")
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
		CopyToClipboard: *clipboardPtr,
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

// getInput reads JSON input from file, URL, or stdin
func getInput(trustAllURLs bool) ([]byte, error) {
	args := flag.Args()

	// If no arguments, read from stdin
	if len(args) == 0 {
		return io.ReadAll(os.Stdin)
	}

	input := args[0]

	// Check if input is a URL
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		// Security prompt for URLs unless trust-all is enabled
		if !trustAllURLs {
			fmt.Printf("Do you trust the URL: %s? [y/N] ", input)
			var response string
			fmt.Scanln(&response)

			if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
				return nil, fmt.Errorf("URL access denied by user")
			}
		}

		return readFromURL(input)
	}

	// Otherwise, treat as file path
	return os.ReadFile(input)
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
