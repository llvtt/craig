package craigslist

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"html"
	"net/http"
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

func (c *CraigslistClient) CraigslistItemFromRssItem(item *gofeed.Item) (*types.CraigslistItem, error) {
	publishDate, err := time.Parse(time.RFC3339, item.Published)
	if err != nil {
		panic(err)
	}

	price, err := c.getPrice(item)
	if err != nil {
		return nil, utils.WrapError("Could no get price for item: "+item.Link, err)
	}
	level.Info(c.logger).Log("price", price)

	return &types.CraigslistItem{
		Url:          item.Link,
		Title:        html.UnescapeString(item.Title),
		Description:  html.UnescapeString(item.Description),
		ThumbnailUrl: extractThumbnail(item),
		IndexDate:    time.Now(),
		PublishDate:  publishDate,
		Price:        price,
	}, nil
}

func (c *CraigslistClient) getPrice(item *gofeed.Item) (float32, error) {
	url := item.Link
	res, err := http.Get(url)
	if err != nil {
		return 0, utils.WrapError("Could not load page: "+url, err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return 0, errors.New(fmt.Sprintf("Could not load page. Status code was: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 0, utils.WrapError("Could not parse page body at url: "+url, err)
	}

	selection := doc.Find(".price")
	priceText := selection.Text()
	if strings.HasPrefix(priceText, "$") {
		priceText = priceText[1:]
	}

	var price float32
	price64, err := strconv.ParseFloat(priceText, 32)
	if err != nil {
		return 0, utils.WrapError(fmt.Sprintf("Could not parse price from text %s. Url was %s", priceText, url), err)
	}
	price = float32(price64)


	level.Info(c.logger).Log("price", price)
	return price, nil
}

type CraigslistClient struct {
	region   string
	category string
	options  *SearchOptions
	parser   *gofeed.Parser
	byUrl    map[string]*types.CraigslistItem
	byTitle  map[string]*types.CraigslistItem
	logger   log.Logger
}

func NewCraigslistClient(region string, logger log.Logger) *CraigslistClient {
	client := &CraigslistClient{region, "", &SearchOptions{}, gofeed.NewParser(),
		make(map[string]*types.CraigslistItem), make(map[string]*types.CraigslistItem), logger}
	return client
}

type SearchOptions struct {
	HasPicture    bool
	SubRegion     string
	Neighborhoods []int
}

type param []string
type params []param

func (c *CraigslistClient) parameterString(p params) string {
	var paramParts []string
	for _, param := range p {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s", param[0], param[1]))
	}
	paramParts = append(paramParts, "format=rss")
	return fmt.Sprintf("?%s", strings.Join(paramParts, "&"))
}

func (c *CraigslistClient) optionsToParams(p params) params {
	for _, nh := range c.options.Neighborhoods {
		p = append(p, param{"nh", strconv.Itoa(nh)})
	}
	return p
}

func (c *CraigslistClient) buildUrl(path string, p params) string {
	return fmt.Sprintf(
		"http://%s.craigslist.org%s%s%s%s",
		c.region,
		path,
		prependSlash(c.options.SubRegion),
		prependSlash(c.category),
		prependSlash(c.parameterString(c.optionsToParams(p))),
	)
}

func (c *CraigslistClient) get(path string, p params) (feed *gofeed.Feed, err error) {
	url := fmt.Sprintf(c.buildUrl(path, p))
	level.Info(c.logger).Log("msg", "Getting url: "+url)
	feed, err = c.parser.ParseURL(url)
	return
}

func (c *CraigslistClient) Category(category string) *CraigslistClient {
	c.category = category
	return c
}

func (c *CraigslistClient) Options(options *SearchOptions) *CraigslistClient {
	c.options = options
	return c
}

func (c *CraigslistClient) Search(searchTerm string) (Listing, error) {
	query := strings.Replace(searchTerm, " ", "+", -1)
	resultsFound := 1
	var results Listing
	for startItem := 0; resultsFound > 0; startItem += resultsFound {
		feed, err := c.get("/search", params{param{"query", query}, param{"s", strconv.Itoa(startItem)}})
		if err != nil {
			return nil, utils.WrapError("Could not execute craigslist search request. ", err)
		}
		for _, item := range feed.Items {
			rssItem, err := c.CraigslistItemFromRssItem(item)
			if err != nil {
				// skip the item, don't fail the whole request
				level.Error(c.logger).Log(fmt.Sprintf("Could not convert rss item into craigslist item. Item was %s", item), err)
			}
			results = append(results, rssItem)
			startItem += 1
		}
		resultsFound = len(feed.Items)
	}
	return results, nil
}
