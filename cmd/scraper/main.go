package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/slack"
	"log"
)

func main() {
	ctx := context.Background()
	sess := session.Must(session.NewSession())
	manager := db.NewDynamoDBAccessManager(dynamodb.New(sess))
	indexer := craigslist.NewDynamoDBIndexer(manager.Table("items"))
	scraper := craigslist.NewScraper()
	slacker := slack.NewSlacker()
	newItemCount := 0
	for item, err := scraper.Next(); err == nil; item, err = scraper.Next() {
		if newItem, err := indexer.Index(ctx, item); err != nil {
			log.Println("error", err.Error())
		} else if newItem {
			log.Println("indexed new item", item)
			newItemCount++
		}
	}
	fmt.Print("new items indexed:", newItemCount)
}
