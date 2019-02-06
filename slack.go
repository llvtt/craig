package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
}

func (self *SlackClient) sendSlackMessage(message *SlackMessage) {
	payload, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(self.endpoint, "application/json", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("possible bad request, response was %s\n", string(responseBytes))
	}
}

func messageTextForItem(item *CraigslistItem) string {
	return fmt.Sprintf(
		"*%s*\n%s\n%s",
		item.Title,
		item.Url,
		item.Description)
}

func (self *SlackClient) SendString(format string, args ...interface{}) {
	self.sendSlackMessage(&SlackMessage{Text: fmt.Sprintf(format, args...)})
}

func (self *SlackClient) SendItem(endpoint string, item *CraigslistItem) {
	var attachments []*Attachment
	for _, imageUrl := range item.GetImageUrls() {
		attachments = append(
			attachments,
			&Attachment{
				ImageUrl: imageUrl,
				Fallback: imageUrl,
			})
	}
	fmt.Println("sending slack message for item " + item.Title)
	self.sendSlackMessage(
		&SlackMessage{
			Text:        messageTextForItem(item),
			Attachments: attachments})
}
