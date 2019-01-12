package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/feeds"
)

type Results []*feeds.RssItem

type Client struct {
	region   string
	category string
	options  *SearchOptions
}

func NewClient(region string) *Client {
	return &Client{region, "", &SearchOptions{}}
}

type SearchOptions struct {
	HasPicture bool
}

type params map[string]string

func parameterString(p params) string {
	var paramParts []string
	for name, value := range p {
		paramParts = append(paramParts, fmt.Sprintf("%s=%s",
			url.QueryEscape(name),
			url.QueryEscape(value)))
	}
	paramParts = append(paramParts, "format=rss")
	return fmt.Sprintf("?%s", strings.Join(paramParts, "&"))
}

func (self *Client) get(path string, p params) (feed feeds.RssFeedXml, err error) {
	var (
		resp *http.Response
		body []byte
	)
	url := fmt.Sprintf("http://%s.craigslist.org%s%s", self.region, path, parameterString(p))
	fmt.Println(url)
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	fmt.Printf("%s", body)
	err = xml.Unmarshal(body, &feed)
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

func (self *Client) Search(terms ...string) (Results, error) {
	path := fmt.Sprintf("/search/%s", self.category)
	feed, err := self.get(path, params{"query": strings.Join(terms, "+")})
	return feed.Channel.Items, err
}

func main() {
	client := NewClient("sfbay")
	results, err := client.Category("ata").Options(&SearchOptions{HasPicture: true}).Search("queen", "bed")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Sprintf("found result: %v\n", result.Title)
	}
}
