package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/llvtt/craig/types"
)

type SearchRequest struct{}

type SearchReply struct {
	// todo
	Err string `json:"err,omitempty"`
}

type ListSearchesRequest struct{}

type ListSearchesReply struct {
	Searches []types.CraigslistSearch `json:"searches"`
	Err      string                   `json:"err,omitempty"`
}

type HealthRequest struct{}

type HealthReply struct {
	Status string `json:"status"`
	Err    string `json:"err,omitempty"`
}

func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SearchRequest
	return req, nil
}

func decodeListSearchesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ListSearchesRequest
	return req, nil
}

func decodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req HealthRequest
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
