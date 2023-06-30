package api

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := pat.New()

	mux.Add(http.MethodGet, "/ping", http.HandlerFunc(s.healthcheck))

	mux.Add(http.MethodPost, "/stores", http.HandlerFunc(s.createStore))
	mux.Add(http.MethodGet, "/stores", http.HandlerFunc(s.discoverStore))
	mux.Add(http.MethodPatch, "/stores/{id}/freeze", http.HandlerFunc(s.freezeStore))
	mux.Add(http.MethodPatch, "/stores/{id}/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.Add(http.MethodDelete, "/stores/{id}/delete", http.HandlerFunc(s.deleteStore))

	mux.Add(http.MethodPost, "/stores/{id}/items", http.HandlerFunc(s.addStoreItem))
	mux.Add(http.MethodGet, "/stores/{id}/items", http.HandlerFunc(s.listStoreItems))
	mux.Add(http.MethodPatch, "/stores/{id}/items/{item_id}/update", http.HandlerFunc(s.updateStoreItems))
	mux.Add(http.MethodPatch, "/stores/{id}/items/{item_id}/buy", http.HandlerFunc(s.buyStoreItems))
	mux.Add(http.MethodPatch, "/stores/{id}/items/{item_id}/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.Add(http.MethodPatch, "/stores/{id}/items/{item_id}/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	mux.Add(http.MethodDelete, "/stores/{id}/items/{item_id}/delete", http.HandlerFunc(s.deleteStoreItems))

	mux.Add(http.MethodPost, "/store/{id}/owners", http.HandlerFunc(s.addNewOwner))
	mux.Add(http.MethodGet, "/store/{id}/owners", http.HandlerFunc(s.listOwners))
	mux.Add(http.MethodDelete, "/store/{id}/owners", http.HandlerFunc(s.deleteOwner))

	return mux
}