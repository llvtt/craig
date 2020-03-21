package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	clambda "github.com/llvtt/craig/lambda"
	"github.com/llvtt/craig/server"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/log"
)

const conf = `
{
  "db_type": "json",
  "db_dir": "/tmp/craig_prod",
  "region": "sfc",
  "searches": [
    {
      "category": "zip",
      "terms": [""],
      "nh": [3]
    },
    {
      "category": "ata",
      "terms": ["end table", "lamp", "mirror", "queen bed"]
    }
  ]
}
`
func HandlerGeneric(ctx context.Context, event interface{}) (string, error) {
	// scrape craigslist trigger event looks like:
	/*
		map[account:860626312307
		    detail:map[]
		    detail-type:Scheduled Event
		    id:e137a93d-05ce-0978-3201-d6795dff8b30
		    region:us-west-2
		    resources:[arn:aws:events:us-west-2:860626312307:rule/ScrapeCraigslistTriggerRule]
		    source:aws.events
		    time:2020-03-21T17:04:58Z
		    version:0
		]
	*/
	fmt.Printf("Handler invoked with input: %v\n", event)
	fmt.Printf("input has type: %T\n", event)
	return "", nil
}

func Handler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	fmt.Printf("Handler invoked with input: %v\n", event)
	fmt.Printf("input has type: %T\n", event)
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	config, err := parseConfig(conf)
	if err != nil {
		panic(utils.WrapError("Could not start craig!", err).Error())
	}

	svc, err := server.NewService(config, logger)
	if err != nil {
		panic(utils.WrapError("Could not start craig!", err).Error())
	}

	lambdaServer := clambda.NewLambdaServer(svc)

	return lambdaServer.Search(ctx, event)
}

func main() {
	fmt.Println("Craig main")

	// Lambda only allows specifying one handler per lambda function
	// TODO figure out how we can have use the same binary for multiple different functions
	// we'll need a function to respond to api gateway requests as well as cloudwatch events
	// we could configure which handler is started with env variables?
	// or have the handler function take a generic event interface{} type and try to parse it?
	// lambda does not have a nice way to parse events AFAICT

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}

func parseConfig(conf string) (*types.CraigConfig, error) {
	var config types.CraigConfig
	err := json.Unmarshal([]byte(conf), &config)
	if err != nil {
		return nil, utils.WrapError("could not parse config", err)
	}
	return &config, nil
}
