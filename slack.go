package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// See https://api.slack.com/docs/message-attachments
const SLACK_MAX_ATTACHMENTS = 100

type SlackMessage struct {
	Text        string        `json:"text"`
	Attachments []*Attachment `json:"attachments"`
}

type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Fallback  string `json:"fallback"`
	ImageUrl  string `json:"image_url"`
	Text      string `json:"text"`
}

func AttachmentFromCraiglistItem(item *CraigslistItem) *Attachment {
	return &Attachment{
		Title:     item.Title,
		TitleLink: item.Url,
		Fallback:  fmt.Sprintf("%s - %s", item.Title, item.Url),
		ImageUrl:  item.ThumbnailUrl,
		Text:      item.Description,
	}
}

func (self *Client) messageAttachmentsFromItems(items []*CraigslistItem) []*Attachment {
	attachments := make([]*Attachment, len(items))
	for _, item := range items {
		attachments = append(attachments, AttachmentFromCraiglistItem(item))
	}
	return attachments
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (self *Client) NotifySlack(endpoint, term string, items []*CraigslistItem) {
	messagesToSend := max(len(items)/SLACK_MAX_ATTACHMENTS, 1)
	for i := 0; i < messagesToSend; i++ {
		itemsToSend := items[i*SLACK_MAX_ATTACHMENTS : min((i+1)*SLACK_MAX_ATTACHMENTS, len(items))]
		attachments := self.messageAttachmentsFromItems(itemsToSend)
		messageText := fmt.Sprintf("New results for *%s* found on my list!", term)
		if i > 0 {
			messageText = fmt.Sprintf("More results for *%s*", term)
		}
		message := &SlackMessage{
			Text:        messageText,
			Attachments: attachments,
		}
		payload, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		resp, err := http.Post(endpoint, "application/json", bytes.NewReader(payload))
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}
}
