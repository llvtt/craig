package main

import (
	sqlite3 "github.com/mattn/go-sqlite3"
)

func (self *Client) InitTable() {
	createTableStmt := `
create table if not exists items (
  url varchar primary key,
  title varchar,
  thumbnail_url varchar,
  index_date timestamp,
  publish_date timestamp
)`
	if _, err := self.db.Exec(createTableStmt); err != nil {
		panic(err)
	}
}

// Insert inserts a new RSS Item into the database.
func (self *Client) Insert(item *CraigslistItem) bool {
	insertStmt := `
insert into items (title, url, thumbnail_url, index_date, publish_date)
values(:title, :url, :thumbnail_url, :index_date, :publish_date)
`
	_, err := self.db.NamedExec(insertStmt, item)
	if err != nil && err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}
