package craigslist

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/db"
	"github.com/llvtt/craig/slack"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"os"
)

type Searcher interface {
	Search(ctx context.Context) error
}

type searcher struct {
	conf             *types.CraigConfig
	craigslistClient CraigslistClient
	imageScraper     ImageScraper
	slacker          *slack.Slacker
	itemDbAccess     db.DataAccess
	priceDbAccess    db.DataAccess
	logger           log.Logger
}

const (
	defaultSlackChannel = "cltest"
)

func NewSearcher(conf *types.CraigConfig, logger log.Logger) (Searcher, error) {
	craigslistClient := NewCraigslistClient("sfbay", logger)
	slackChannel := os.Getenv("SLACK_CHANNEL")
	if len(slackChannel) == 0 {
		slackChannel = defaultSlackChannel
	}
	slacker := slack.NewSlacker(slackChannel)
	sess := session.Must(session.NewSession())
	dbMgmr := db.NewDynamoDBAccessManager(dynamodb.New(sess))
	itemDbAccess := dbMgmr.Table("items")
	priceDbAccess := dbMgmr.Table("price_logs")
	imageScraper := NewImageScraper(logger)

	return &searcher{
		conf:             conf,
		craigslistClient: craigslistClient,
		imageScraper:     imageScraper,
		slacker:          slacker,
		itemDbAccess:     itemDbAccess,
		priceDbAccess:    priceDbAccess,
		logger:           logger,
	}, nil

}


func (s *searcher) Search(ctx context.Context) error {
	options := &SearchOptions{HasPicture: true, SubRegion: s.conf.Region}
	for _, search := range s.conf.Searches {
		options.Neighborhoods = search.Neighborhoods
		categoryClient := s.craigslistClient.Category(search.Category).Options(options)
		for _, term := range search.Terms {
			// gather all matching results from craigslist
			var newResults Listing
			var priceDrops []*types.PriceDrop
			level.Info(s.logger).Log("msg", "searching term: "+term)
			listing, err := categoryClient.Search(term)
			if err != nil {
				return utils.WrapError(fmt.Sprintf("Could not search term: %s", term), err)
			}

			// check to see which results are new and which are not
			for _, item := range listing {
				if item == nil {
					level.Error(s.logger).Log("msg", "item was nil!")
					continue
				}
				newItem, err := s.insertItem(ctx, item)
				if err != nil {
					return err
				}
				if newItem {
					newResults = append(newResults, item)
				}

				priceDrop, err := s.insertPrice(ctx, item, newItem)
				if err != nil {
					return err
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
				err := s.slacker.PostMessage(ctx, announcement)
				if err != nil {
					return utils.WrapError("caught error when sending slack message", err)
				}
				messagesSent := 0
				for _, result := range newResults {
					urls, err := s.imageScraper.GetImageUrls(result)
					if err != nil {
						return utils.WrapError("Could not send item to craigslist", err)
					}
					// TODO format image posts properly
					err = s.slacker.PostMessage(ctx, "item: %v, image_urls: %v", result, urls)
					if err != nil {
						level.Error(s.logger).Log("caught error when sending slack message", err)
						continue
					}
					messagesSent++
				}
				level.Info(s.logger).Log("msg", fmt.Sprintf("sent %d slack messages", messagesSent))
			}

			if len(priceDrops) > 0 {
				err := s.slacker.PostMessage(ctx, "Found %d items with price drops! :fire: :money_with_wings: :fire: ", len(priceDrops))
				if err != nil {
					return utils.WrapError("Could not send string to craigslist", err)
				}
				for _, priceDrop := range priceDrops {
					urls, err := s.imageScraper.GetImageUrls(priceDrop.Item)
					if err != nil {
						return utils.WrapError("Could not send item to craigslist", err)
					}
					err = s.slacker.PostMessage(ctx, "priceDrop: %v, image_urls: %v", priceDrop, urls)
					if err != nil {
						return utils.WrapError("Could not send price drop to craigslist", err)
					}
				}
			}
		}
	}
	return nil
}

func (s *searcher) insertItem(ctx context.Context, item *types.CraigslistItem) (newItem bool, err error) {
	// insert item into item table
	level.Debug(s.logger).Log("result", item)
	previousRecord := &types.CraigslistItem{}
	err = s.itemDbAccess.Upsert(ctx, item, previousRecord)
	if err != nil {
		return false, utils.WrapError("could not insert searched item", err)
	}
	return previousRecord.Url == "", err
}

func (s *searcher) insertPrice(ctx context.Context, item *types.CraigslistItem, newItem bool) (priceDrop *types.PriceDrop, err error) {
	var priceLog *types.CraigslistPriceLog
	if newItem {
		level.Info(s.logger).Log("msg", "inserting price for new item", "item", item)
		// initialize price log for new item
		priceLog = NewPriceLog(item)
	} else {
		level.Info(s.logger).Log("msg", "updating price log for existing item", "item", item)
		// check for price drops
		priceLog = &types.CraigslistPriceLog{}
		priceLogQuery := &types.CraigslistPriceLogGet{ItemUrl: item.Url}
		err := s.priceDbAccess.Get(ctx, priceLogQuery, priceLog)
		if err != nil {
			return nil, utils.WrapError("could not get price log", err)
		}

		if priceLog.ItemUrl == "" {
			// no price log exists for this item for some reason.
			priceLog = NewPriceLog(item)
		}

		level.Info(s.logger).Log("msg", "got price log for item", "priceLog", priceLog)

		if priceLog.CurrentPrice > item.Price {
			// there's been a price drop
			priceDrop = &types.PriceDrop{
				Item:          item,
				CurrentPrice:  item.Price,
				MaxPrice:      priceLog.MaxPrice,
				PreviousPrice: priceLog.CurrentPrice,
				// assumes the cheapest listing is the last listing in the log.
				// todo: actually iterate over the list and find the min or keep the list sorted
				PreviousPricePublishDate: priceLog.Prices[len(priceLog.Prices)-1].PublishDate,
				// assumes the max price is the first listing in the log.
				// todo: actually iterate over the list and find the max or keep the list sorted
				MaxPricePublishDate: priceLog.Prices[0].PublishDate,
			}
			priceLog.CurrentPrice = item.Price
			// append new low price to end of list
			priceLog.Prices = append(priceLog.Prices, &types.CraigslistPriceEntry{Price: item.Price, PublishDate: item.PublishDate})
			level.Debug(s.logger).Log("msg", fmt.Sprintf("Added new price for item %s. Current price is %d. Price drop object: %v", item.Title, priceLog.CurrentPrice, priceDrop))
		} else if priceLog.MaxPrice < item.Price {
			priceLog.MaxPrice = item.Price
			// prepend new max price to beginning of list
			priceLog.Prices = append([]*types.CraigslistPriceEntry{{Price: item.Price, PublishDate: item.PublishDate}}, priceLog.Prices...)
			level.Debug(s.logger).Log("msg", fmt.Sprintf("Updated max price for item %s. Max price is %d", item.Title, priceLog.MaxPrice))
		}
	}

	err = s.priceDbAccess.Upsert(ctx, priceLog)
	if err != nil {
		return nil, utils.WrapError("could not insert price log", err)
	}
	return nil, nil
}

func NewPriceLog(item *types.CraigslistItem) *types.CraigslistPriceLog {
	return &types.CraigslistPriceLog{
		ItemUrl: item.Url,
		Prices: []*types.CraigslistPriceEntry{
			{Price: item.Price, PublishDate: item.PublishDate},
		},
		CurrentPrice: item.Price,
		MaxPrice:     item.Price,
	}
}
