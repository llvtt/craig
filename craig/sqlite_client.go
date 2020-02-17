package craig

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
	sqlite3 "github.com/mattn/go-sqlite3"
)

const initializeDbStatements = `
create table if not exists items (
  url varchar primary key,
  title varchar,
  thumbnail_url varchar,
  index_date timestamp,
  publish_date timestamp,
  price int
);

create unique index on items (url, title);

create table if not exists search_terms (
  term varchar,
  neighborhood varchar,
  category varchar
);

create unique index on search_terms (term, neighborhood, category);
`

const insertStmt = `
insert into items (title, url, thumbnail_url, index_date, publish_date, price)
values(:title, :url, :thumbnail_url, :index_date, :publish_date, :price)
`

type SqliteClient struct {
	dbFile string
	db     *sqlx.DB
	logger *log.Logger
}

func wrapSqlError(statement string, err error) error {
	return utils.WrapError(
		fmt.Sprintf("error executing statement: %s", statement),
		err)
}

func (c *SqliteClient) InitDB() (err error) {
	if c.db, err = sqlx.Connect("sqlite3", c.dbFile); err != nil {
		return utils.WrapError("could not open db file", err)
	}
	if _, err = c.db.Exec(initializeDbStatements); err != nil {
		return wrapSqlError(initializeDbStatements, err)
	}
	return nil
}

func (c *SqliteClient) InsertSearchedItem(item *types.CraigslistItem) (bool, error) {
	_, err := c.db.NamedExec(insertStmt, item)
	if err == nil {
		return true, nil
	}
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintPrimaryKey:
			fallthrough
		case sqlite3.ErrConstraintUnique:
			return false, nil
		}
	}
	return false, utils.WrapError("could not insert item", err)
}
