package server

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/llvtt/craig/types"
)

type Endpoints struct {
	SearchEndpoint       endpoint.Endpoint
	ListSearchesEndpoint endpoint.Endpoint
	HealthEndpoint       endpoint.Endpoint
}

func NewEndpoints(s CraigService, logger log.Logger) Endpoints {
	return Endpoints{
		SearchEndpoint: makeSearchEndpoint(s),
		ListSearchesEndpoint: makeListSearchesEndpoint(s,logger),
		HealthEndpoint: makeHealthEndpoint(s),
	}
}

func makeSearchEndpoint(s CraigService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// don't need anything from the request, just need to validate it
		_ = request.(SearchRequest)
		err = s.Search(ctx)
		if err != nil {
			return SearchReply{Err: err.Error()}, nil
		}
		return SearchReply{}, nil
	}
}

func makeListSearchesEndpoint(s CraigService, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		level.Info(logger).Log("msg", "In make list searches endpoint")
		// don't need anything from the request, just need to validate it
		_ = request.(ListSearchesRequest)
		searches, err := s.ListSearches(ctx)
		if err != nil {
			return ListSearchesReply{Searches: nil, Err: err.Error()}, nil
		}
		return ListSearchesReply{Searches: searches}, nil
	}
}

func makeHealthEndpoint(s CraigService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(HealthRequest)
		health, err := s.Health(ctx)
		if err != nil {
			return nil, err
		}
		return HealthReply{Status: health}, nil
	}
}


func (e Endpoints) Search(ctx context.Context) error {
	req := SearchRequest{}
	resp, err := e.SearchEndpoint(ctx, req)
	if err != nil {
		return err
	}
	reply := resp.(SearchReply)
	if reply.Err != "" {
		return errors.New(reply.Err)
	}
	return nil
}

func (e Endpoints) ListSearches(ctx context.Context) (*[]types.CraigslistSearch, error) {
	req := ListSearchesRequest{}
	resp, err := e.ListSearchesEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	reply := resp.(ListSearchesReply)
	if reply.Err != "" {
		return nil,errors.New(reply.Err)
	}
	return reply.Searches, nil
}

func (e Endpoints) Health(ctx context.Context) (string, error) {
	req := HealthRequest{}
	resp, err := e.HealthEndpoint(ctx, req)
	if err != nil {
		return "error", err
	}
	reply := resp.(HealthReply)
	if reply.Err != "" {
		return "error", errors.New(reply.Err)
	}
	return reply.Status, nil
}

