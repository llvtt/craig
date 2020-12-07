package http

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer is a good little http
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware) // @see https://stackoverflow.com/a/51456342

	r.Methods("POST").Path("/search").Handler(httptransport.NewServer(
		endpoints.SearchEndpoint,
		decodeSearchRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/searches").Handler(httptransport.NewServer(
		endpoints.ListSearchesEndpoint,
		decodeListSearchesRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoints.HealthEndpoint,
		decodeHealthRequest,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

