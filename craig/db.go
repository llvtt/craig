package craig

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	log "github.com/go-kit/kit/log"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
)

type DBClient interface {
	// Insert inserts a new RSS Item into the database.
	// returns false if the item existed already in the database, otherwise return true
	InsertSearchedItem(item *types.CraigslistItem) (bool, error)
}

type JsonDBClient struct {
	dbFile  string
	byUrl   map[string]*types.CraigslistItem
	byTitle map[string]*types.CraigslistItem
	logger  log.Logger
}

func NewDBClient(conf *types.CraigConfig, logger log.Logger) (DBClient, error) {
	var client DBClient
	switch conf.DBType {
	case "json":
		var jsonClient *JsonDBClient
		jsonClient = &JsonDBClient{conf.DBFile, make(map[string]*types.CraigslistItem), make(map[string]*types.CraigslistItem), logger}
		err := jsonClient.initDB()
		if err != nil {
			return nil, err
		}
		client = jsonClient
	case "sqlite":
		var sqlClient *SqliteClient
		sqlClient = &SqliteClient{
			conf.DBFile,
			nil,
			&logger,
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
			return utils.WrapError("could not deserialize record in db", err)
		}
		c.byUrl[record.Url] = &record
		c.byTitle[record.Title] = &record
	}
	return nil
}

func (c JsonDBClient) flushDB() error {
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
	err := c.flushDB()
	if err != nil {
		return false, utils.WrapError("Could not flush db when inserting item", err)
	}
	return true, nil
}
