package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type DBClient interface {
	initDB()
	flushDB()
	InsertSearchedItem(item *CraigslistItem) bool
}

type JsonDBClient struct {
	dbFile  string
	byUrl   map[string]*CraigslistItem
	byTitle map[string]*CraigslistItem
}

func NewDBClient(conf *CraigConfig) DBClient {
	var client DBClient
	switch conf.DBType {
	case "json":
		client = JsonDBClient{conf.DBFile, make(map[string]*CraigslistItem), make(map[string]*CraigslistItem)}
	case "":
		log.Fatal("No db type specified. Must specify db_type in config file.")
	default:
		log.Fatal("Invalid db type: " + conf.DBType)
	}
	client.initDB()
	return client
}

func (self JsonDBClient) initDB() {
	file, _ := os.Open(self.dbFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record CraigslistItem
		panicOnErr(json.Unmarshal(scanner.Bytes(), &record))
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
		panicOnErr(writer.Write(bytes))
		panicOnErr(writer.WriteString("\n"))
	}
	writer.Flush()
}

// Insert inserts a new RSS Item into the database.
func (self JsonDBClient) InsertSearchedItem(item *CraigslistItem) bool {
	if _, ok := self.byUrl[item.Url]; ok {
		return false
	} else if _, ok := self.byTitle[item.Title]; ok {
		return false
	}
	self.byUrl[item.Url] = item
	self.byTitle[item.Title] = item
	return true
}
