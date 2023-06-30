package api

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type discoverStoreByOwnerPathVar struct {
	UserID int64 `json:"user_id" validate:"required,min=1"`
}

// listUserStores maps to endpoint "GET /users/{id}/stores"
func (s *StoreHub) listUserStores(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar discoverStoreByOwnerPathVar
	var err error

	// parse path variables
	pathVar.UserID, err = s.retrieveIDParam(r, "id")
	if err != nil || pathVar.UserID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid store id")
		return
	}

	// validate path variables
	if err := s.bindJSONWithValidation(w, r, &pathVar, validator.New()); err != nil {
		return
	}

	// authorise
	authPayload := s.contextGetToken(r)
	if authPayload.UserID != pathVar.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// db query
	stores, err := s.dbStore.GetStoreByOwner(r.Context(), authPayload.UserID)
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
			"result":  envelop{
				"stores": stores,
			},
		},
	}, nil)
}