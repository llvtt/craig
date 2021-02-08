package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/log"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/types"
	"os"
)

type Handler struct {
}

func (h Handler) Invoke(ctx context.Context, event []byte) (result []byte, err error) {
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
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var cloudWatchEvent events.CloudWatchEvent
	if err := json.Unmarshal(event, &cloudWatchEvent); err != nil {
		return nil, fmt.Errorf("unrecognized event type: %v", event)
	}

	// TODO: load this from DB instead of hard coding here
	config := &types.CraigConfig{
		Region: "sfc",
		Searches: []types.CraigslistSearch{
			{Category: "zip", Terms: []string{}, Neighborhoods: []int{3}},
			{Category: "ata", Terms: []string{"end table", "lamp", "mirror", "queen bed"}},
		},
	}
	searcher, err := craigslist.NewSearcher(config, logger)
	if err != nil {
		return
	}

	err = searcher.Search(ctx)

	return
}

////////////
// entry-point for lambda functions.
////////////
func main() {
	fmt.Println("Scrape craig main")

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.StartHandler(Handler{})
}