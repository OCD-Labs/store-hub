package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/ping", s.healthcheck)

	// store
	mux.Handler(http.MethodPost, "/stores", s.authenticate(http.HandlerFunc(s.createStore)))
	mux.HandlerFunc(http.MethodPatch, "/stores/{id}/freeze", http.HandlerFunc(s.freezeStore))
	mux.HandlerFunc(http.MethodPatch, "/stores/{id}/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.HandlerFunc(http.MethodDelete, "/stores/{id}/delete", http.HandlerFunc(s.deleteStore))
	mux.HandlerFunc(http.MethodGet, "/stores", http.HandlerFunc(s.discoverStore))
	mux.HandlerFunc(http.MethodGet, "/stores/{id}/items", http.HandlerFunc(s.listStoreItems))

	// inventory
	mux.Handler(http.MethodPost, "/stores/{id}/items", s.authenticate(http.HandlerFunc(s.addStoreItem)))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/update", http.HandlerFunc(s.updateStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/buy", http.HandlerFunc(s.buyStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	mux.HandlerFunc(http.MethodDelete, "/stores/{store_id}/items/{item_id}/delete", http.HandlerFunc(s.deleteStoreItems))

	// store ownership
	mux.Handler(http.MethodPost, "/store/{id}/owners", http.HandlerFunc(s.addNewOwner))
	mux.HandlerFunc(http.MethodGet, "/store/{id}/owners", http.HandlerFunc(s.listOwners))
	mux.HandlerFunc(http.MethodDelete, "/store/{id}/owners", http.HandlerFunc(s.deleteOwner))

	// user
	mux.HandlerFunc(http.MethodGet, "/users/{id}/stores", http.HandlerFunc(s.discoverStoreByOwner))
	mux.Handler(http.MethodGet, "/users/{user_id}/stores/{store_id}/items", s.authenticate(http.HandlerFunc(s.listOwnedStoreItems)))

	return s.recoverPanic(s.enableCORS(s.httpLogger(mux)))
}