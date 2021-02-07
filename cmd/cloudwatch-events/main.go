package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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
	var cloudWatchEvent events.CloudWatchEvent
	if err := json.Unmarshal(event, &cloudWatchEvent); err == nil {
		fmt.Println(cloudWatchEvent)
		str := "it works!"
		return []byte(str), nil
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