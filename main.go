package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CraigslistConfig struct {
	SlackEndpoint string   `json:"slackEndpoint"`
	SearchTerms   []string `json:"searchTerms"`
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

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s: must specify config file\n", os.Args[0])
		os.Exit(1)
	}
	config := parseConfig(os.Args[1])
	c := NewClient("sfbay")
	c.InitTable()
	for _, term := range config.SearchTerms {
		results := c.Category("ata").Options(&SearchOptions{HasPicture: true}).Search(term)
		var newResults Listing
		for _, result := range results {
			if c.Insert(result) {
				newResults = append(newResults, result)
			}
		}
		if len(newResults) > 0 {
			c.NotifySlack(config.SlackEndpoint, term, newResults)
		}
	}
}
