package server

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	craig "github.com/llvtt/craig/craig-core"
	"github.com/llvtt/craig/types"
	"github.com/llvtt/craig/utils"
)

type CraigService interface {
	Search(ctx context.Context) error
	ListSearches(ctx context.Context) ([]types.CraigslistSearch, error)
	Health(ctx context.Context) (string, error)
}

type service struct {
	config   *types.CraigConfig
	logger   log.Logger
	searcher craig.Searcher
}

func NewService(config *types.CraigConfig, logger log.Logger) (CraigService, error) {
	searcher, err := craig.NewSearcher(config, logger)
	if err != nil {
		return nil, utils.WrapError("could not initialize craig service", err)
	}
	return &service{config: config, logger: logger, searcher: searcher}, nil
}

func (s *service) Search(ctx context.Context) error {
	level.Info(s.logger).Log("msg", "Called method: Search")
	return s.searcher.Search()
}

func (s *service) ListSearches(ctx context.Context) ([]types.CraigslistSearch, error) {
	level.Info(s.logger).Log("msg", "Called method: ListSearches")
	return s.config.Searches, nil
}

func (s *service) Health(ctx context.Context) (string, error) {
	level.Info(s.logger).Log("msg", "Called method: Health")
	return "healthy", nil
}
