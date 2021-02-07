package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNext(t *testing.T) {
	assertions := assert.New(t)
	records := []map[string]*dynamodb.AttributeValue{
		{"Name": &dynamodb.AttributeValue{S: aws.String("Jeff")}},
		{"Name": &dynamodb.AttributeValue{S: aws.String("Bert")}},
	}
	type person struct { Name string }
	var p person
	it := &DynamoAccessIterator{scannedItems: records}
	assertions.NoError(it.Next(&p))
	assertions.Equal("Jeff", p.Name)
	assertions.NoError(it.Next(&p))
	assertions.Equal("Bert", p.Name)
	assertions.Equal(IteratorExhausted, it.Next(&p))
	assertions.Equal(IteratorExhausted, it.Next(&p))
}
