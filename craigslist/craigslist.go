package craigslist

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/llvtt/craig/types"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	_ "github.com/mattn/go-sqlite3"
)

type Listing []*types.CraigslistItem

func prependSlash(urlPart string) string {
	if urlPart == "" {
		return urlPart
	}
	return "/" + urlPart
}

func extractThumbnail(item *gofeed.Item) string {
	enclosureList := item.Extensions["enc"]["enclosure"]
	if len(enclosureList) == 0 {
		return ""
	}
	return enclosureList[0].Attrs["resource"]
}

func CraigslistItemFromRssItem(item *gofeed.Item) *types.CraigslistItem {
	publishDate, err := time.Parse(time.RFC3339, item.Published)
	if err != nil {
		panic(err)
	}
	return &types.CraigslistItem{
		Url:          item.Link,
		Title:        html.UnescapeString(item.Title),
		Description:  html.UnescapeString(item.Description),
		ThumbnailUrl: extractThumbnail(item),
		IndexDate:    time.Now(),
		PublishDate:  publishDate,
	}
}

type CraigslistClient struct {
	region   string
	category string
	options  *SearchOptions
	parser   *gofeed.Parser
	byUrl    map[string]*types.CraigslistItem
	byTitle  map[string]*types.CraigslistItem
}

func NewCraigslistClient(region string, logger log.Logger) *CraigslistClient {
	client := &CraigslistClient{region, "", &SearchOptions{}, gofeed.NewParser(),
		make(map[string]*types.CraigslistItem), make(map[string]*types.CraigslistItem)}
	return client
}

type SearchOptions struct {
	HasPicture    bool
	SubRegion     string
	Neighborhoods []int
}

type param []string
type params []param

func (self *CraigslistClient) parameterString(p params) string {
	var paramParts []string
	for _, param := range p {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", param[0], param[1]))
	}
	paramParts = append(paramParts, "format=rss")
	return fmt.Sprintf("?%s", strings.Join(paramParts, "&"))
}

func (self *CraigslistClient) optionsToParams(p params) params {
	for _, nh := range self.options.Neighborhoods {
		p = append(p, param{"nh", strconv.Itoa(nh)})
	}
	return p
}

func (self *CraigslistClient) buildUrl(path string, p params) string {
	return fmt.Sprintf(
		"http://%s.craigslist.org%s%s%s%s",
		self.region,
		path,
		prependSlash(self.options.SubRegion),
		prependSlash(self.category),
		prependSlash(self.parameterString(self.optionsToParams(p))),
	)
}

func (self *CraigslistClient) get(path string, p params) (feed *gofeed.Feed, err error) {
	url := fmt.Sprintf(self.buildUrl(path, p))
	fmt.Println(url)
	feed, err = self.parser.ParseURL(url)
	return
}

func (self *CraigslistClient) Category(category string) *CraigslistClient {
	self.category = category
	return self
}

func (self *CraigslistClient) Options(options *SearchOptions) *CraigslistClient {
	self.options = options
	return self
}

func (self *CraigslistClient) Search(searchTerm string) (results Listing) {
	query := strings.Replace(searchTerm, " ", "+", -1)
	resultsFound := 1
	for startItem := 0; resultsFound > 0; startItem += resultsFound {
		feed, err := self.get("/search", params{param{"query", query}, param{"s", strconv.Itoa(startItem)}})
		if err != nil {
			panic(err)
		}
		for _, item := range feed.Items {
			results = append(results, CraigslistItemFromRssItem(item))
			startItem += 1
		}
		resultsFound = len(feed.Items)
	}
	return
}
