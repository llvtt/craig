package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/llvtt/craig/server"
	"github.com/mitchellh/mapstructure"
)

type LambdaServer interface {
	Handle(ctx context.Context, event interface{}) (string, error)
}

type lambdaServer struct {
	svc server.CraigService
}

func NewLambdaServer(svc server.CraigService) LambdaServer {
	return &lambdaServer{svc: svc}
}

func (s *lambdaServer) Handle(ctx context.Context, event interface{}) (string, error) {
	var e events.CloudWatchEvent
	err := mapstructure.Decode(event, &e)
	if err != nil {
		fmt.Printf("Caught err while parsing craigslist event %s", err.Error())
	}
	err = s.svc.Search(ctx)
	if err != nil {
		fmt.Printf("Caught err while searching craigslit %s", err.Error())
	}
	return "", nil
}
