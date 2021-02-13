package craigslist

import (
	"context"
	"fmt"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/internal/util"
	"github.com/llvtt/craig/types"
	"time"
)

type Indexer interface{}

type DynamoDBIndexer struct {
	mgr db.DataAccessManager
}

func NewDynamoDBIndexer(mgr db.DataAccessManager) *DynamoDBIndexer {
	return &DynamoDBIndexer{mgr}
}

func (idx *DynamoDBIndexer) Index(ctx context.Context, scraper CraigslistScraper) error {
	var (
		err          error
		item         *types.CraigslistItem
		newItemCount int
	)

	indexDate := time.Now()
	for item, err = scraper.Next(); err == nil; item, err = scraper.Next() {
		var previousItem types.CraigslistItem
		item.IndexDate = indexDate
		if upsertErr := idx.mgr.Table("items").Upsert(ctx, item, &previousItem); upsertErr != nil {
			return upsertErr
		}
		if previousItem.IsEmpty() {
			newItemCount++
		}
	}
	if err != util.IteratorExhausted {
		return err
	}

	fmt.Println("new item count", newItemCount)
	return nil
}
