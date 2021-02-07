package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/llvtt/craig/db"
	"github.com/slack-go/slack"
)

type Search struct {
	Query   string    `dynamodbav:"query"`
	Created time.Time `dynamodbav:"created"`
}

var (
	sess           *session.Session
	tableMgr       *db.DynamoDBAccessManager
	searchesClient db.DataAccess

	slacker            *slack.Client
	slackSigningSecret string
)

func init() {
	sess = session.Must(session.NewSession())
	tableMgr = db.NewDynamoDBAccessManager(dynamodb.New(sess))
	searchesClient = tableMgr.Table("searches")

	slacker = slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

func postMessage(ctx context.Context, messageText string) error {
	respChannel, respTs, err := slacker.PostMessageContext(ctx, "cltest",
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText(messageText, false))

	fmt.Println("respChannel", respChannel)
	fmt.Println("respTs", respTs)

	return err
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

	//httpRequest, err := slackbot.HttpRequest(&req)
	//if err != nil {
	//	return err
	//}
	// TODO: parse slash commands from httpRequest

	err = postMessage(ctx, "Looking for searches!")
	err = postMessage(ctx, fmt.Sprintf("request = %+v", req))

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
