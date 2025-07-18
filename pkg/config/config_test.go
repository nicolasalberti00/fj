package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Check default values
	if cfg.IndentSpaces != 2 {
		t.Errorf("DefaultConfig().IndentSpaces = %v, want %v", cfg.IndentSpaces, 2)
	}

	if cfg.SortKeys != false {
		t.Errorf("DefaultConfig().SortKeys = %v, want %v", cfg.SortKeys, false)
	}

	if cfg.CopyToClipboard != true {
		t.Errorf("DefaultConfig().CopyToClipboard = %v, want %v", cfg.CopyToClipboard, true)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fj-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to remove temp directory: %v", err)
		}
	}(tempDir)

	// Create a test config
	testCfg := Config{
		IndentSpaces:    4,
		SortKeys:        true,
		CopyToClipboard: false,
		OutputDir:       "/test/output",
		TrustAllURLs:    true,
		MaxMemoryMB:     1024,
		MaxProcessors:   2,
		LogToFile:       true,
		LogFilePath:     "/test/log.txt",
	}

	// Override getConfigPath for testing
	originalGetConfigPath := getConfigPath
	defer func() { getConfigPath = originalGetConfigPath }()

	getConfigPath = func() (string, error) {
		return filepath.Join(tempDir, "config.json"), nil
	}

	// Save the config
	if err := SaveConfig(testCfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Load the config
	loadedCfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Compare the configs
	if loadedCfg.IndentSpaces != testCfg.IndentSpaces {
		t.Errorf("LoadConfig().IndentSpaces = %v, want %v", loadedCfg.IndentSpaces, testCfg.IndentSpaces)
	}

	if loadedCfg.SortKeys != testCfg.SortKeys {
		t.Errorf("LoadConfig().SortKeys = %v, want %v", loadedCfg.SortKeys, testCfg.SortKeys)
	}

	if loadedCfg.CopyToClipboard != testCfg.CopyToClipboard {
		t.Errorf("LoadConfig().CopyToClipboard = %v, want %v", loadedCfg.CopyToClipboard, testCfg.CopyToClipboard)
	}

	if loadedCfg.OutputDir != testCfg.OutputDir {
		t.Errorf("LoadConfig().OutputDir = %v, want %v", loadedCfg.OutputDir, testCfg.OutputDir)
	}

	if loadedCfg.TrustAllURLs != testCfg.TrustAllURLs {
		t.Errorf("LoadConfig().TrustAllURLs = %v, want %v", loadedCfg.TrustAllURLs, testCfg.TrustAllURLs)
	}

	if loadedCfg.MaxMemoryMB != testCfg.MaxMemoryMB {
		t.Errorf("LoadConfig().MaxMemoryMB = %v, want %v", loadedCfg.MaxMemoryMB, testCfg.MaxMemoryMB)
	}

	if loadedCfg.MaxProcessors != testCfg.MaxProcessors {
		t.Errorf("LoadConfig().MaxProcessors = %v, want %v", loadedCfg.MaxProcessors, testCfg.MaxProcessors)
	}

	if loadedCfg.LogToFile != testCfg.LogToFile {
		t.Errorf("LoadConfig().LogToFile = %v, want %v", loadedCfg.LogToFile, testCfg.LogToFile)
	}

	if loadedCfg.LogFilePath != testCfg.LogFilePath {
		t.Errorf("LoadConfig().LogFilePath = %v, want %v", loadedCfg.LogFilePath, testCfg.LogFilePath)
	}
}
