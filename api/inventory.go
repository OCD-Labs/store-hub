package api

import (
	"database/sql"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type addStoreItemRequestBody struct {
	Name               string   `json:"name" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Price              string   `json:"price" validate:"required"`
	ImageURLs          []string `json:"image_urls" validate:"required"`
	Category           string   `json:"category" validate:"category"`
	DiscountPercentage string   `json:"discount_percentage" validate:"required"`
	SupplyQuantity     int64    `json:"supply_quantity" validate:"required"`
}

type addStoreItemPathVar struct {
	StoreID int64 `json:"id" validate:"required,min=1"`
}

// discoverStoreByOwner maps to endpoint "POST /stores/{id}/items"
func (s *StoreHub) addStoreItem(w http.ResponseWriter, r *http.Request) {
	var pathVar addStoreItemPathVar
	var err error

	// parse path variables
	pathVar.StoreID, err = s.retrieveIDParam(r, "id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
		return
	}

	// parse request body
	var reqBody addStoreItemRequestBody
	if err := s.readJSON(w, r, &reqBody); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, "failed to parse request")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// validate request body
	if err := s.bindJSONWithValidation(w, r, &reqBody, validator.New()); err != nil {
		return
	}

	authPayload := s.contextGetToken(r) // authorize

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db query
	arg := db.CreateStoreItemParams{
		Name:               reqBody.Name,
		Description:        reqBody.Description,
		Price:              reqBody.Price,
		StoreID:            pathVar.StoreID,
		ImageUrls:          reqBody.ImageURLs,
		Category:           reqBody.Category,
		SupplyQuantity:     reqBody.SupplyQuantity,
		DiscountPercentage: reqBody.DiscountPercentage,
		Extra:              []byte("{}"),
	}
	item, err := s.dbStore.CreateStoreItem(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "add a new item",
			"result":  envelop{
				"item": item,
			},
		},
	}, nil)
}

type listOwnedStoreItemsQueryStr struct {
	ItemName string `json:"item_name"` // TODO: add category field
	Page     int    `json:"page" validate:"min=1,max=10000000"`
	PageSize int    `json:"page_size" validate:"min=1,max=20"`
	Sort     string `json:"sort"`
}

type listOwnedStoreItemsPathVar struct {
	StoreID int64 `json:"store_id" validate:"required,min=1"`
}

// listOwnedStoreItems maps to endpoint "GET /users/{user_id}/stores/{store_id}/items"
func (s *StoreHub) listOwnedStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar listOwnedStoreItemsPathVar
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
	var reqQueryStr listOwnedStoreItemsQueryStr

	reqQueryStr.ItemName = s.readStr(queryStr, "item_name", "")
	reqQueryStr.Sort = s.readStr(queryStr, "sort", "")

	reqQueryStr.Page, _ = s.readInt(queryStr, "page", 1)
	reqQueryStr.PageSize, _ = s.readInt(queryStr, "page_size", 15)

	// validate query string
	if err := s.bindJSONWithValidation(w, r, &reqQueryStr, validator.New()); err != nil {
		return
	}

	authPayload := s.contextGetToken(r) // authorize

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
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
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some store items",
			"result": envelop{
				"items":   items,
				"metadata": pagination,
			},
		},
	}, nil)

}

type updateStoreItemsRequestBody struct {
	Name               *string  `json:"name" validate:"min=2"`
	Description        *string  `json:"description" validate:"min=2"`
	Price              *string  `json:"price"`
	ImageURLs          []string `json:"image_urls"`
	Category           *string  `json:"category"`
	DiscountPercentage *string  `json:"discount_percentage"`
	SupplyQuantity     *int64   `json:"supply_quantity"`
}

type updateStoreItemsPathVar struct {
	StoreID int64 `json:"store_id" validate:"required,min=1"`
	ItemID int64 `json:"item_id" validate:"required,min=1"`
}

// updateStoreItems maps to endpoint "PATCH /stores/{store_id}/items/{item_id}/update"
func (s *StoreHub) updateStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar updateStoreItemsPathVar
	var err error

	// parse path variables
	pathVar.StoreID, err = s.retrieveIDParam(r, "store_id")
	if err != nil || pathVar.StoreID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	pathVar.ItemID, err = s.retrieveIDParam(r, "item_id")
	if err != nil || pathVar.ItemID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid item id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
		return
	}

	// parse request body
	var reqBody updateStoreItemsRequestBody
	if err := s.readJSON(w, r, &reqBody); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, "failed to parse request")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// validate request body
	if err := s.bindJSONWithValidation(w, r, &reqBody, validator.New()); err != nil {
		return
	}

	authPayload := s.contextGetToken(r) // authorize

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db query
	arg := db.UpdateItemParams{
		ItemID: pathVar.ItemID,
	}
	if reqBody.Name != nil {
		arg.Name = sql.NullString{
			String: *reqBody.Name,
			Valid: true,
		}
	}
	if reqBody.Description != nil {
		arg.Name = sql.NullString{
			String: *reqBody.Description,
			Valid: true,
		}
	}
	if reqBody.Price != nil {
		arg.Name = sql.NullString{
			String: *reqBody.Price,
			Valid: true,
		}
	}
	if reqBody.ImageURLs != nil {
		arg.ImageUrls = reqBody.ImageURLs
	}
	if reqBody.Category != nil {
		arg.Category = sql.NullString{
			String: *reqBody.Category,
			Valid: true,
		}
	}
	if reqBody.DiscountPercentage != nil {
		arg.DiscountPercentage = sql.NullString{
			String: *reqBody.DiscountPercentage,
			Valid: true,
		}
	}
	if reqBody.SupplyQuantity != nil {
		arg.SupplyQuantity = sql.NullInt64{
			Int64: *reqBody.SupplyQuantity,
			Valid: true,
		}
	}

	item, err := s.dbStore.UpdateItem(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to update item details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "updated item's details",
			"result": envelop{
				"item":   item,
			},
		},
	}, nil)
}

func (s *StoreHub) buyStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) freezeStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) unfreezeStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) deleteStoreItems(w http.ResponseWriter, r *http.Request) {

}