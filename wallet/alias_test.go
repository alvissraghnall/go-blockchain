package wallet

import (
	"testing"
	"unicode"
  "strings"
)

func TestGenerateAlias(t *testing.T) {
	tests := []struct {
		name        string
		wantLength  int
		wantHyphens int
	}{
		{"Basic generation", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAlias()
			if err != nil {
				t.Errorf("GenerateAlias() error = %v", err)
				return
			}

			parts := strings.Split(got, "-")
			if len(parts) != tt.wantLength {
				t.Errorf("GenerateAlias() = %v, want %d parts", got, tt.wantLength)
			}

  		hyphenCount := strings.Count(got, "-")
			if hyphenCount != tt.wantHyphens {
				t.Errorf("GenerateAlias() = %v, want exactly %d hyphen", got, tt.wantHyphens)
			}

			for _, part := range parts {
				for _, r := range part {
					if !unicode.IsLower(r) && !unicode.IsLetter(r) {
						t.Errorf("GenerateAlias() contains invalid characters = %v", got)
						break
					}
				}
			}
		})
	}
}

func TestGenerateAlias_Uniqueness(t *testing.T) {
	const iterations = 100
	aliases := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		alias, err := GenerateAlias()
		if err != nil {
			t.Errorf("GenerateAlias() error = %v", err)
			continue
		}

		if aliases[alias] {
			t.Errorf("GenerateAlias() produced duplicate alias = %v", alias)
		}
		aliases[alias] = true
	}
}

func TestGenerateAlias_Entropy(t *testing.T) {
	const iterations = 1000
	wordCounts := make(map[string]int)

	for i := 0; i < iterations; i++ {
		alias, err := GenerateAlias()
		if err != nil {
			t.Errorf("GenerateAlias() error = %v", err)
			continue
		}

		words := strings.Split(alias, "-")
		for _, word := range words {
			wordCounts[word]++
		}
	}

	if len(wordCounts) < iterations/10 {
		t.Errorf("GenerateAlias() seems to be repeating words too often, got %d unique words out of %d iterations", len(wordCounts), iterations)
	}
}

func TestGenerateAlias_ErrorHandling(t *testing.T) {
	t.Skip("Error injection test requires mocking bip39 package")
}
