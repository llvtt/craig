package craigslist

import (
	"context"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/types"
	"time"
)

type Indexer interface {
	// Index a CraigslistItem
	// Return true if the item was new, and any error
	Index(item *types.CraigslistItem) (bool, error)
}

type DynamoDBIndexer struct {
	access db.DataAccess
}

func NewDynamoDBIndexer(acc db.DataAccess) *DynamoDBIndexer {
	return &DynamoDBIndexer{acc}
}

func (idx *DynamoDBIndexer) Index(ctx context.Context, item *types.CraigslistItem) (newRecord bool, err error) {
	item.IndexDate = time.Now()

	var previousItem types.CraigslistItem
	if err = idx.access.Upsert(ctx, item, &previousItem); err != nil {
		return
	}

	newRecord = previousItem.IsEmpty()
	return
}
