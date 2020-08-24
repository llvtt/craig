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

// TODO load conf from conf file
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

func SearchHandler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
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
	return initCraig().Search(ctx, event)
}

func GatewayHandler(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	// TODO
	// Parse API gateway events and respond accordingly
	return "", nil
}

func initCraig() clambda.LambdaServer {
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
	return clambda.NewLambdaServer(logger, svc)
}

// entry-point for scrape lambda function.
func main() {
	fmt.Println("Scrape craig main")

	// Lambda only allows specifying one handler per lambda function
	// we'll need a different entry-point function to respond to api gateway requests

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(SearchHandler)
}

// entry-point for gateway lambda function.
func gatewayMain() {
	fmt.Println("gateway craig main")

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(GatewayHandler)
}

func parseConfig(conf string) (*types.CraigConfig, error) {
	var config types.CraigConfig
	err := json.Unmarshal([]byte(conf), &config)
	if err != nil {
		return nil, utils.WrapError("could not parse config", err)
	}
	return &config, nil
}
