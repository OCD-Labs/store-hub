package api

import (
	"database/sql"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type createStoreRequestBody struct {
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	ProfileImageUrl string `json:"profile_image_url" validate:"required"`
	Category        string `json:"category" validate:"required"`
}

// createStore maps to endpoint "POST /stores".
func (s *StoreHub) createStore(w http.ResponseWriter, r *http.Request) {
	var reqBody createStoreRequestBody
	if err := s.readJSON(w, r, &reqBody); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, "failed to parse request")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// verify request
	if err := s.bindJSONWithValidation(w, r, &reqBody, validator.New()); err != nil {
		return
	}

	// authorise
	authPayload := s.contextGetToken(r)

	// db query
	arg := db.CreateStoreTxParams{
		CreateStoreParams: db.CreateStoreParams{
			Name:            reqBody.Name,
			Description:     reqBody.Description,
			ProfileImageUrl: reqBody.ProfileImageUrl,
			Category:        reqBody.Category,
		},
		OwnerID:     authPayload.UserID,
		AccessLevel: 1,
	}
	_, err := s.dbStore.CreateStoreTx(r.Context(), arg)
	if err != nil { // TODO: Handle error due to Postgres constraints
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to create new store")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// check if user is a previous store owner
	if authPayload.UserRole != util.STOREOWNER {
		arg := db.UpdateUserParams{
			ID: sql.NullInt64{
				Int64: authPayload.UserID,
				Valid: true,
			},
			Status: sql.NullString{
				String: util.STOREOWNER,
				Valid:  true,
			},
		}
		_, err := s.dbStore.UpdateUser(r.Context(), arg)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to upgrade user to a store owner")
			log.Error().Err(err).Msg("error occurred")
			return
		}
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "created a new store",
		},
	}, nil)
}

type discoverStoreQueryStr struct {
	Page     int    `json:"page" validate:"min=1,max=10000000"`
	PageSize int    `json:"page_size" validate:"min=1,max=20"`
	Sort     string `json:"sort"`
}

func (s *StoreHub) discoverStore(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) freezeStore(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) unfreezeStore(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) deleteStore(w http.ResponseWriter, r *http.Request) {

}