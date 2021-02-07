package db

import "context"

type DataAccess interface {
	// List returns an Iterator over fetched items.
	List(context.Context) (Iterator, error)
	// Upsert upserts the `input` record and unmarshalls any overwritten record into `output`.
	Upsert(ctx context.Context, input interface{}, output ...interface{}) error
}

type dataAccessErr string

func (ie dataAccessErr) Error() string {
	return string(ie)
}

// IteratorExhausted is returned when an Iterator has reached the end of its documents.
const IteratorExhausted dataAccessErr = "iterator exhausted"

type Iterator interface {
	// Unmarshal the next record from the Iterator.
	// Returns IteratorExhausted when there are no more records.
	// Returns any underlying error with unmarshalling.
	Next(out interface{}) error
}
