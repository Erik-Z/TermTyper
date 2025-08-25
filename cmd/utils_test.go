package cmd

import (
	"testing"
)

func TestMapToKeysSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    map[int]bool
		expected []int
	}{
		{
			name:     "empty map",
			input:    map[int]bool{},
			expected: []int{},
		},
		{
			name:     "single key",
			input:    map[int]bool{1: true},
			expected: []int{1},
		},
		{
			name:     "multiple keys",
			input:    map[int]bool{1: true, 2: false, 3: true},
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToKeysSlice(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}

			for _, expectedKey := range tt.expected {
				found := false
				for _, resultKey := range result {
					if resultKey == expectedKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected key %d not found in result", expectedKey)
				}
			}
		})
	}
}

func TestDeleteLastChar(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		expected []rune
	}{
		{
			name:     "empty slice",
			input:    []rune{},
			expected: []rune{},
		},
		{
			name:     "single character",
			input:    []rune{'a'},
			expected: []rune{},
		},
		{
			name:     "multiple characters",
			input:    []rune{'h', 'e', 'l', 'l', 'o'},
			expected: []rune{'h', 'e', 'l', 'l'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deleteLastChar(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(result))
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("at index %d: expected %c, got %c", i, expected, result[i])
				}
			}
		})
	}
}

func TestContainsChar(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		char     rune
		expected bool
	}{
		{
			name:     "empty slice",
			input:    []rune{},
			char:     'a',
			expected: false,
		},
		{
			name:     "character found",
			input:    []rune{'h', 'e', 'l', 'l', 'o'},
			char:     'e',
			expected: true,
		},
		{
			name:     "character not found",
			input:    []rune{'h', 'e', 'l', 'l', 'o'},
			char:     'x',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsChar(tt.input, tt.char)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAverageStringLen(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: 0,
		},
		{
			name:     "single string",
			input:    []string{"hello"},
			expected: 5,
		},
		{
			name:     "multiple strings",
			input:    []string{"hi", "hello", "world"},
			expected: 4, // (2 + 5 + 5) / 3 = 4
		},
		{
			name:     "strings with different lengths",
			input:    []string{"a", "bb", "ccc"},
			expected: 2, // (1 + 2 + 3) / 3 = 2
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := averageStringLen(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestAverageLineLen(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: 0,
		},
		{
			name:     "single line",
			input:    []string{"hello"},
			expected: 5,
		},
		{
			name:     "multiple lines",
			input:    []string{"hi", "hello", "world"},
			expected: 3, // (2 + 5) / 2 = 3 (last line dropped)
		},
		{
			name:     "two lines",
			input:    []string{"short", "very long line"},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := averageLineLen(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
