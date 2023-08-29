package api

import (
	"net/http"

	"github.com/OCD-Labs/store-hub/util"
	"github.com/julienschmidt/httprouter"
)

func (s *StoreHub) setupRoutes() http.Handler {
	mux := httprouter.New()

	fsysHandler := http.FileServer(http.FS(s.swaggerFiles))
	mux.Handler(http.MethodGet, "/api/v1/swagger/*any", http.StripPrefix("/api/v1/swagger/", fsysHandler))

	mux.HandlerFunc(http.MethodPost, "/api/v1/ping/:user_id/:store_id", s.healthcheck)

	// storefront
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores", http.HandlerFunc(s.discoverStores))
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores/:store_id/items", http.HandlerFunc(s.listStoreItems))
	mux.HandlerFunc(http.MethodGet, "/api/v1/stores/:store_id/items/:item_id", http.HandlerFunc(s.getStoreItems))
	mux.Handler(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/buy", s.authenticate(http.HandlerFunc(s.buyStoreItems)))

	// inventory
	mux.Handler(http.MethodPost, "/api/v1/inventory/stores", s.authenticate(http.HandlerFunc(s.createStore)))
	mux.Handler(http.MethodGet, "/api/v1/inventory/stores", s.authenticate(http.HandlerFunc(s.listUserStores)))
	mux.Handler(
		http.MethodPost,
		"/api/v1/inventory/stores/:store_id/items",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.PRODUCTINVENTORYACCESS,
			)(
				http.HandlerFunc(s.addStoreItem),
			),
		),
	)
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/items",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.PRODUCTINVENTORYACCESS,
			)(
				http.HandlerFunc(s.listOwnedStoreItems),
			),
		),
	)
	mux.Handler(
		http.MethodPatch,
		"/api/v1/inventory/stores/:store_id/items/:item_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.PRODUCTINVENTORYACCESS,
			)(
				http.HandlerFunc(s.updateStoreItems),
			),
		),
	)
	mux.Handler(
		http.MethodDelete,
		"/api/v1/inventory/stores/:store_id/items/:item_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.PRODUCTINVENTORYACCESS,
			)(
				http.HandlerFunc(s.deleteStoreItems),
			),
		),
	)
	mux.Handler(
		http.MethodPatch,
		"/api/v1/inventory/stores/:store_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
			)(
				http.HandlerFunc(s.updateStoreProfile),
			),
		),
	)
	mux.Handler(
		http.MethodDelete,
		"/api/v1/inventory/stores/:store_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
			)(
				http.HandlerFunc(s.deleteStore),
			),
		),
	)

	mux.Handler(
		http.MethodPost,
		"/api/v1/inventory/stores/:store_id/send-access-invitation",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
			)(
				http.HandlerFunc(s.sendAccessInvitation),
			),
		),
	)
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/accept-access-invitation",
		s.authenticate(
			http.HandlerFunc(s.grantStoreAccess),
		),
	)
	mux.Handler(
		http.MethodPatch,
		"/api/v1/inventory/stores/:store_id/revoke-access",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
			)(
				http.HandlerFunc(s.revokeUserAccess),
			),
		),
	)
	mux.Handler(
		http.MethodDelete,
		"/api/v1/inventory/stores/:store_id/revoke-all-access",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
			)(
				http.HandlerFunc(s.revokeAllUserAccess),
			),
		),
	)

	// orders
	mux.Handler(http.MethodPost, "/api/v1/inventory/stores/:store_id/orders", s.authenticate(http.HandlerFunc(s.createOrder)))
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/orders",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.ORDERSACCESS,
			)(
				http.HandlerFunc(s.listSellerOrders),
			),
		),
	)
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/orders/:order_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.ORDERSACCESS,
			)(
				http.HandlerFunc(s.getSellerOrder),
			),
		),
	)
	mux.Handler(
		http.MethodPatch,
		"/api/v1/inventory/stores/:store_id/orders/:order_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.ORDERSACCESS,
			)(
				http.HandlerFunc(s.updateSellerOrder),
			),
		),
	)

	// sales
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/sales",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.SALESACCESS,
			)(
				http.HandlerFunc(s.listStoreSales),
			),
		),
	)
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/sales/:sale_id",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.SALESACCESS,
			)(
				http.HandlerFunc(s.getSale),
			),
		),
	)
	mux.Handler(
		http.MethodGet,
		"/api/v1/inventory/stores/:store_id/sales-overview",
		s.authenticate(
			s.CheckAccessLevel(
				util.FULLACCESS,
				util.SALESACCESS,
			)(
				http.HandlerFunc(s.listSalesOverview),
			),
		),
	)

	// user
	mux.HandlerFunc(http.MethodPost, "/api/v1/users", s.createUser)
	mux.HandlerFunc(http.MethodPost, "/api/v1/auth/login", s.login)
	mux.Handler(http.MethodPost, "/api/v1/auth/logout", s.authenticate(http.HandlerFunc(s.logout)))
	mux.Handler(http.MethodGet, "/api/v1/users/:user_id", s.authenticate(http.HandlerFunc(s.getUser)))

	// review
	mux.Handler(http.MethodPut, "/api/v1/accounts/:account_id/reviews/:order_id", s.authenticate(http.HandlerFunc(s.addReview)))

	// TODO:
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/freeze", http.HandlerFunc(s.freezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/unfreeze", http.HandlerFunc(s.unfreezeStore))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/freeze", http.HandlerFunc(s.freezeStoreItems))
	mux.HandlerFunc(http.MethodPatch, "/api/v1/stores/:store_id/items/:item_id/unfreeze", http.HandlerFunc(s.unfreezeStoreItems))
	// mux.HandlerFunc(http.MethodGet, "/stores/{name}", http.HandlerFunc(s.getStore))

	return s.recoverPanic(s.enableCORS(s.httpLogger(mux)))
}
