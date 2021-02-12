package db

import (
	"context"
	"github.com/llvtt/craig/internal/util"
)

type DataAccess interface {
	// List returns an Iterator over fetched items.
	List(context.Context) (util.Iterator, error)
	// Upsert upserts the `input` record and unmarshalls any overwritten record into `output`.
	Upsert(ctx context.Context, input interface{}, output ...interface{}) error
}
