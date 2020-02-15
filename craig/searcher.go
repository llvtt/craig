package craig

import (
	"fmt"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/slack"
	"github.com/llvtt/craig/types"
)

func Search(conf *types.CraigConfig) error {
	craigslistClient := craigslist.NewCraigslistClient("sfbay")
	slackClient := slack.NewSlackClient()
	dbClient := NewDBClient(conf)

	options := &craigslist.SearchOptions{HasPicture: true, SubRegion: conf.Region}
	for _, search := range conf.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			var newResults craigslist.Listing
			for _, result := range categoryClient.Search(term) {
				if dbClient.InsertSearchedItem(result) {
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
	// TODO proper error handling
	return nil
}
