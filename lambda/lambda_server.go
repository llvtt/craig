package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/llvtt/craig/server"
)

type LambdaServer interface {
	Search(ctx context.Context, event events.CloudWatchEvent) (string, error)
}

type lambdaServer struct {
	svc server.CraigService
}

func NewLambdaServer(svc server.CraigService) LambdaServer {
	return &lambdaServer{svc: svc}
}

func (s *lambdaServer) Search(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	fmt.Printf("%v", event.Detail)
	err := s.svc.Search(ctx)
	if err != nil {
		fmt.Printf("Caught err while searching craigslit %s", err.Error())
	}
	return "", nil
}
