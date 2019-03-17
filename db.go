package main

import (
	sqlite3 "github.com/mattn/go-sqlite3"
)

func (self *Client) initTable() {
	createTableStmt := `
create table if not exists items (
  url varchar primary key,
  title varchar,
  thumbnail_url varchar,
  index_date timestamp,
  publish_date timestamp
);

create unique index if not exists unique_title on items (title);
`
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
	if err == nil {
		return true
	}
	if sqliteErr, ok := err.(sqlite3.Error); !ok {
		panic(err)
	} else {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintPrimaryKey:
			fallthrough
		case sqlite3.ErrConstraintUnique:
			return false
		}
	}
	panic(err)
}
