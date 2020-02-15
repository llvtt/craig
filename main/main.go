package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/llvtt/craig/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const DEFAULT_CONFIG_FILE_NAME = "config.json"

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
	)

	//configFilePath := flag.String(
	//	"config-file",
	//	DEFAULT_CONFIG_FILE_NAME,
	//	"The path to the config file.")

	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	flag.Parse()
	ctx := context.Background()
	svc := server.NewService(logger)

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

	//log.Printf("Config file path is: %v", *configFilePath)
	//craig.StartServer(craig.ParseConfig(*configFilePath))


	//server.NewHTTPServer()
	level.Error(logger).Log(<-errChan)
}
