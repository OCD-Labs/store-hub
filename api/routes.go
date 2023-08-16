package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := httprouter.New()

	fsysHandler := http.FileServer(http.FS(s.swaggerFiles))
	mux.Handler(http.MethodGet, "/api/v1/swagger/*any", http.StripPrefix("/api/v1/swagger/", fsysHandler))

	mux.HandlerFunc(http.MethodPost, "/api/v1/ping/:user_id/:store_id", s.healthcheck)

	// store
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores", http.HandlerFunc(s.discoverStores))
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores/:store_id/items", http.HandlerFunc(s.listStoreItems))
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores/:store_id/items/:item_id", http.HandlerFunc(s.getStoreItems))
	mux.Handler(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/buy", s.authenticate(http.HandlerFunc(s.buyStoreItems)))

	// inventory
	mux.Handler(http.MethodPost, "/api/v1/users/:user_id/stores", s.authenticate(http.HandlerFunc(s.createStore)))
	mux.Handler(http.MethodPost, "/api/v1/users/:user_id/stores/:store_id/items", s.authenticate(http.HandlerFunc(s.addStoreItem)))
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id/stores/:store_id/items", s.authenticate(http.HandlerFunc(s.listOwnedStoreItems)))
	mux.Handler(http.MethodPatch, "/api/v1/users/:user_id/stores/:store_id/items/:item_id", s.authenticate(http.HandlerFunc(s.updateStoreItems)))
	mux.Handler(http.MethodPost, "/api/v1/users/:user_id/store/:store_id/owners", s.authenticate(http.HandlerFunc(s.addNewOwner)))
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id/stores", s.authenticate(http.HandlerFunc(s.listUserStores)))
	mux.Handler(http.MethodDelete, "/api/v1/users/:user_id/stores/:store_id/items/:item_id", s.authenticate(http.HandlerFunc(s.deleteStoreItems)))
	mux.Handler(http.MethodDelete, "/api/v1/users/:user_id/store/:store_id/owners", s.authenticate(http.HandlerFunc(s.deleteOwner)))
	mux.Handler(http.MethodPatch, "/api/v1/users/:user_id/stores/:store_id", s.authenticate(http.HandlerFunc(s.updateStoreProfile)))
	mux.Handler(http.MethodDelete, "/api/v1/users/:user_id/stores/:store_id", s.authenticate(http.HandlerFunc(s.deleteStore)))

	// orders
	mux.Handler(http.MethodPost, "/api/v1/seller/orders", s.authenticate(http.HandlerFunc(s.createOrder)))
	mux.Handler(http.MethodGet, "/api/v1/seller/orders", s.authenticate(http.HandlerFunc(s.listSellerOrders)))
	mux.Handler(http.MethodGet, "/api/v1/seller/orders/:order_id", s.authenticate(http.HandlerFunc(s.getSellerOrder)))
	mux.Handler(http.MethodPatch, "/api/v1/seller/orders/:order_id", s.authenticate(http.HandlerFunc(s.updateSellerOrder)))

	// sales
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id/stores/:store_id/sales", s.authenticate(http.HandlerFunc(s.listAllSales)))
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id/stores/:store_id/sales/:sale_id", s.authenticate(http.HandlerFunc(s.getSale)))

	// user
	mux.HandlerFunc(http.MethodPost, "/api/v1/users", s.createUser)
	mux.HandlerFunc(http.MethodPost, "/api/v1/auth/login", s.login)
	mux.Handler(http.MethodPost, "/api/v1/auth/logout", s.authenticate(http.HandlerFunc(s.logout)))
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id", s.authenticate(http.HandlerFunc(s.getUser)))

	// TODO:
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/freeze", http.HandlerFunc(s.freezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	// mux.HandlerFunc(http.MethodGet, "/stores/{name}", http.HandlerFunc(s.getStore))

	return s.recoverPanic(s.enableCORS(s.httpLogger(mux)))
}
