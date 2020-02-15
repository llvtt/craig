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
	DBFile   string             `json:"db_file"`
}

type CraigslistItem struct {
	Url          string    `json:"url"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ThumbnailUrl string    `json:"thumbnail_url"`
	IndexDate    time.Time `json:"index_date"`
	PublishDate  time.Time `json:"publish_date"`
}
