package craigslist

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/slack"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
)

type Searcher interface {
	Search() error
}

type searcher struct {
	conf             *types.CraigConfig
	craigslistClient CraigslistClient
	imageScraper     ImageScraper
	slackClient      *slack.SlackClient
	//dbClient         db.DBClient
	logger           log.Logger
}

func NewSearcher(conf *types.CraigConfig, logger log.Logger) (Searcher, error) {
	craigslistClient := NewCraigslistClient("sfbay", logger)
	slackClient, err := slack.NewSlackClient(logger)
	if err != nil {
		return nil, utils.WrapError("could not initialize slack client", err)
	}
	//dbClient, err := db.NewDBClient(conf, logger)
	if err != nil {
		return nil, utils.WrapError("could not initialize searcher", err)
	}

	imageScraper := NewImageScraper(logger)

	return &searcher{
		conf: conf,
		craigslistClient: craigslistClient,
		imageScraper: imageScraper,
		slackClient: slackClient,
		//dbClient: dbClient,
		logger: logger,
	}, nil

}


func (s *searcher) Search() error {
	options := &SearchOptions{HasPicture: true, SubRegion: s.conf.Region}
	for _, search := range s.conf.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := s.craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			// gather all matching results from craigslist
			var newResults Listing
			var priceDrops []*types.PriceDrop
			listing, err := categoryClient.Search(term)
			if err != nil {
				return utils.WrapError(fmt.Sprintf("Could not search term: %s", term), err)
			}

			// check to see which results are new and which are not
			for _, result := range listing {
				if result == nil {
					level.Error(s.logger).Log("result was nil!")
					continue
				}
				// TODO get this logic working again with the dynamodb db impl
				//level.Debug(s.logger).Log("result", result)
				//inserted, err := s.dbClient.InsertSearchedItem(result)
				//if err != nil {
				//	return utils.WrapError("could not insert searched item", err)
				//}
				//if inserted {
				//	newResults = append(newResults, result)
				//}
				//
				//// check for price drops
				//priceDrop, err := s.dbClient.InsertPrice(result)
				//if err != nil {
				//	return utils.WrapError("could not insert price into db", err)
				//}
				//if priceDrop != nil {
				//	priceDrops = append(priceDrops, priceDrop)
				//}
			}


			// post any new results to slack
			if len(newResults) > 0 {
				var announcement string
				if len(term) > 0 {
					announcement = fmt.Sprintf("Found %d new items matching *%s* on my list!", len(newResults), term)
				} else {
					announcement = fmt.Sprintf("Found %d new *free* items on my list!", len(newResults))
				}
				err := s.slackClient.SendString(announcement)
				if err != nil {
					return utils.WrapError("caught error when sending slack message", err)
				}
				messagesSent := 0
				for _, result := range newResults {
					urls, err := s.imageScraper.GetImageUrls(result)
					if err != nil {
						return utils.WrapError("Could not send item to craigslist", err)
					}
					err = s.slackClient.SendItem(result, urls)
					if err != nil {
						level.Error(s.logger).Log("caught error when sending slack message", err)
						continue
					}
					messagesSent++
				}
				level.Info(s.logger).Log("msg", fmt.Sprintf("sent %d slack messages", messagesSent))
			}

			if len(priceDrops) > 0 {
				announcement := fmt.Sprintf("Found %d items with price drops! :fire: :money_with_wings: :fire: ", len(priceDrops))
				err := s.slackClient.SendString(announcement)
				if err != nil {
					return utils.WrapError("Could not send string to craigslist", err)
				}
				for _, priceDrop := range priceDrops {
					urls, err := s.imageScraper.GetImageUrls(priceDrop.Item)
					if err != nil {
						return utils.WrapError("Could not send item to craigslist", err)
					}
					err = s.slackClient.SendPriceDrop(priceDrop, urls)
					if err != nil {
						return utils.WrapError("Could not send price drop to craigslist", err)
					}
				}
			}
		}
	}
	return nil
}
