package craig

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)


type DBTestSuite struct {
	suite.Suite
	ctx    context.Context
	logger log.Logger
}

func (suite *DBTestSuite) SetupTest()  {
	suite.logger = log.With(log.NewJSONLogger(os.Stdout), "app", "flow")
}


func TestDBTestSuite(t *testing.T) {
	suite.Run(t, &DBTestSuite{
		ctx: context.Background(),
	})
}

func (suite *DBTestSuite) TestNew_DB_Client() {
	//c, err := NewDBClient(nil, suite.logger)


}

