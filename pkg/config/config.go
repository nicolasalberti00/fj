package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	IndentSpaces    int    `json:"indent_spaces"`
	SortKeys        bool   `json:"sort_keys"`
	SilentMode      bool   `json:"silent_mode"`
	CopyToClipboard bool   `json:"copy_to_clipboard"`
	SaveToDir       bool   `json:"save_to_dir"`
	OutputDir       string `json:"output_dir"`
	TrustAllURLs    bool   `json:"trust_all_urls"`
	MaxMemoryMB     int    `json:"max_memory_mb"`
	MaxProcessors   int    `json:"max_processors"`
	LogToFile       bool   `json:"log_to_file"`
	LogFilePath     string `json:"log_file_path"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return Config{
		IndentSpaces:    2,
		SortKeys:        false,
		CopyToClipboard: false,
		OutputDir:       filepath.Join(homeDir, "fj_output"),
		TrustAllURLs:    false,
		MaxMemoryMB:     0, // 0 means no limit
		MaxProcessors:   0, // 0 means use all available
		LogToFile:       false,
		LogFilePath:     filepath.Join(homeDir, ".fj", "fj.log"),
	}
}

// LoadConfig loads configuration from file
func LoadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); errors.Is(err, fs.ErrNotExist) {
		// Create default config
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return config, fmt.Errorf("failed to create default config: %v", err)
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// getConfigPathFunc is the function type for getting the config path
type getConfigPathFunc func() (string, error)

// getConfigPath returns the path to the config file
var getConfigPath getConfigPathFunc = func() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	configDir := filepath.Join(homeDir, ".config", "fj")

	return filepath.Join(configDir, "config.json"), nil
}
