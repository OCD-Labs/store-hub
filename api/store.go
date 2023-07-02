package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type discoverStoreQueryStr struct {
	StoreName string `json:"store_name"`
	Page      int    `json:"page" validate:"min=1,max=10000000"`
	PageSize  int    `json:"page_size" validate:"min=1,max=20"`
	Sort      string `json:"sort"`
}

// discoverStoreByOwner maps to endpoint "GET /stores?<query_string>"
func (s *StoreHub) discoverStores(w http.ResponseWriter, r *http.Request) {
	// parse request
	queryStr := r.URL.Query()
	var reqQueryStr discoverStoreQueryStr

	reqQueryStr.StoreName = s.readStr(queryStr, "store_name", "")
	reqQueryStr.Sort = s.readStr(queryStr, "sort", "id")

	reqQueryStr.Page, _ = s.readInt(queryStr, "page", 1)
	reqQueryStr.PageSize, _ = s.readInt(queryStr, "page_size", 15)

	// validate query string
	if err := s.bindJSONWithValidation(w, r, &reqQueryStr, validator.New()); err != nil {
		return
	}

	fmt.Printf("\n%+v\n", reqQueryStr)

	// db query
	arg := db.ListAllStoresParams{
		Search: reqQueryStr.StoreName,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"id", "store_name", "-id", "-store_name"},
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
	ItemName string `json:"item_name"` // TODO: add category field
	Page     int    `json:"page" validate:"min=1,max=10000000"`
	PageSize int    `json:"page_size" validate:"min=1,max=20"`
	Sort     string `json:"sort"`
}

type listStoreItemsPathVar struct {
	StoreID int64 `json:"store_id" validate:"required,min=1"`
}

// listStoreItems maps to endpoint "GET /stores/{id}/items"
func (s *StoreHub) listStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar listStoreItemsPathVar
	var err error

	// parse path variables
	pathVar.StoreID, err = s.retrieveIDParam(r, "store_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
		return
	}

	// parse request
	queryStr := r.URL.Query()
	var reqQueryStr listStoreItemsQueryStr

	reqQueryStr.ItemName = s.readStr(queryStr, "item_name", "")
	reqQueryStr.Sort = s.readStr(queryStr, "sort", "id")

	reqQueryStr.Page, _ = s.readInt(queryStr, "page", 1)
	reqQueryStr.PageSize, _ = s.readInt(queryStr, "page_size", 15)

	// validate query string
	if err := s.bindJSONWithValidation(w, r, &reqQueryStr, validator.New()); err != nil {
		return
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
	StoreID int64 `json:"store_id" validate:"required,min=1"`
	ItemID int64 `json:"item_id" validate:"required,min=1"`
}

// buyStoreItems maps to endpoint "PATCH /stores/:store_id/items/:item_id/buy"
func (s *StoreHub) buyStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar buyStoreItemsPathVar
	var err error

	// parse path variables
	pathVar.StoreID, err = s.retrieveIDParam(r, "store_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// parse path variables
	pathVar.ItemID, err = s.retrieveIDParam(r, "item_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
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
		arg := db.UpdateItemParams {
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
}

type getStoreItemsPathVar struct {
	StoreID int64 `json:"store_id" validate:"required,min=1"`
	ItemID int64 `json:"item_id" validate:"required,min=1"`
}

// maps to endpoint "GET /stores/:store_id/items/:item_id"
func (s *StoreHub) getStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar getStoreItemsPathVar
	var err error

	// parse path variables
	pathVar.StoreID, err = s.retrieveIDParam(r, "store_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// parse path variables
	pathVar.ItemID, err = s.retrieveIDParam(r, "item_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
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