package api

import (
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type addStoreItemRequestBody struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price string `json:"price" validate:"required"`
	ImageURLs []string `json:"image_urls" validate:"required"`
	Category string `json:"category" validate:"category"`
	DiscountPercentage string `json:"discount_percentage" validate:"required"`
	SupplyQuantity int64 `json:"supply_quantity" validate:"required"`
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
	check, err := s.dbStore.IsStoreOwner(r.Context(),db.IsStoreOwnerParams{
		StoreID: pathVar.StoreID,
		UserID: authPayload.UserID,
	})
	if check != 1 {
		s.errorResponse(w, r, http.StatusForbidden, "no access to store")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db query
	arg := db.CreateStoreItemParams{
		Name: reqBody.Name,
		Description: reqBody.Description,
		Price: reqBody.Price,
		StoreID: pathVar.StoreID,
		ImageUrls: reqBody.ImageURLs,
		Category: reqBody.Category,
		SupplyQuantity: reqBody.SupplyQuantity,
		DiscountPercentage: reqBody.DiscountPercentage,
		Extra: []byte("{}"),
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
			"result": item,
		},
	}, nil)
}

func (s *StoreHub) listStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) updateStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) buyStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) freezeStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) unfreezeStoreItems(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) deleteStoreItems(w http.ResponseWriter, r *http.Request) {

}