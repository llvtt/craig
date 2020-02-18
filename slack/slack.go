package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type SlackMessage struct {
	Text        string        `json:"text"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Fallback  string `json:"fallback"`
	ImageUrl  string `json:"image_url"`
	Text      string `json:"text"`
}

type SlackClient struct {
	endpoint string
	imageScraper craigslist.ImageScraper
	logger log.Logger
}

func NewSlackClient(logger log.Logger) (*SlackClient, error) {
	endpoint := os.Getenv("CRAIG_SLACK_ENDPOINT")
	if len(endpoint) == 0 {
		return nil, errors.New("CRAIG_SLACK_ENDPOINT is empty!")
	}
	return &SlackClient{endpoint, craigslist.NewImageScraper(logger), logger}, nil
}

func (c *SlackClient) sendSlackMessage(message *SlackMessage) {
	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(c.endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 300 {
		level.Warn(c.logger).Log("msg", fmt.Sprintf("possible bad request, response was %s\n", string(responseBytes)))
	}
}

func messageTextForItem(item *types.CraigslistItem) string {
	return fmt.Sprintf(
		"*%s*\n%s\n%s",
		item.Title,
		item.Url,
		item.Description)
}

func messageTextForPriceDrop(priceDrop *types.PriceDrop) string {
	item := priceDrop.Item
	text := fmt.Sprintf("*Price just dropped from $%d to $%d.*  ", priceDrop.PreviousPrice / 100, priceDrop.CurrentPrice / 100)
	if priceDrop.PreviousPrice != priceDrop.MaxPrice {
		text += fmt.Sprintf("*_Original price was $%d._*", priceDrop.MaxPrice / 100)
	}
	text += "\n"
	now := time.Now()
	ageOfItemInDays := now.Sub(priceDrop.MaxPricePublishDate).Round(time.Hour * 24)
	text += fmt.Sprintf(" _Post has been up for the past %d days_\n", int(ageOfItemInDays.Hours() / 24))
	text += messageTextForItem(item)

	return text
}

func (c *SlackClient) SendString(format string, args ...interface{}) {
	c.sendSlackMessage(&SlackMessage{Text: fmt.Sprintf(format, args...)})
}

func (c *SlackClient) SendItem(item *types.CraigslistItem) error {
	var attachments []*Attachment
	urls, err := c.imageScraper.GetImageUrls(item)
	if err != nil {
		return utils.WrapError("Could not send item to craigslist", err)
	}
	for _, imageUrl := range urls {
		attachments = append(
			attachments,
			&Attachment{
				ImageUrl: imageUrl,
				Fallback: imageUrl,
			})
	}
	level.Info(c.logger).Log("msg", "sending slack message for item " + item.Title)
	c.sendSlackMessage(
		&SlackMessage{
			Text:        messageTextForItem(item),
			Attachments: attachments})
	return nil
}

func (c *SlackClient) SendPriceDrop(priceDrop *types.PriceDrop) error {
	item := priceDrop.Item
	var attachments []*Attachment
	urls, err := c.imageScraper.GetImageUrls(item)
	if err != nil {
		return utils.WrapError("Could not send item to craigslist", err)
	}
	for _, imageUrl := range urls {
		attachments = append(
			attachments,
			&Attachment{
				ImageUrl: imageUrl,
				Fallback: imageUrl,
			})
	}
	level.Info(c.logger).Log("msg", "sending slack message for item " + item.Title)
	c.sendSlackMessage(
		&SlackMessage{
			Text:        messageTextForPriceDrop(priceDrop),
			Attachments: attachments})
	return nil
}

