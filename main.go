package main

import (
	"fmt"
)

const CONFIG_FILE_NAME = "config.json"

func main() {
	config := parseConfig(CONFIG_FILE_NAME)
	craigslistClient := NewCraigslistClient("sfbay")
	slackClient := NewSlackClient()

	options := &SearchOptions{HasPicture: true, SubRegion: config.Region}
	for _, search := range config.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			var newResults Listing
			for _, result := range categoryClient.Search(term) {
				if craigslistClient.Insert(result) {
					newResults = append(newResults, result)
				}
			}
			if len(newResults) > 0 {
				announcement := fmt.Sprintf("Found %d new *free* items on my list!", len(newResults))
				if len(term) > 0 {
					announcement = fmt.Sprintf("Found %d new items matching *%s* on my list!", len(newResults), term)
				}
				slackClient.SendString(announcement)
				for _, result := range newResults {
					slackClient.SendItem(result)
				}
			}
		}
	}

	craigslistClient.flushDB()
}
