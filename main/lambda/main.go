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

type Handler struct {
}

func (h Handler) Invoke(ctx context.Context, event []byte) ([]byte, error) {
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
	craig := initCraig()
	var cloudWatchEvent events.CloudWatchEvent
	if err := json.Unmarshal(event, &cloudWatchEvent); err == nil {
		result, err := craig.Search(ctx, cloudWatchEvent)
		return []byte(result), err
	} else {
		return nil, fmt.Errorf("unrecognized event type: %v", event)
	}
}

////////////
// entry-point for lambda functions.
////////////
func main() {
	fmt.Println("Scrape craig main")

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.StartHandler(Handler{})
}

func parseConfig(conf string) (*types.CraigConfig, error) {
	var config types.CraigConfig
	err := json.Unmarshal([]byte(conf), &config)
	if err != nil {
		return nil, utils.WrapError("could not parse config", err)
	}
	return &config, nil
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

