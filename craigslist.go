package main

import (
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"

	_ "github.com/mattn/go-sqlite3"
)

type Listing []*CraigslistItem

type CraigslistItem struct {
	Url          string    `db:"url"`
	Title        string    `db:"title"`
	Description  string    `db:"description"`
	ThumbnailUrl string    `db:"thumbnail_url"`
	IndexDate    time.Time `db:"index_date"`
	PublishDate  time.Time `db:"publish_date"`
}

func extractThumbnail(item *gofeed.Item) string {
	enclosureList := item.Extensions["enc"]["enclosure"]
	if len(enclosureList) == 0 {
		return ""
	}
	return enclosureList[0].Attrs["resource"]
}

func CraigslistItemFromRssItem(item *gofeed.Item) *CraigslistItem {
	publishDate, err := time.Parse(time.RFC3339, item.Published)
	if err != nil {
		panic(err)
	}
	return &CraigslistItem{
		Url:          item.Link,
		Title:        html.UnescapeString(item.Title),
		Description:  html.UnescapeString(item.Description),
		ThumbnailUrl: extractThumbnail(item),
		IndexDate:    time.Now(),
		PublishDate:  publishDate,
	}
}

type Client struct {
	region   string
	category string
	options  *SearchOptions
	parser   *gofeed.Parser
	db       *sqlx.DB
}

func NewClient(region string) *Client {
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	client := &Client{region, "", &SearchOptions{}, gofeed.NewParser(), db}
	client.initTable()
	return client
}

type SearchOptions struct {
	HasPicture bool
}

type params map[string]string

func (self *Client) parameterString(p params) string {
	var paramParts []string
	for name, value := range p {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s",
			url.QueryEscape(name),
			url.QueryEscape(value)))
	}
	paramParts = append(paramParts, "format=rss")
	return fmt.Sprintf("?%s", strings.Join(paramParts, "&"))
}

func (self *Client) get(path string, p params) (feed *gofeed.Feed, err error) {
	url := fmt.Sprintf("http://%s.craigslist.org%s%s", self.region, path, self.parameterString(p))
	feed, err = self.parser.ParseURL(url)
	return
}

func (self *Client) Category(category string) *Client {
	self.category = category
	return self
}

func (self *Client) Options(options *SearchOptions) *Client {
	self.options = options
	return self
}

func (self *Client) Search(searchTerm string) (results Listing) {
	query := strings.Replace(searchTerm, " ", "+", -1)
	path := fmt.Sprintf("/search/%s", self.category)
	resultsFound := 1
	for startItem := 0; resultsFound > 0; startItem += resultsFound {
		feed, err := self.get(path, params{"query": query, "s": strconv.Itoa(startItem)})
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
