package server

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/craig"
	"github.com/llvtt/craig/types"
)

type CraigService interface {
	Search(ctx context.Context) error
	ListSearches(ctx context.Context) (*[]types.CraigslistSearch, error)
	Health(ctx context.Context) (string, error)
}

type service struct {
	config *types.CraigConfig
	logger log.Logger
}

func NewService(config *types.CraigConfig, logger log.Logger) *service {
	return &service{config: config, logger: logger}
}

func (s *service) Search(ctx context.Context) error {
	level.Info(s.logger).Log("msg", "Called method: Search")
	return craig.Search(s.config, s.logger)
}

func (s *service) ListSearches(ctx context.Context) (*[]types.CraigslistSearch, error) {
	level.Info(s.logger).Log("msg", "Called method: ListSearches")
	return &s.config.Searches, nil
}

func (s *service) Health(ctx context.Context) (string, error) {
	level.Info(s.logger).Log("msg", "Called method: Health")
	return "healthy", nil
}
