package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const tableStatement = `
create table if not exists searches (
  id auto primary key,
  region varchar,
  term varchar,
  category varchar,
  neighborhood int,
  created_at timestamp
);
`

type SqliteClient struct {
	dbFile string
	db     *sqlx.DB
	logger log.Logger
}

func NewSqliteClient() (client *SqliteClient, err error) {
	var db *sqlx.DB
	if db, err = sqlx.Connect("sqlite3", "/tmp/database.sqlite3"); err != nil {
		return nil, err
	}

	client = &SqliteClient{
		"/tmp/database.sqlite3",
		db,
		log.NewJSONLogger(os.Stdout),
	}
	return
}

func Handler(ctx context.Context, event interface{}) (string, error) {
	fmt.Sprintf("Handler invoked with input: %v\n", event)
	if client, err := NewSqliteClient(); err != nil {
		return "", err
	} else if _, err := client.db.Exec(tableStatement); err != nil {
		return "", err
	}
	return "database initialized!", nil
}

func main() {
	fmt.Println("Craig main")
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)
}
