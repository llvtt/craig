package util

type iterationError string

func (ie iterationError) Error() string {
	return string(ie)
}

// IteratorExhausted is returned when an Iterator has reached the end of its documents.
const IteratorExhausted iterationError = "iterator exhausted"

type Iterator interface {
	// Unmarshal the next record from the Iterator.
	// Returns IteratorExhausted when there are no more records.
	// Returns any underlying error with unmarshalling.
	Next(out interface{}) error
}
