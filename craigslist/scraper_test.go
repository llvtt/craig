package craigslist

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

func clTestFixture() io.Reader {
	projectRoot, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fixturePath := path.Join(projectRoot, "test", "craigslist-result-page.html")
	file, err := os.Open(fixturePath)
	if err != nil {
		panic(err)
	}
	return file
}

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
