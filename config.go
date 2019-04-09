package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type CraigslistSearch struct {
	Category      string   `json:"category"`
	Terms         []string `json:"terms"`
	Neighborhoods []int    `json:"nh"`
}

type CraigslistConfig struct {
	SlackEndpoint string             `json:"slackEndpoint"`
	Region        string             `json:"region"`
	Searches      []CraigslistSearch `json:"searches"`
}

func parseConfig(filename string) *CraigslistConfig {
	var config CraigslistConfig
	if file, err := os.Open(filename); err != nil {
		panic(err)
	} else if contents, err := ioutil.ReadAll(file); err != nil {
		panic(err)
	} else if err := json.Unmarshal(contents, &config); err != nil {
		panic(err)
	} else {
		return &config
	}
}
