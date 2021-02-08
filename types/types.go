package types

import (
	"time"
)

type CraigslistSearch struct {
	Category      string   `json:"category"`
	Terms         []string `json:"terms"`
	Neighborhoods []int    `json:"nh"`
}

type CraigConfig struct {
	Region   string             `json:"region"`
	Searches []CraigslistSearch `json:"searches"`
	DBType   string             `json:"db_type"`
	DBDir    string             `json:"db_dir"`
}

type CraigslistItem struct {
	Url          string    `json:"url" dynamodbav:"url"`
	Title        string    `json:"title" dynamodbav:"title"`
	Description  string    `json:"description" dynamodbav:"description"`
	ThumbnailUrl string    `json:"thumbnail_url" dynamodbav:"thumbnail_url"`
	IndexDate    time.Time `json:"index_date" dynamodbav:"index_date"`
	PublishDate  time.Time `json:"publish_date" dynamodbav:"publish_date"`
	Price        int       `json:"price" dynamodbav:"price"`
}

type CraigslistPriceLogGet struct {
	ItemUrl      string                  `json:"item_url" dynamodbav:"item_url"`
}

type CraigslistPriceLog struct {
	ItemUrl      string                  `json:"item_url" dynamodbav:"item_url"`
	Prices       []*CraigslistPriceEntry `json:"prices" dynamodbav:"prices"`
	MaxPrice     int                     `json:"max_price_cents" dynamodbav:"max_price_cents"`
	CurrentPrice int                     `json:"current_price_cents" dynamodbav:"current_price_cents"`
}

type CraigslistPriceEntry struct {
	Price       int       `json:"price"`
	PublishDate time.Time `json:"publish_date"`
}

type PriceDrop struct {
	Item                     *CraigslistItem
	CurrentPrice             int
	PreviousPrice            int
	MaxPrice                 int
	PreviousPricePublishDate time.Time
	MaxPricePublishDate      time.Time
}

