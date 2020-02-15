package craig

import (
	"bufio"
	"encoding/json"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	"log"
	"os"
)

type DBClient interface {
	// Insert inserts a new RSS Item into the database.
	// returns false if the item existed already in the database, otherwise return true
	InsertSearchedItem(item *types.CraigslistItem) bool
}

type JsonDBClient struct {
	dbFile  string
	byUrl   map[string]*types.CraigslistItem
	byTitle map[string]*types.CraigslistItem
}

func NewDBClient(conf *types.CraigConfig) DBClient {
	var client DBClient
	switch conf.DBType {
	case "json":
		var jsonClient JsonDBClient
		jsonClient = JsonDBClient{conf.DBFile, make(map[string]*types.CraigslistItem), make(map[string]*types.CraigslistItem)}
		jsonClient.initDB()
		client = jsonClient
	case "":
		log.Fatal("No db type specified. Must specify db_type in config file.")
	default:
		log.Fatal("Invalid db type: " + conf.DBType)
	}
	return client
}

func (self JsonDBClient) initDB() {
	file, _ := os.Open(self.dbFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record types.CraigslistItem
		utils.PanicOnErr(json.Unmarshal(scanner.Bytes(), &record))
		self.byUrl[record.Url] = &record
		self.byTitle[record.Title] = &record
	}
}

func (self JsonDBClient) flushDB() {
	file, _ := os.Create(self.dbFile)
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, record := range self.byUrl {
		bytes, _ := json.Marshal(&record)
		utils.PanicOnErr(writer.Write(bytes))
		utils.PanicOnErr(writer.WriteString("\n"))
	}
	writer.Flush()
}

func (self JsonDBClient) InsertSearchedItem(item *types.CraigslistItem) bool {
	if _, ok := self.byUrl[item.Url]; ok {
		return false
	} else if _, ok := self.byTitle[item.Title]; ok {
		return false
	}
	self.byUrl[item.Url] = item
	self.byTitle[item.Title] = item
	self.flushDB()
	return true
}
