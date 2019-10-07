package main

import (
	"flag"
	"log"
)

const DEFAULT_CONFIG_FILE_NAME = "config.json"

func main() {
	configFilePath := flag.String(
		"config-file",
		DEFAULT_CONFIG_FILE_NAME,
		"The path to the config file.")

	flag.Parse()

	log.Printf("Config file path is: %v", *configFilePath)
	startServer(parseConfig(*configFilePath))
}
