package formatter

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		opts    Options
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			input:   `{"name":"John","age":30,"city":"New York"}`,
			opts:    Options{IndentSpaces: 2, SortKeys: false},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"name":"John","age":30,"city":"New York"`,
			opts:    Options{IndentSpaces: 2, SortKeys: false},
			wantErr: true,
		},
		{
			name:    "Empty JSON object",
			input:   `{}`,
			opts:    Options{IndentSpaces: 2, SortKeys: false},
			wantErr: false,
		},
		{
			name:    "Empty JSON array",
			input:   `[]`,
			opts:    Options{IndentSpaces: 2, SortKeys: false},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format([]byte(tt.input), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the output is valid JSON
				var js interface{}
				if err := json.Unmarshal(got, &js); err != nil {
					t.Errorf("Format() produced invalid JSON: %v", err)
				}
			}
		})
	}
}

func TestSortJSONKeys(t *testing.T) {
	input := map[string]interface{}{
		"c": 3,
		"a": 1,
		"b": 2,
	}

	expected := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	result := sortJSONKeys(input)

	// Convert to JSON for comparison
	resultJSON, _ := json.Marshal(result)
	expectedJSON, _ := json.Marshal(expected)

	if !reflect.DeepEqual(resultJSON, expectedJSON) {
		t.Errorf("sortJSONKeys() = %v, want %v", string(resultJSON), string(expectedJSON))
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			input:   `{"name":"John","age":30}`,
			want:    true,
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"name":"John","age":30`,
			want:    false,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   ``,
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAutoCorrect(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Already valid JSON",
			input:   `{"name":"John","age":30}`,
			wantErr: false,
		},
		{
			name:    "Trailing comma",
			input:   `{"name":"John","age":30,}`,
			wantErr: false,
		},
		{
			name:    "Missing quotes around key",
			input:   `{name:"John","age":30}`,
			wantErr: false,
		},
		{
			name:    "Severely malformed JSON",
			input:   `{name:"John","age:30`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AutoCorrect([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("AutoCorrect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the output is valid JSON
				var js interface{}
				if err := json.Unmarshal(got, &js); err != nil {
					t.Errorf("AutoCorrect() produced invalid JSON: %v", err)
				}
			}
		})
	}
}
