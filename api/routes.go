package api

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := pat.New()

	mux.Add(http.MethodGet, "/ping", http.HandlerFunc(s.healthcheck))

	return mux
}