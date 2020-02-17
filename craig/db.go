package craig

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"os"
	"strings"
)

type DBClient interface {
	// Insert inserts a new RSS Item into the database.
	// returns false if the item existed already in the database at the same price, otherwise return true
	InsertSearchedItem(item *types.CraigslistItem) (bool, error)

	// Inserts the item into the price log. If the item has appeared before but at a different price,
	// return a price drop object indicating how much the price has fallen
	InsertPrice(item *types.CraigslistItem) (*types.PriceDrop, error)
}

type JsonDBClient struct {
	dbFile          string
	priceLogFile    string
	byUrl           map[string]*types.CraigslistItem
	byTitle         map[string]*types.CraigslistItem
	priceLogByUrl   map[string]*types.CraigslistPriceLog
	priceLogByTitle map[string]*types.CraigslistPriceLog
	logger          log.Logger
}

func NewDBClient(conf *types.CraigConfig, logger log.Logger) (DBClient, error) {
	if !strings.HasSuffix(conf.DBDir, "/") {
		conf.DBDir = conf.DBDir + "/"
	}
	err := os.MkdirAll(conf.DBDir, 0755)
	if err != nil {
		return nil, utils.WrapError("could not create db dir: "+conf.DBDir, err)
	}

	var client DBClient
	switch conf.DBType {
	case "json":
		dbFile := conf.DBDir + "database.json"
		priceLogFile := conf.DBDir + "price_log.json"

		var jsonClient *JsonDBClient
		jsonClient = &JsonDBClient{
			dbFile,
			priceLogFile,
			make(map[string]*types.CraigslistItem),
			make(map[string]*types.CraigslistItem),
			make(map[string]*types.CraigslistPriceLog),
			make(map[string]*types.CraigslistPriceLog),
			logger}
		err := jsonClient.initDB()
		if err != nil {
			return nil, err
		}
		client = jsonClient
	case "sqlite":
		dbFile := conf.DBDir + "database.sqlite3"
		var sqlClient *SqliteClient
		sqlClient = &SqliteClient{
			dbFile,
			nil,
			logger,
		}
		if err := sqlClient.InitDB(); err != nil {
			return nil, err
		}
		client = sqlClient
	case "":
		return nil, errors.New("no db type specified. must specify db_type in config file")
	default:
		return nil, errors.New("invalid db type: " + conf.DBType)
	}
	return client, nil
}

func (c *JsonDBClient) initDB() error {
	err := c.initItems()
	if err != nil {
		return utils.WrapError("Could not init items db", err)
	}
	err = c.initPriceLog()
	if err != nil {
		return utils.WrapError("Could not init price log db", err)
	}
	return nil
}

func (c *JsonDBClient) initItems() error {
	file, err := os.OpenFile(c.dbFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return utils.WrapError("could not open db file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record types.CraigslistItem
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return utils.WrapError("could not deserialize record in item db", err)
		}
		c.byUrl[record.Url] = &record
		c.byTitle[record.Title] = &record
	}
	return nil
}

func (c *JsonDBClient) initPriceLog() error {
	file, err := os.OpenFile(c.priceLogFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return utils.WrapError("could not open db file: "+c.priceLogFile, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record types.CraigslistPriceLog
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			return utils.WrapError("could not deserialize record in price log db", err)
		}
		c.priceLogByUrl[record.Item.Url] = &record
		c.priceLogByTitle[record.Item.Title] = &record
	}
	return nil
}

func (c JsonDBClient) flushDB() error {
	err := c.flushItems()
	if err != nil {
		return utils.WrapError("Could not flush items db", err)
	}
	err = c.flushPriceLog()
	if err != nil {
		return utils.WrapError("Could not flush price log db", err)
	}
	return nil
}

func (c JsonDBClient) flushItems() error {
	file, err := os.Create(c.dbFile)
	if err != nil {
		return utils.WrapError("could not create db file", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, record := range c.byUrl {
		bytes, err := json.Marshal(&record)
		if err != nil {
			return utils.WrapError("could not json serialize record", err)
		}
		_, err = writer.Write(bytes)
		if err != nil {
			return utils.WrapError("could not write to db", err)
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return utils.WrapError("could not write to db", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return utils.WrapError("could not flush to db", err)
	}
	return nil
}

func (c JsonDBClient) flushPriceLog() error {
	file, err := os.Create(c.priceLogFile)
	if err != nil {
		return utils.WrapError("could not create price log db file", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, record := range c.priceLogByUrl {
		bytes, err := json.Marshal(&record)
		if err != nil {
			return utils.WrapError("could not json serialize price log record", err)
		}
		_, err = writer.Write(bytes)
		if err != nil {
			return utils.WrapError("could not write to price log db", err)
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			return utils.WrapError("could not write to price log db", err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return utils.WrapError("could not flush to price log db", err)
	}
	return nil
}

func (c *JsonDBClient) InsertSearchedItem(item *types.CraigslistItem) (bool, error) {
	// check to see if we've posted about this item already
	// if the item already exists in the database, return false and do nothing
	if _, ok := c.byUrl[item.Url]; ok {
		return false, nil
	} else if _, ok := c.byTitle[item.Title]; ok {
		return false, nil
	}
	c.byUrl[item.Url] = item
	c.byTitle[item.Title] = item
	err := c.flushItems()
	if err != nil {
		return false, utils.WrapError("Could not flush db when inserting item", err)
	}
	return true, nil
}

func (c *JsonDBClient) InsertPrice(item *types.CraigslistItem) (*types.PriceDrop, error) {
	var priceLog *types.CraigslistPriceLog
	if _, ok := c.priceLogByUrl[item.Url]; ok {
		priceLog = c.priceLogByUrl[item.Url]
	} else if _, ok := c.priceLogByTitle[item.Title]; ok {
		priceLog = c.priceLogByTitle[item.Title]
	} else {
		// no price log for this entry, add it to the price log
		priceLog = &types.CraigslistPriceLog{
			Item:         item,
			Prices:       []*types.CraigslistPriceEntry{{
				Price:       item.Price,
				PublishDate: item.PublishDate,
			}},
			MaxPrice:     item.Price,
			CurrentPrice: item.Price,
		}
		c.priceLogByUrl[item.Url] = priceLog
		c.priceLogByTitle[item.Title] = priceLog
		err := c.flushPriceLog()
		if err != nil {
			return nil, utils.WrapError("Could not flush price db when inserting item", err)
		}
		return nil, nil
	}

	if priceLog.CurrentPrice > item.Price {
		// there's been a price drop
		priceDrop := &types.PriceDrop{
			Item: item,
			CurrentPrice: item.Price,
			MaxPrice: priceLog.MaxPrice,
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
		level.Debug(c.logger).Log("msg", fmt.Sprintf("Added new price for item %s. Current price is %d. Price drop object: %v", item.Title, priceLog.CurrentPrice, priceDrop))
		err := c.flushPriceLog()
		if err != nil {
			return nil, utils.WrapError("Could not flush price db when inserting item", err)
		}
		return priceDrop, nil
	} else if priceLog.MaxPrice < item.Price {
		priceLog.MaxPrice = item.Price
		// prepend new max price to beginning of list
		priceLog.Prices = append([]*types.CraigslistPriceEntry{{Price: item.Price, PublishDate: item.PublishDate}}, priceLog.Prices...)
		level.Debug(c.logger).Log("msg", fmt.Sprintf("Updated max price for item %s. Max price is %d", item.Title, priceLog.MaxPrice))
		err := c.flushPriceLog()
		if err != nil {
			return nil, utils.WrapError("Could not flush price db when inserting item", err)
		}
	}

	return nil, nil
}
