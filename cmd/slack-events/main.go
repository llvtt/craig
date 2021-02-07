package main

import (
	"context"
	"fmt"
	slackbot "github.com/llvtt/craig/slack"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/db"
)

type Search struct {
	Query   string    `dynamodbav:"query"`
	Created time.Time `dynamodbav:"created"`
}

const (
	defaultSlackChannel = "cltest"
)

var (
	sess           *session.Session
	tableMgr       *db.DynamoDBAccessManager
	searchesClient db.DataAccess

	slacker *slackbot.Slacker
)

func init() {
	sess = session.Must(session.NewSession())
	tableMgr = db.NewDynamoDBAccessManager(dynamodb.New(sess))
	searchesClient = tableMgr.Table("searches")

	slackChannel := os.Getenv("SLACK_CHANNEL")
	if len(slackChannel) == 0 {
		slackChannel = defaultSlackChannel
	}
	slacker = slackbot.NewSlacker(slackChannel)
}

// help [command]
// search create <"terms">, ["category", "region", "neighborhood"]
// search list
// search delete <id>

func handler(ctx context.Context, req events.APIGatewayProxyRequest) error {
	var (
		search         Search
		searchIterator db.Iterator
		err            error
	)

	slashCommand, err := slacker.ParseCommand(&req)
	if err != nil {
		return err
	}
	err = slacker.PostMessage(ctx, "received slash command: %+v", slashCommand)
	err = slacker.PostMessage(ctx, "Looking for searches!")
	err = slacker.PostMessage(ctx, "request = %+v", req)

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
