package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPopulateStringCombinationsAtLength(t *testing.T) {

	tests := []struct {
		chars                string
		uniqueCharsCount     int
		generateLength       int
		expectedCombinations int
	}{
		{
			chars:                "*+",
			generateLength:       2,
			uniqueCharsCount:     2,
			expectedCombinations: 4,
		},
		{
			chars:                "xx",
			generateLength:       2,
			uniqueCharsCount:     1,
			expectedCombinations: 1,
		},
		{
			chars:                "*+",
			generateLength:       3,
			uniqueCharsCount:     2,
			expectedCombinations: 8,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			results := make(map[string]bool)
			PopulateStringCombinationsAtLength(results, tc.chars, "", tc.generateLength)
			assert.Equal(t, tc.expectedCombinations, len(results))
			assert.Equal(t, int(math.Pow(float64(tc.uniqueCharsCount), float64(tc.generateLength))), len(results))
		})
	}

}
