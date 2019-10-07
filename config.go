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

type CraigConfig struct {
	Region   string             `json:"region"`
	Searches []CraigslistSearch `json:"searches"`
	DBType   string             `json:"db_type"`
	DBFile   string             `json:"db_file"`
}

func parseConfig(filename string) *CraigConfig {
	var config CraigConfig
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
