package api

import (
	"database/sql"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/rs/zerolog/log"
)

type createStoreRequestBody struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ProfileImageUrl string `json:"profile_image_url" validate:"required"`
	Category        string `json:"category" validate:"required"`
	StoreAccountID  string `json:"store_account_id" validate:"required,min=2,max=64"`
}

type createStorePathVar struct {
	UserID int64 `path:"user_id" validate:"required,min=1"`
}

// createStore maps to endpoint "POST /users/{id}/stores".
func (s *StoreHub) createStore(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar createStorePathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody createStoreRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// authorise
	authPayload := s.contextGetToken(r)
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

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
	}
	result, err := s.dbStore.CreateStoreTx(r.Context(), arg)
	if err != nil { // TODO: Handle error due to Postgres constraints
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to create new store")
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

type addStoreItemRequestBody struct {
	Name               string   `json:"name" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Price              string   `json:"price" validate:"required"`
	ImageURLs          []string `json:"image_urls" validate:"required"`
	Category           string   `json:"category" validate:"required"` // TODO: change DB schema to tags
	DiscountPercentage string   `json:"discount_percentage" validate:"required"`
	SupplyQuantity     int64    `json:"supply_quantity" validate:"required"`
	// TODO: Add currency DB schema
}

type addStoreItemPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// discoverStoreByOwner maps to endpoint "POST /users/{user_id}/stores/{store_id}/items"
func (s *StoreHub) addStoreItem(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar addStoreItemPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody addStoreItemRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
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
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// listOwnedStoreItems maps to endpoint "GET /users/{user_id}/stores/{store_id}/items"
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

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
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
	Category           *string  `json:"category"`
	DiscountPercentage *string  `json:"discount_percentage"`
	SupplyQuantity     *int64   `json:"supply_quantity"`
}

type updateStoreItemsPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	ItemID  int64 `path:"item_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// updateStoreItems maps to endpoint "PATCH /users/{user_id}/stores/{store_id}/items/{item_id}"
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

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
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

	item, err := s.dbStore.UpdateItem(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to update item details")
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
	UserID  int64 `path:"user_id" validate:"required"`
}

// deleteStoreItems maps to endpoint "DELETE /users/{user_id}/stores/{store_id}/items/{item_id}"
func (s *StoreHub) deleteStoreItems(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar deleteStoreItemsPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	err = s.dbStore.DeleteItem(r.Context(), db.DeleteItemParams{
		StoreID: pathVar.StoreID,
		ItemID:  pathVar.ItemID,
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete item")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusNoContent, envelop{ // TODO: remove response body as 204 status code means no body
		"status": "success",
		"data": envelop{
			"message": "deleted item and its details",
		},
	}, nil)
}

type addNewOwnerRequestBody struct {
	AccountID string `json:"account_id" validate:"required,min=2,max=64"`
}

type addNewOwnerPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// addNewOwner maps to endpoint "POST /users/{user_id}/store/{store_id}/owners"
func (s *StoreHub) addNewOwner(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar addNewOwnerPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody addNewOwnerRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if check.AccessLevel != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "higher access level needed for this action")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db query
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	arg := db.CreateStoreOwnerParams{
		AccessLevel: check.AccessLevel + 1,
		StoreID:     pathVar.StoreID,
		UserID:      user.ID,
	}
	newOwner, err := s.dbStore.CreateStoreOwner(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to add owner")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "added a new owner",
			"result": envelop{
				"owner": newOwner,
			},
		},
	}, nil)
}

type deleteOwnerPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

type deleteOwnerRequestBody struct {
	AccountID string `json:"account_id" validate:"required,min=2,max=64"`
}

// deleteOwner maps to endpoint "DELETE /users/{user_id}/store/{store_id}/owners"
func (s *StoreHub) deleteOwner(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar deleteOwnerPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody deleteOwnerRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if check.AccessLevel != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "higher access level needed for this action")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db query
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	err = s.dbStore.DeleteStoreOwner(r.Context(), db.DeleteStoreOwnerParams{
		UserID:  user.ID,
		StoreID: pathVar.StoreID,
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete owner")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusNoContent, envelop{
		"status": "success",
		"data": envelop{
			"message": "remove user from store ownership",
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

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check ownership
	check, err := s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID:  authPayload.UserID,
	})
	if check.OwnershipCount != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "access to store denied")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if check.AccessLevel != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "higher access level needed for this action")
		log.Error().Err(err).Msg("error occurred")
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
