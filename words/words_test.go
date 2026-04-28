package words

import (
	"strings"
	"testing"
)

func TestGenerateWithPunctuation(t *testing.T) {
	gen := NewGenerator()
	gen.Count = 50
	gen.Punctuation = true

	result := gen.Generate("Common words")

	if len(result) == 0 {
		t.Error("Generated runes should not be empty")
	}

	resultStr := string(result)

	t.Logf("Generated text: %s", resultStr)

	hasPunctuation := false
	for _, r := range result {
		if strings.ContainsRune(".,!?;()\"", r) {
			hasPunctuation = true
			break
		}
	}

	if !hasPunctuation {
		t.Error("Generated text should contain punctuation")
	}

	// Check that parentheses are balanced
	openParens := strings.Count(resultStr, "(")
	closeParens := strings.Count(resultStr, ")")
	if openParens != closeParens {
		t.Errorf("Unbalanced parentheses: %d open vs %d close", openParens, closeParens)
	}

	// Check that quotes are balanced
	quotes := strings.Count(resultStr, "\"")
	if quotes%2 != 0 {
		t.Errorf("Unbalanced quotes: %d quotes (should be even)", quotes)
	}

	// Check that parentheses and quotes have proper spacing
	// There should be a space before opening ( and "
	runes := []rune(resultStr)
	insideQuotes := false
	for i, r := range runes {
		if r == '"' {
			if insideQuotes {
				insideQuotes = false
			} else {
				// Opening quote - should have space before it
				if i > 0 && runes[i-1] != ' ' && runes[i-1] != '(' && runes[i-1] != ')' {
					t.Errorf("Missing space before opening quote at position %d: ...%s...", i, resultStr[max(0, i-10):min(len(resultStr), i+10)])
				}
				insideQuotes = true
			}
		}
		if r == '(' {
			// Opening parenthesis - should have space before it
			if i > 0 && runes[i-1] != ' ' && runes[i-1] != '"' {
				t.Errorf("Missing space before '(' at position %d: ...%s...", i, resultStr[max(0, i-10):min(len(resultStr), i+10)])
			}
		}
	}
}

func TestGenerateWithoutPunctuation(t *testing.T) {
	gen := NewGenerator()
	gen.Count = 50
	gen.Punctuation = false

	result := gen.Generate("Common words")

	// Should not contain punctuation except spaces
	for _, r := range result {
		if strings.ContainsRune(".,!?;()\"", r) {
			t.Errorf("Generated text should not contain punctuation when Punctuation is false, found: %c", r)
		}
	}
}

func TestCapitalizationAfterEndPunct(t *testing.T) {
	gen := NewGenerator()
	gen.Count = 100
	gen.Punctuation = true

	result := gen.Generate("Common words")
	resultStr := string(result)

	t.Logf("Generated text: %s", resultStr)

	// Check that words after sentence-ending punctuation are capitalized
	runes := []rune(resultStr)
	for i := 0; i < len(runes)-1; i++ {
		if runes[i] == '.' || runes[i] == '!' || runes[i] == '?' {
			// Next non-space, non-quote character should be uppercase
			for j := i + 1; j < len(runes); j++ {
				if runes[j] == '"' {
					// Skip opening quote
					continue
				}
				if runes[j] != ' ' {
					if runes[j] >= 'a' && runes[j] <= 'z' {
						t.Errorf("Word after sentence-ending punctuation should be capitalized: found %c at position %d (context: ...%s...)", runes[j], j, resultStr[max(0, i-5):min(len(resultStr), i+10)])
					}
					break
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
