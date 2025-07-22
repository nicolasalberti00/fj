package formatter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Options defines formatting options
type Options struct {
	IndentSpaces int
	SortKeys     bool
}

// Format formats JSON data according to the provided options
func Format(data []byte, opts Options) ([]byte, error) {
	var jsonObj interface{}

	// Parse JSON
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, fmt.Errorf("invalid JSON: %v", err)
	}

	// Sort keys if requested
	if opts.SortKeys {
		jsonObj = sortJSONKeys(jsonObj)
	}

	// Create indentation string
	indent := strings.Repeat(" ", opts.IndentSpaces)

	// Marshal with indentation
	formattedJSON, err := json.MarshalIndent(jsonObj, "", indent)
	if err != nil {
		return nil, fmt.Errorf("error formatting JSON: %v", err)
	}

	return formattedJSON, nil
}

// sortJSONKeys recursively sorts keys in JSON objects
func sortJSONKeys(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		// Create a new sorted map
		sortedMap := make(map[string]interface{})

		// Get all keys
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}

		// Sort keys
		sort.Strings(keys)

		// Add sorted keys to new map
		for _, k := range keys {
			sortedMap[k] = sortJSONKeys(v[k])
		}

		return sortedMap
	case []interface{}:
		// Process each element in the array
		for i, val := range v {
			v[i] = sortJSONKeys(val)
		}
	}

	return data
}

// ValidateJSON checks if the provided data is valid JSON
func ValidateJSON(data []byte) (bool, error) {
	var js interface{}
	err := json.Unmarshal(data, &js)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AutoCorrect attempts to fix common JSON syntax errors
// This is a simple implementation and won't handle all cases
func AutoCorrect(data []byte) ([]byte, error) {
	str := string(data)

	// Try to fix missing quotes around keys
	// This is a very simplified approach
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])

			// If key doesn't start and end with quotes, add them
			if !strings.HasPrefix(key, "\"") && !strings.HasSuffix(key, "\"") {
				lines[i] = strings.Replace(line, key, "\""+key+"\"", 1)
			}
		}
	}
	str = strings.Join(lines, "\n")

	// Try to fix trailing commas
	str = strings.ReplaceAll(str, ",\n}", "\n}")
	str = strings.ReplaceAll(str, ",\n]", "\n]")

	// Validate the corrected JSON
	if _, err := ValidateJSON([]byte(str)); err != nil {
		return nil, fmt.Errorf("auto-correction failed: %v", err)
	}

	return []byte(str), nil
}
