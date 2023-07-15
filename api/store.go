package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/rs/zerolog/log"
)

type discoverStoreQueryStr struct {
	StoreName string `querystr:"store_name"`
	Page      int    `querystr:"page" validate:"max=10000000"`
	PageSize  int    `querystr:"page_size" validate:"max=20"`
	Sort      string `querystr:"sort"`
}

// discoverStoreByOwner maps to endpoint "GET /stores?<query_string>"
func (s *StoreHub) discoverStores(w http.ResponseWriter, r *http.Request) {
	// parse request
	var reqQueryStr discoverStoreQueryStr
	if err := s.shouldBindQuery(w, r, &reqQueryStr); err != nil {
		return
	}

	// TODO: check out how to use struct tag to achieve same result
	if reqQueryStr.Page < 1 {
		reqQueryStr.Page = 1
	}
	if reqQueryStr.PageSize < 1 {
		reqQueryStr.PageSize = 15
	}
	if reqQueryStr.Sort == "" {
		reqQueryStr.Sort = "id"
	}

	// db query
	arg := db.ListAllStoresParams{
		Search: reqQueryStr.StoreName,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"id", "name", "-id", "-name"},
		},
	}

	stores, pagination, err := s.dbStore.ListAllStores(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve stores")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some stores",
			"result": envelop{
				"stores":   stores,
				"metadata": pagination,
			},
		},
	}, nil)
}

type listStoreItemsQueryStr struct {
	ItemName string `querystr:"item_name"` // TODO: add category field
	Page     int    `querystr:"page" validate:"max=10000000"`
	PageSize int    `querystr:"page_size" validate:"max=20"`
	Sort     string `querystr:"sort"`
}

type listStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

// listStoreItems maps to endpoint "GET /stores/{id}/items"
func (s *StoreHub) listStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar listStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request
	var reqQueryStr listStoreItemsQueryStr
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
		reqQueryStr.Sort = "id"
	}

	// db query
	arg := db.ListStoreItemsParams{
		StoreID:  pathVar.StoreID,
		ItemName: reqQueryStr.ItemName,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"id", "category", "price", "name", "-id", "-category", "-name", "-price"},
		},
	}
	items, pagination, err := s.dbStore.ListStoreItems(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve store items")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some store items",
			"result": envelop{
				"items":    items,
				"metadata": pagination,
			},
		},
	}, nil)
}

type buyStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	ItemID  int64 `path:"item_id" validate:"required,min=1"`
}

// buyStoreItems maps to endpoint "PATCH /stores/:store_id/items/:item_id/buy"
func (s *StoreHub) buyStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar buyStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	item, err := s.dbStore.GetItem(r.Context(), pathVar.ItemID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "item not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch item's details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if item.SupplyQuantity >= 1 {
		arg := db.UpdateItemParams{
			ItemID: pathVar.ItemID,
			SupplyQuantity: sql.NullInt64{
				Int64: item.SupplyQuantity - 1,
				Valid: true,
			},
		}
		updatedItem, err := s.dbStore.UpdateItem(r.Context(), arg)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to update item")
			log.Error().Err(err).Msg("error occurred")
			return
		}
		s.writeJSON(w, http.StatusOK, envelop{
			"status": "success",
			"data": envelop{
				"message": "item sold",
				"result": envelop{
					"new_item": updatedItem,
				},
			},
		}, nil)
	} else {
		s.errorResponse(w, r, http.StatusNotFound, "item no longer in stock")
	}

	// TODO: Add swagger documentation for this endpoint
}

type getStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	ItemID  int64 `path:"item_id" validate:"required,min=1"`
}

// maps to endpoint "GET /stores/:store_id/items/:item_id"
func (s *StoreHub) getStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar getStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	item, err := s.dbStore.GetItem(r.Context(), pathVar.ItemID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "item not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch item's details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found item",
			"result": envelop{
				"item": item,
			},
		},
	}, nil)
}

func (s *StoreHub) freezeStore(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) unfreezeStore(w http.ResponseWriter, r *http.Request) {

}
