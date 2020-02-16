package craig

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/llvtt/craig/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

const (
	TEST_DB_FILE = "./test.db.json"
	TS_FORMAT    = time.RFC822
	INDEX_DATE   = "16 Feb 20 01:00 PST"
	PUBLISH_DATE = "16 Feb 20 00:00 PST"
)

type DBTestSuite struct {
	suite.Suite
	ctx    context.Context
	logger log.Logger
}

func (suite *DBTestSuite) SetupTest()  {
	suite.logger = log.With(log.NewJSONLogger(os.Stdout), "app", "flow")
	// clean up after previous runs if necessary
	os.Remove(TEST_DB_FILE)
}


func (suite *DBTestSuite) TearDownTest()  {
	os.Remove(TEST_DB_FILE)
}


func TestDBTestSuite(t *testing.T) {
	suite.Run(t, &DBTestSuite{
		ctx: context.Background(),
	})
}

func (suite *DBTestSuite) TestNew_DB_Client_Invalid_Config() {
	t := suite.T()
	conf := &types.CraigConfig{}
	c, err := NewDBClient(conf, suite.logger)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no db type specified. must specify db_type in config file")
	assert.Nil(t, c)

	conf = &types.CraigConfig{DBType: "derp"}
	c, err = NewDBClient(conf, suite.logger)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid db type: derp")
	assert.Nil(t, c)

	conf = &types.CraigConfig{DBType: "json"}
	c, err = NewDBClient(conf, suite.logger)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not open db file")
	assert.Nil(t, c)
}

func (suite *DBTestSuite) TestNew_DB_Client() {
	t := suite.T()
	conf := &types.CraigConfig{DBType: "json", DBFile: TEST_DB_FILE}
	c, err := NewDBClient(conf, suite.logger)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}


func (suite *DBTestSuite) TestDB_Inserts() {
	t := suite.T()
	conf := &types.CraigConfig{DBType: "json", DBFile: TEST_DB_FILE}
	c, err := NewDBClient(conf, suite.logger)
	assert.Nil(t, err)

	indexDate, _ := time.Parse(TS_FORMAT, INDEX_DATE)
	publishDate, _ := time.Parse(TS_FORMAT, PUBLISH_DATE)
	// test inserting an item
	item := &types.CraigslistItem{
		Url: "fake url",
		Title: "fake title",
		Description: "fake description",
		ThumbnailUrl: "thumbnail url",
		IndexDate: indexDate,
		PublishDate: publishDate,
		Price: 10.0,
	}
	inserted, err := c.InsertSearchedItem(item)
	assert.NoError(t, err)
	assert.True(t, inserted)

	// try inserting another item with the same title, it should not insert anything
	item2 := &types.CraigslistItem{
		Url: "fake url2",
		Title: "fake title",
		Description: "fake description. different listing but same title",
		ThumbnailUrl: "thumbnail url",
		IndexDate: indexDate,
		PublishDate: publishDate,
		Price: 10.0,
	}
	inserted, err = c.InsertSearchedItem(item2)
	assert.NoError(t, err)
	assert.False(t, inserted)

	// try inserting another item with the same url, it should not insert anything
	item3 := &types.CraigslistItem{
		Url: "fake url",
		Title: "fake title, different listing same url",
		Description: "fake description",
		ThumbnailUrl: "thumbnail url",
		IndexDate: indexDate,
		PublishDate: publishDate,
		Price: 10.0,
	}
	inserted, err = c.InsertSearchedItem(item3)
	assert.NoError(t, err)
	assert.False(t, inserted)
}

