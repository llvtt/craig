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
	Url          string    `json:"url",db:"url"`
	Title        string    `json:"title",db:"title"`
	Description  string    `json:"description",db:"description"`
	ThumbnailUrl string    `json:"thumbnail_url",db:"thumbnail_url"`
	IndexDate    time.Time `json:"index_date",db:"index_date"`
	PublishDate  time.Time `json:"publish_date",db:"publish_date"`
	Price        int       `json:"price",db:"price"`
}

func (item *CraigslistItem) IsEmpty() bool {
	return item.Url == "" && item.Title == ""
}

type CraigslistPriceLog struct {
	Item         *CraigslistItem         `json:"item"`
	Prices       []*CraigslistPriceEntry `json:"prices"`
	MaxPrice     int                     `json:"max_price_cents"`
	CurrentPrice int                     `json:"current_price_cents"`
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
