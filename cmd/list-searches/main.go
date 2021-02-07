package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/db"
	"time"
)

type Search struct {
	Query   string    `dynamodbav:"query"`
	Created time.Time `dynamodbav:"created"`
}

var (
	searchesClient db.DataAccess
	sess           *session.Session
	dynamo         *dynamodb.DynamoDB
)

func init() {
	sess = session.Must(session.NewSession())
	dynamo = dynamodb.New(sess)
	searchesClient = db.NewDynamoAccess("searches", dynamo)
}

func handler(ctx context.Context) error {
	var (
		search         Search
		searchIterator db.Iterator
		err            error
	)
	for searchIterator, err = searchesClient.List(ctx); err == nil; err = searchIterator.Next(&search) {
		fmt.Printf("search: %+v", search)
	}
	if err != db.IteratorExhausted {
		return err
	}

	newSearch := Search{"blender", time.Now()}
	var replacedSearch Search
	if err = searchesClient.Upsert(ctx, &newSearch, &replacedSearch); err != nil {
		return err
	}

	if replacedSearch.Query == "" {
		fmt.Println("No search was replaced")
	} else {
		fmt.Printf("Search was replaced: %+v\n", replacedSearch)
	}
	fmt.Printf("Upserted: %+v\n", newSearch)

	return nil
}

func main() {
	lambda.Start(handler)
}
