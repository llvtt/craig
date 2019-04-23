package main

import (
	"bufio"
	"encoding/json"
	"os"
)

func (self *Client) initDB() {
	file, _ := os.Open("database.json")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record CraigslistItem
		panicOnErr(json.Unmarshal(scanner.Bytes(), &record))
		self.byUrl[record.Url] = &record
		self.byTitle[record.Title] = &record
	}
}

func (self *Client) flushDB() {
	file, _ := os.Create("database.json")
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
func (self *Client) Insert(item *CraigslistItem) bool {
	if _, ok := self.byUrl[item.Url]; ok {
		return false
	} else if _, ok := self.byTitle[item.Title]; ok {
		return false
	}
	self.byUrl[item.Url] = item
	self.byTitle[item.Title] = item
	return true
}
