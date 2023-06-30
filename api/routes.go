package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := httprouter.New()

	mux.HandlerFunc(http.MethodGet, "/api/v1/ping", s.healthcheck)

	// store
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores", http.HandlerFunc(s.discoverStores))
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores/{id}/items", http.HandlerFunc(s.listStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/{store_id}/items/{item_id}/buy", http.HandlerFunc(s.buyStoreItems))

	// inventory
	mux.Handler(http.MethodPost, "/api/v1/users/{id}/stores", s.authenticate(http.HandlerFunc(s.createStore)))
	mux.Handler(http.MethodPost, "/api/v1/users/{user_id}/stores/{store_id}/items", s.authenticate(http.HandlerFunc(s.addStoreItem)))
	mux.Handler(http.MethodGet, "/api/v1/users/{user_id}/stores/{store_id}/items", s.authenticate(http.HandlerFunc(s.listOwnedStoreItems)))
	mux.Handler(http.MethodPatch, "/api/v1/users/{user_id}/stores/{store_id}/items/{item_id}", s.authenticate(http.HandlerFunc(s.updateStoreItems)))
	mux.Handler(http.MethodPost, "/api/v1/users/{user_id}/store/{store_id}/owners", http.HandlerFunc(s.addNewOwner))
	mux.Handler(http.MethodGet, "/api/v1/users/{id}/stores", s.authenticate(http.HandlerFunc(s.listUserStores)))
	mux.Handler(http.MethodDelete, "/api/v1/users/{user_id}/stores/{store_id}/items/{item_id}", s.authenticate(http.HandlerFunc(s.deleteStoreItems)))
	mux.HandlerFunc(http.MethodDelete, "/api/v1/users/{user_id}/store/{store_id}/owners", http.HandlerFunc(s.deleteOwner))
	mux.Handler(http.MethodDelete, "/api/v1/users/{user_id}/stores/{store_id}", s.authenticate(http.HandlerFunc(s.deleteStore)))

	// user
	mux.HandlerFunc(http.MethodPost, "/api/v1/users", s.createUser)
	mux.HandlerFunc(http.MethodGet, "/api/v1/auth/login", s.login)
	mux.Handler(http.MethodPost, "/api/v1/auth/logout", s.authenticate(http.HandlerFunc(s.logout)))
	mux.Handler(http.MethodGet, "/api/v1/users/{id}", s.authenticate(http.HandlerFunc(s.getUser)))

	// TODO:
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/{id}/freeze", http.HandlerFunc(s.freezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/{id}/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/{store_id}/items/{item_id}/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/{store_id}/items/{item_id}/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	// mux.HandlerFunc(http.MethodGet, "/stores/{name}", http.HandlerFunc(s.getStore))

	return s.recoverPanic(s.enableCORS(s.httpLogger(mux)))
}