package craig_core

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/llvtt/craig/craigslist"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/slack"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
)

type Searcher interface {
	Search() error
}

type searcher struct {
	conf             *types.CraigConfig
	craigslistClient craigslist.CraigslistClient
	slackClient      *slack.SlackClient
	dbClient         db.DBClient
	logger           log.Logger
}

func NewSearcher(conf *types.CraigConfig, logger log.Logger) (Searcher, error) {
	craigslistClient := craigslist.NewCraigslistClient("sfbay", logger)
	slackClient, err := slack.NewSlackClient(logger)
	if err != nil {
		return nil, utils.WrapError("could not initialize slack client", err)
	}
	dbClient, err := db.NewDBClient(conf, logger)
	if err != nil {
		return nil, utils.WrapError("could not initialize searcher", err)
	}

	return &searcher{
		conf: conf,
		craigslistClient: craigslistClient,
		slackClient: slackClient,
		dbClient: dbClient,
		logger: logger,
	}, nil

}


func (s *searcher) Search() error {
	options := &craigslist.SearchOptions{HasPicture: true, SubRegion: s.conf.Region}
	for _, search := range s.conf.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := s.craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			// gather all matching results from craigslist
			var newResults craigslist.Listing
			var priceDrops []*types.PriceDrop
			listing, err := categoryClient.Search(term)
			if err != nil {
				return utils.WrapError(fmt.Sprintf("Could not search term: %s", term), err)
			}

			// check to see which results are new and which are not
			for _, result := range listing {
				inserted, err := s.dbClient.InsertSearchedItem(result)
				if err != nil {
					return utils.WrapError("could not insert searched item", err)
				}
				if inserted {
					newResults = append(newResults, result)
				}

				// check for price drops
				priceDrop, err := s.dbClient.InsertPrice(result)
				if err != nil {
					return utils.WrapError("could not insert price into db", err)
				}
				if priceDrop != nil {
					priceDrops = append(priceDrops, priceDrop)
				}
			}


			// post any new results to slack
			if len(newResults) > 0 {
				var announcement string
				if len(term) > 0 {
					announcement = fmt.Sprintf("Found %d new items matching *%s* on my list!", len(newResults), term)
				} else {
					announcement = fmt.Sprintf("Found %d new *free* items on my list!", len(newResults))
				}
				s.slackClient.SendString(announcement)
				for _, result := range newResults {
					s.slackClient.SendItem(result)
				}
			}

			if len(priceDrops) > 0 {
				announcement := fmt.Sprintf("Found %d items with price drops! :fire: :money_with_wings: :fire: ", len(priceDrops))
				s.slackClient.SendString(announcement)
				for _, priceDrop := range priceDrops {
					s.slackClient.SendPriceDrop(priceDrop)
				}
			}
		}
	}
	return nil
}
