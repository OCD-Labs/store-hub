package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// TODO: Create index on search columns like category, tag price etc.
// TODO: When listing storefront items for a store, discuss catering for out of stock and delete items.
// TODO: push the boundary of th NUMERIC db type. Find out the limit of value will break db column with NUMERIC type.

type createStoreRequestBody struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ProfileImageUrl string `json:"profile_image_url" validate:"required"`
	Category        string `json:"category" validate:"required"`
	StoreAccountID  string `json:"store_account_id" validate:"required,min=2,max=64"`
}

// createStore maps to endpoint "POST /inventory/stores".
func (s *StoreHub) createStore(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var reqBody createStoreRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetMustToken(r)

	// db query
	arg := db.CreateStoreTxParams{
		CreateStoreParams: db.CreateStoreParams{
			Name:            reqBody.Name,
			Description:     reqBody.Description,
			ProfileImageUrl: reqBody.ProfileImageUrl,
			Category:        reqBody.Category,
			StoreAccountID:  reqBody.StoreAccountID,
		},
		OwnerID: authPayload.UserID,
		AfterCreate: func(ctx context.Context, store db.Store) (err error) {
			subaccount := fmt.Sprintf("%s-%d.%s", util.SanitizeAccountID(reqBody.StoreAccountID, s.configs.NEARNetwork), store.ID, s.configs.NEARAccountID)
			taskNEARTxPayload := &worker.PayloadNEARTx{
				Args: []string{"create-account", subaccount, "--masterAccount", s.configs.NEARAccountID, "--initialBalance", "1"},
			}

			nearTxOpts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return s.taskDistributor.DistributeTaskNEARTx(ctx, taskNEARTxPayload, nearTxOpts...)
		},
	}
	result, err := s.dbStore.CreateStoreTx(r.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				s.errorResponse(w, r, http.StatusConflict, "A store with the same AccountID already exists.")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to create new store")
			}
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to create new store")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "created a new store",
			"result":  result,
		},
	}, nil)
}

// listUserStores maps to endpoint "GET /inventory/stores"
func (s *StoreHub) listUserStores(w http.ResponseWriter, r *http.Request) {
	// authorise
	authPayload := s.contextGetMustToken(r)

	// db query
	stores, err := s.dbStore.ListUserStoresWithAccess(r.Context(), authPayload.UserID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve stores")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found your stores",
			"result": envelop{
				"stores": stores,
			},
		},
	}, nil)
}

type addStoreItemRequestBody struct {
	Name               string   `json:"name" validate:"required"`
	Description        string   `json:"description" validate:"required"` // TODO: Check the DB schema for the NUMERIC type if it's enough to accommodate big price
	Price              string   `json:"price" validate:"required"`       // TODO: validate the value contain in the string is valid number
	ImageURLs          []string `json:"image_urls" validate:"required"`
	Category           string   `json:"category" validate:"required"` // TODO: change DB schema to tags
	DiscountPercentage string   `json:"discount_percentage" validate:"required"`
	SupplyQuantity     int64    `json:"supply_quantity" validate:"required"`
	CoverImgURL        string   `json:"cover_img_url" validate:"required"`
	Status             string   `json:"status" validate:"required,oneof=VISIBLE HIDDEN"`
}

type addStoreItemPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

// addStoreItem maps to endpoint "POST /inventory/stores/{store_id}/items"
func (s *StoreHub) addStoreItem(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar addStoreItemPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody addStoreItemRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
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
		CoverImgUrl:        reqBody.CoverImgURL,
		Extra:              []byte("{}"),
		Status:             reqBody.Status,
	}
	item, err := s.dbStore.CreateStoreItem(r.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				s.errorResponse(w, r, http.StatusConflict, "Referenced store doesn't exist.")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item")
			}
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "add a new item",
			"result": envelop{
				"item": item,
			},
		},
	}, nil)
}

type listOwnedStoreItemsQueryStr struct {
	ItemName string `querystr:"item_name"` // TODO: add tags field
	Page     int    `querystr:"page" validate:"max=10000000"`
	PageSize int    `querystr:"page_size" validate:"max=20"`
	Sort     string `querystr:"sort"`
}

type listOwnedStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

// listOwnedStoreItems maps to endpoint "GET /inventory/stores/{store_id}/items"
func (s *StoreHub) listOwnedStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar listOwnedStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse query string
	var reqQueryStr listOwnedStoreItemsQueryStr
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

type updateStoreItemsRequestBody struct { // TODO: write custom validation tags for Name and Description fields
	Name               *string  `json:"name"`
	Description        *string  `json:"description"`
	Price              *string  `json:"price"`
	ImageURLs          []string `json:"image_urls"`
	CoverImgURL        *string  `json:"cover_img_url"`
	Category           *string  `json:"category"`
	DiscountPercentage *string  `json:"discount_percentage"`
	SupplyQuantity     *int64   `json:"supply_quantity"`
	Status             *string  `json:"status"`
}

type updateStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	ItemID  int64 `path:"item_id" validate:"required,min=1"`
}

// updateStoreItems maps to endpoint "PATCH /inventory/stores/{store_id}/items/{item_id}"
func (s *StoreHub) updateStoreItems(w http.ResponseWriter, r *http.Request) {
	var pathVar updateStoreItemsPathVar

	// parse path variables
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody updateStoreItemsRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
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
			Valid:  true,
		}
	}
	if reqBody.Description != nil {
		arg.Description = sql.NullString{
			String: *reqBody.Description,
			Valid:  true,
		}
	}
	if reqBody.Price != nil {
		arg.Price = sql.NullString{
			String: *reqBody.Price,
			Valid:  true,
		}
	}
	if reqBody.ImageURLs != nil {
		arg.ImageUrls = reqBody.ImageURLs
	}
	if reqBody.CoverImgURL != nil {
		arg.CoverImgUrl = sql.NullString{
			String: *reqBody.CoverImgURL,
			Valid:  true,
		}
	}
	if reqBody.Category != nil {
		arg.Category = sql.NullString{
			String: *reqBody.Category,
			Valid:  true,
		}
	}
	if reqBody.DiscountPercentage != nil {
		arg.DiscountPercentage = sql.NullString{
			String: *reqBody.DiscountPercentage,
			Valid:  true,
		}
	}
	if reqBody.SupplyQuantity != nil {
		arg.SupplyQuantity = sql.NullInt64{
			Int64: *reqBody.SupplyQuantity,
			Valid: true,
		}
	}
	if reqBody.Status != nil && (*reqBody.Status == "HIDDEN" || *reqBody.Status == "VISIBLE" ){
		arg.Status = sql.NullString{
			String: *reqBody.Status,
			Valid:  true,
		}
	}

	item, err := s.dbStore.UpdateItem(r.Context(), arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "item not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to update item details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "updated item's details",
			"result": envelop{
				"item": item,
			},
		},
	}, nil)
}

type deleteStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required"`
	ItemID  int64 `path:"item_id" validate:"required"`
}

// deleteStoreItems maps to endpoint "DELETE /inventory/stores/{store_id}/items/{item_id}"
func (s *StoreHub) deleteStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar deleteStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	err := s.dbStore.DeleteItem(r.Context(), db.DeleteItemParams{
		StoreID: pathVar.StoreID,
		ItemID:  pathVar.ItemID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "item not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete item")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusNoContent, nil, nil)
}

type updateStoreProfilePathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

type updateStoreProfileRquestBody struct {
	Name            *string  `json:"name"`
	Description     *string  `json:"description"`
	ProfileImageUrl *string  `json:"profile_image_url"`
	Category        *string  `json:"category"`
	Tags            []string `json:"tags"` // TODO: Ask if updating account_id of a store is necessary
}

// updateStoreProfile maps to "PATCH /api/v1/users/:user_id/stores/:store_id"
func (s *StoreHub) updateStoreProfile(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar updateStoreProfilePathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody updateStoreProfileRquestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	arg := db.UpdateStoreParams{
		StoreID: pathVar.StoreID,
	}

	if reqBody.Name != nil {
		arg.Name = sql.NullString{
			String: *reqBody.Name,
			Valid:  true,
		}
	}
	if reqBody.Description != nil {
		arg.Description = sql.NullString{
			String: *reqBody.Description,
			Valid:  true,
		}
	}
	if reqBody.ProfileImageUrl != nil {
		arg.ProfileImageUrl = sql.NullString{
			String: *reqBody.ProfileImageUrl,
			Valid:  true,
		}
	}
	if reqBody.Category != nil {
		arg.Category = sql.NullString{
			String: *reqBody.Category,
			Valid:  true,
		}
	}
	// TODO: Add if section for Tags
	// if reqBody.Tags != nil {
	// 	arg.Tags = sql.NullString{
	// 		String: *reqBody.Tags,
	// 		Valid: true,
	// 	}
	// }

	updatedStore, err := s.dbStore.UpdateStore(r.Context(), arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "Store not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to update store profile")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "updated store profile",
			"result": envelop{
				"store": updatedStore,
			},
		},
	}, nil)
}

type deleteStorePathVar struct {
	StoreID int64 `json:"store_id" validate:"required,min=1"`
	UserID  int64 `json:"user_id" validate:"required,min=1"`
}

// deleteStore maps to endpoint "DELETE /users/{user_id}/stores/{store_id}"
func (s *StoreHub) deleteStore(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar deleteStorePathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	authPayload := s.contextGetMustToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// TODO:
	// 	1. Delete all its items
	// 	2. Delete all its owners' records
	// 	3. then delete the store
}

func (s *StoreHub) freezeStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) unfreezeStoreItems(w http.ResponseWriter, r *http.Request) {

}
