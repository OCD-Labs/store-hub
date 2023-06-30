package api

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type discoverStoreByOwnerPathVar struct {
	UserID int `json:"user_id" validate:"required"`
}

// discoverStoreByOwner maps to endpoint "GET /users/{id}/stores"
func (s *StoreHub) discoverStoreByOwner(w http.ResponseWriter, r *http.Request) {
	// parse request
	queryStr := r.URL.Query()
	var reqQueryStr discoverStoreByOwnerPathVar

	reqQueryStr.UserID, _ = s.readInt(queryStr, "id", 0)

	// validate query string
	if err := s.bindJSONWithValidation(w, r, &reqQueryStr, validator.New()); err != nil {
		return
	}

	// authorise
	authPayload := s.contextGetToken(r)
	if int(authPayload.UserID) != reqQueryStr.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "access denied")
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
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "found your stores",
			"result":  envelop{
				"stores": stores,
			},
		},
	}, nil)
}