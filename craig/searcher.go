package craig

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/slack"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
)

func Search(conf *types.CraigConfig, logger log.Logger) error {
	craigslistClient := craigslist.NewCraigslistClient("sfbay", logger)
	slackClient := slack.NewSlackClient(logger)
	dbClient, err := NewDBClient(conf, logger)
	if err != nil {
		return utils.WrapError("could not perform search", err)
	}

	options := &craigslist.SearchOptions{HasPicture: true, SubRegion: conf.Region}
	for _, search := range conf.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			var newResults craigslist.Listing
			for _, result := range categoryClient.Search(term) {
				inserted, err := dbClient.InsertSearchedItem(result)
				if err != nil {
					return utils.WrapError("could not insert searched item", err)
				}
				if inserted {
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
	return nil
}
