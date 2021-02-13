package main

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/db"
)

func main() {
	sess := session.Must(session.NewSession())
	manager := db.NewDynamoDBAccessManager(dynamodb.New(sess))
	indexer := craigslist.NewDynamoDBIndexer(manager)
	if err := indexer.Index(context.Background(), craigslist.NewScraper()); err != nil {
		panic(err)
	}
}
