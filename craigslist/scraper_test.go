package craigslist

import (
	"github.com/llvtt/craig/types"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

func clTestFixture() io.Reader {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	projectRoot := strings.TrimSpace(string(output))
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

func TestParseItems(t *testing.T) {
	assertions := assert.New(t)
	scraper := NewScraper()

	results, resultCount, err := scraper.parseItems(clTestFixture())
	assertions.NoError(err)

	assertions.Equal(14, resultCount, results)

	included := &types.CraigslistItem{
		Url:         "https://sfbay.craigslist.org/sfc/bik/d/san-francisco-stolen-bionx-jamis-dakar/7274129390.html",
		Title:       "STOLEN: BionX Jamis Dakar Dragon 29er electric mountain bike",
		Description: "",
		IndexDate:   time.Time{},
		PublishDate: time.Date(2021, 2, 7, 13, 21, 0, 0, time.UTC),
		Price:       100,
	}
	assertions.Equal(included, results[0])
}
