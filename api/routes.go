package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/ping", s.healthcheck)

	// store
	mux.HandlerFunc(http.MethodGet, "/stores", http.HandlerFunc(s.discoverStores))
	mux.HandlerFunc(http.MethodGet, "/stores/{id}/items", http.HandlerFunc(s.listStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/buy", http.HandlerFunc(s.buyStoreItems))

	// inventory
	mux.Handler(http.MethodPost, "/users/{id}/stores", s.authenticate(http.HandlerFunc(s.createStore)))
	mux.Handler(http.MethodPost, "/users/{user_id}/stores/{store_id}/items", s.authenticate(http.HandlerFunc(s.addStoreItem)))
	mux.Handler(http.MethodGet, "/users/{user_id}/stores/{store_id}/items", s.authenticate(http.HandlerFunc(s.listOwnedStoreItems)))
	mux.Handler(http.MethodPatch, "/users/{user_id}/stores/{store_id}/items/{item_id}", s.authenticate(http.HandlerFunc(s.updateStoreItems)))
	mux.Handler(http.MethodPost, "/users/{user_id}/store/{store_id}/owners", http.HandlerFunc(s.addNewOwner))
	mux.Handler(http.MethodGet, "/users/{id}/stores", s.authenticate(http.HandlerFunc(s.listUserStores)))
	mux.Handler(http.MethodDelete, "/users/{user_id}/stores/{store_id}/items/{item_id}", s.authenticate(http.HandlerFunc(s.deleteStoreItems)))
	mux.HandlerFunc(http.MethodDelete, "/users/{user_id}/store/{store_id}/owners", http.HandlerFunc(s.deleteOwner))
	mux.Handler(http.MethodDelete, "/users/{user_id}/stores/{store_id}", s.authenticate(http.HandlerFunc(s.deleteStore)))

	// TODO:
	mux.HandlerFunc(http.MethodPatch, "/stores/{id}/freeze", http.HandlerFunc(s.freezeStore))
	mux.HandlerFunc(http.MethodPatch, "/stores/{id}/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/stores/{store_id}/items/{item_id}/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	// mux.HandlerFunc(http.MethodGet, "/stores/{name}", http.HandlerFunc(s.getStore))

	return s.recoverPanic(s.enableCORS(s.httpLogger(mux)))
}