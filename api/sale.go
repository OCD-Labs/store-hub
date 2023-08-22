package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/rs/zerolog/log"
)

type listAllSalesQueryStr struct {
	ItemPriceStart    string    `querystr:"item_price_start"`
	ItemPriceEnd      string    `querystr:"item_price_end"`
	ItemName          string    `querystr:"item_name"`
	CustomerAccountID string    `querystr:"customer_account_id"`
	DeliveryDateStart time.Time `querystr:"delivery_date_start"`
	DeliveryDateEnd   time.Time `querystr:"delivery_date_end"`
	OrderDateStart    time.Time `querystr:"order_date_start"`
	OrderDateEnd      time.Time `querystr:"order_date_end"`
	Page              int       `querystr:"page" validate:"max=10000000"`
	PageSize          int       `querystr:"page_size" validate:"max=20"`
	Sort              string    `querystr:"sort"`
}

type listAllSalesPathVars struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// maps to endpoint "GET /users/{user_id}/stores/{store_id}/sales".
func (s *StoreHub) listStoreSales(w http.ResponseWriter, r *http.Request) {
	var reqQueryStr listAllSalesQueryStr
	if err := s.shouldBindQuery(w, r, &reqQueryStr); err != nil {
		return
	}

	if reqQueryStr.Page < 1 {
		reqQueryStr.Page = 1
	}
	if reqQueryStr.PageSize < 1 {
		reqQueryStr.PageSize = 15
	}
	if reqQueryStr.Sort == "" {
		reqQueryStr.Sort = "-id"
	}

	var pathVars listAllSalesPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	// authorise
	authPayload := s.contextGetToken(r)
	if pathVars.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	arg := db.ListAllSellerSalesParams{
		ItemPriceStart:    reqQueryStr.ItemPriceStart,
		ItemPriceEnd:      reqQueryStr.ItemPriceEnd,
		ItemName:          reqQueryStr.ItemName,
		CustomerAccountID: reqQueryStr.CustomerAccountID,
		DeliveryDateStart: reqQueryStr.DeliveryDateStart,
		DeliveryDateEnd:   reqQueryStr.DeliveryDateEnd,
		OrderDateStart:    reqQueryStr.OrderDateStart,
		OrderDateEnd:      reqQueryStr.OrderDateEnd,
		StoreID:           pathVars.StoreID,
		SellerID:          authPayload.UserID,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"-id", "-item_name", "-delivery_date", "-order_date", "-price", "id", "item_name", "delivery_date", "order_date", "price"},
		},
	}

	sales, pagination, err := s.dbStore.ListAllSellerSales(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve sales")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	storeMetrics, err := s.dbStore.GetStoreMetrics(r.Context(), pathVars.StoreID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch store metrics")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some sales",
			"result": envelop{
				"sales":         sales,
				"metadata":      pagination,
				"store_metrics": storeMetrics,
			},
		},
	}, nil)
}

type getSalePathVars struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
	SaleID  int64 `path:"sale_id" validate:"required,min=1"`
}

// getSale maps to "GET /users/:user_id/stores/:store_id/sales/:sale_id"
func (s *StoreHub) getSale(w http.ResponseWriter, r *http.Request) {
	var pathVars getSalePathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	authPayload := s.contextGetToken(r)
	if pathVars.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	sale, err := s.dbStore.GetSale(r.Context(), db.GetSaleParams{
		StoreID:  pathVars.StoreID,
		SaleID:   pathVars.SaleID,
		SellerID: authPayload.UserID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "sale not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve sale details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found a sale",
			"result": envelop{
				"sale": sale,
			},
		},
	}, nil)
}

type listSalesOverviewPathVars struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

type listSalesOverviewQueryStr struct {
	ItemName     string `querystr:"item_name"`
	RevenueStart string `querystr:"revenue_start"`
	RevenueEnd   string `querystr:"revenue_end"`
	Page         int    `querystr:"page" validate:"max=10000000"`
	PageSize     int    `querystr:"page_size" validate:"max=20"`
	Sort         string `querystr:"sort"`
}

// listSalesOverview maps to "GET /users/:user_id/stores/:store_id/sales-overview"
func (s *StoreHub) listSalesOverview(w http.ResponseWriter, r *http.Request) {
	var pathVars listSalesOverviewPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVars.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	var reqQueryStr listSalesOverviewQueryStr
	if err := s.shouldBindQuery(w, r, &reqQueryStr); err != nil {
		return
	}

	if reqQueryStr.Page < 1 {
		reqQueryStr.Page = 1
	}
	if reqQueryStr.PageSize < 1 {
		reqQueryStr.PageSize = 15
	}
	if reqQueryStr.Sort == "" {
		reqQueryStr.Sort = "-revenue"
	}

	arg := db.SalesOverviewParams{
		ItemName:     reqQueryStr.ItemName,
		RevenueStart: reqQueryStr.RevenueStart,
		RevenueEnd:   reqQueryStr.RevenueEnd,
		StoreID:      pathVars.StoreID,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"number_of_sales", "revenue", "-number_of_sales", "-revenue"},
		},
	}

	saleOverview, pagination, err := s.dbStore.ListSalesOverview(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to list sales overview")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found your sales overview",
			"result": envelop{
				"sales_overview": saleOverview,
				"metadata":       pagination,
			},
		},
	}, nil)
}
