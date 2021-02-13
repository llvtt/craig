package craigslist

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/llvtt/craig/types"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	craigslistUrl    = "https://sfbay.craigslist.org/search/sfc/bia"
	timeFormat       = "2006-01-02 15:04"
	throttleDuration = 10 * time.Second
)

type CraigslistScraper interface {
	Next() (*types.CraigslistItem, error)
}

type iterationError string

func (ie iterationError) Error() string {
	return string(ie)
}

// IteratorExhausted is returned when an Iterator has reached the end of its documents.
const IteratorExhausted iterationError = "iterator exhausted"

// HTML-based scraper for Craigslist
type HTMLScraper struct {
	// Starting index of results on the current craigslist page
	pageResultsStartIndex int
	// Ticker for use in self-throttling
	ticker *time.Ticker
	// Current position in currentResults
	currentResultIndex int
	// The index of the next item to request
	nextItemIndex int
	// Slice of craigslist items for the current page
	currentResults []*types.CraigslistItem
}

func NewScraper() *HTMLScraper {
	return new(HTMLScraper)
}

func (scraper *HTMLScraper) Next() (item *types.CraigslistItem, err error) {
	if scraper.currentResultIndex >= len(scraper.currentResults) {
		// fetch more results
		if err = scraper.getNextPage(); err != nil {
			return
		}
		scraper.currentResultIndex = 0
	}

	if scraper.currentResultIndex < len(scraper.currentResults) {
		item = scraper.currentResults[scraper.currentResultIndex]
		scraper.currentResultIndex++
		log.Println("currentResultIndex", scraper.currentResultIndex)
	}

	return
}

func constructURL(params map[string]interface{}) string {
	var queryString strings.Builder
	queryString.WriteString(craigslistUrl)

	if len(params) > 0 {
		queryString.WriteString("?")
		index := 0
		for name, value := range params {
			queryString.WriteString(url.QueryEscape(name))
			queryString.WriteString("=")
			queryString.WriteString(url.QueryEscape(fmt.Sprint(value)))
			if index+1 < len(params) {
				queryString.WriteString("&")
				index++
			}
		}
	}

	return queryString.String()
}

func (scraper *HTMLScraper) throttle() {
	if scraper.ticker == nil {
		// First iteration does not block
		scraper.ticker = time.NewTicker(throttleDuration)
	} else {
		log.Println("throttling")
		t := <-scraper.ticker.C
		log.Println("awakened from throttle at", t)
	}
}

func parseItem(s *goquery.Selection) (item *types.CraigslistItem, err error) {
	if s.Length() != 1 {
		return nil, fmt.Errorf("WARN - result row has %d children (only 1 expected)", s.Length())
	}
	item = new(types.CraigslistItem)

	item.Url = s.Find("a").AttrOr("href", "")
	resultInfo := s.Find(".result-info")

	if timeString, ok := resultInfo.Find(".result-date").Attr("datetime"); ok {
		item.PublishDate, err = time.Parse(timeFormat, timeString)
		if err != nil {
			return nil, err
		}
	}

	item.Title = resultInfo.Find(".result-heading .result-title").Text()
	priceString := strings.TrimPrefix(resultInfo.Find(".result-meta .result-price").Text(), "$")
	item.Price, err = strconv.Atoi(priceString)
	// Convert to cents
	item.Price *= 100

	return item, nil
}

func (scraper *HTMLScraper) parseItems(reader io.Reader) (results []*types.CraigslistItem, resultsCount int, err error) {
	var doc *goquery.Document
	var decoded io.Reader
	if decoded, err = decodeHTMLBody(reader); err != nil {
		return
	} else {
		doc, err = goquery.NewDocumentFromReader(decoded)
		if err != nil {
			return
		}
	}

	firstResult := doc.Find(".result-row").First()
	resultRows := firstResult.AddSelection(firstResult.NextUntil(".nearby"))
	resultsCount = len(resultRows.Nodes)
	resultRows.Each(func(_ int, selectedNode *goquery.Selection) {
		item, err := parseItem(selectedNode.First())
		if err != nil {
			log.Println("error", err.Error())
		}
		if item != nil {
			results = append(results, item)
		}
	})

	if resultsCount == 0 {
		err = IteratorExhausted
	}

	return
}

func (scraper *HTMLScraper) getNextPage() error {
	scraper.throttle()

	log.Println("fetching page at index", scraper.nextItemIndex)

	requestURL := constructURL(map[string]interface{}{
		"s":     scraper.nextItemIndex,
		"query": "mountain bike",
	})
	request, err := http.NewRequest(http.MethodGet, requestURL, http.NoBody)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unsuccessful request response: %d %s", res.StatusCode, string(body))
	}

	var resultCount int
	scraper.currentResults, resultCount, err = scraper.parseItems(res.Body)
	scraper.nextItemIndex += resultCount

	return err
}

// https://github.com/PuerkitoBio/goquery/wiki/Tips-and-tricks
func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// DecodeHTMLBody returns an decoding reader of the html Body for the specified `charset`
// If `charset` is empty, DecodeHTMLBody tries to guess the encoding from the content
// https://github.com/PuerkitoBio/goquery/wiki/Tips-and-tricks
func decodeHTMLBody(body io.Reader) (io.Reader, error) {
	contentCharset := detectContentCharset(body)
	e, err := htmlindex.Get(contentCharset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	return body, nil
}
