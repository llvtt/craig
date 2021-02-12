package craigslist

import (
	"strings"
	"testing"
)

func TestConstructURL(t *testing.T) {
	testCases := []struct {
		params   map[string]interface{}
		expected string
	}{
		{nil, craigslistUrl},
		{map[string]interface{}{"s": 1}, "?s=1"},
		{map[string]interface{}{"s": 1, "banana": `yellow peeled "healthy"`}, "?s=1&banana=yellow+peeled+%22healthy%22"},
	}

	for _, testCase := range testCases {
		result := constructURL(testCase.params)
		if !strings.HasSuffix(result, testCase.expected) {
			t.Errorf("%s does not end with %s", result, testCase.expected)
		}
	}
}
