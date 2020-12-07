package lambda

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/llvtt/craig/http"
)

type LambdaServer interface {
	Search(ctx context.Context, event events.CloudWatchEvent) (string, error)
}

type lambdaServer struct {
	logger log.Logger
	svc    http.CraigService
}

func NewLambdaServer(logger log.Logger, svc http.CraigService) LambdaServer {
	return &lambdaServer{logger: logger, svc: svc}
}

func (s *lambdaServer) Search(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	fmt.Printf("Handler invoked with input: %v\n", event)
	fmt.Printf("input has type: %T\n", event)
	err := s.svc.Search(ctx)
	if err != nil {
		fmt.Printf("Caught err while searching craigslit %s", err.Error())
	}
	return "", nil
}
