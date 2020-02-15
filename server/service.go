package server

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/types"
)

type CraigService interface {
	Search(ctx context.Context) error
	ListSearches(ctx context.Context) (*[]types.CraigslistSearch, error)
	Health(ctx context.Context) (string, error)
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) *service {
	return &service{logger: logger}
}

func (s *service) Search(ctx context.Context) error {
	level.Info(s.logger).Log("msg", "Called method: Search")
	// todo
	return nil
}

func (s *service) ListSearches(ctx context.Context) (*[]types.CraigslistSearch, error) {
	level.Info(s.logger).Log("msg", "Called method: ListSearches")
	// todo
	return &[]types.CraigslistSearch{}, nil
}

func (s *service) Health(ctx context.Context) (string, error) {
	level.Info(s.logger).Log("msg", "Called method: Health")
	return "healthy", nil
}
