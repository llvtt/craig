package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/types"
)

type Craig struct {
	db  db.DataAccessManager
	ctx context.Context
}

func NewCraig() *Craig {
	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	return &Craig{db.NewDynamoDBAccessManager(dynamo), context.Background()}
}

func (craig *Craig) Run() error {
	var (
		err          error
		item         *types.CraigslistItem
		newItemCount int
	)

	scraper := craigslist.NewScraper()
	for item, err = scraper.Next(); err == nil; item, err = scraper.Next() {
		var previousItem types.CraigslistItem
		if upsertErr := craig.db.Table("items").Upsert(craig.ctx, item, &previousItem); upsertErr != nil {
			return upsertErr
		}
		if previousItem.IsEmpty() {
			newItemCount++
		}
	}
	if err != craigslist.IteratorExhausted {
		return err
	}

	fmt.Println("new item count", newItemCount)
	return nil
}

func main() {
	craig := NewCraig()
	if err := craig.Run(); err != nil {
		panic(err)
	}
}
