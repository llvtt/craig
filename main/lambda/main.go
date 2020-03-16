package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func Handler(ctx context.Context, event interface{}) (string, error) {
	fmt.Printf("Handler invoked with input: %v\n", event)
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

	return lambdaServer.Handle(ctx, event)
}

func main() {
	fmt.Println("Craig main")
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
