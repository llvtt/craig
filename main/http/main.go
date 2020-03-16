package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/server"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const DEFAULT_CONFIG_FILE_NAME = "config.json"

func main() {
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)

	configFilePath := flag.String(
		"config-file",
		DEFAULT_CONFIG_FILE_NAME,
		"The path to the config file.")

	flag.Parse()

	level.Info(logger).Log("msg", "Loading configs from file " + *configFilePath)
	config, err := parseConfig(*configFilePath)
	if err != nil {
		panic(utils.WrapError("Could not start craig!", err).Error())
	}


	ctx := context.Background()
	svc, err := server.NewService(config, logger)
	if err != nil {
		panic(utils.WrapError("Could not start craig!", err).Error())
	}

	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	// mapping endpoints
	endpoints := server.NewEndpoints(svc, logger)

	// HTTP transport
	go func() {
		level.Info(logger).Log("msg", fmt.Sprintf("craig is listening on port: %s", *httpAddr))
		handler := server.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log(<-errChan)
}


func parseConfig(filename string) (*types.CraigConfig, error) {
	var config types.CraigConfig
	if file, err := os.Open(filename); err != nil {
		return nil, utils.WrapError("could not open config file: "+filename, err)
	} else if contents, err := ioutil.ReadAll(file); err != nil {
		return nil, utils.WrapError("could not read config file: "+filename, err)
	} else if err := json.Unmarshal(contents, &config); err != nil {
		return nil, utils.WrapError("could not parse config file: "+filename, err)
	} else {
		return &config, nil
	}
}
