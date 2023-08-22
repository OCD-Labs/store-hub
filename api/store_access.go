package api

import (
	"database/sql"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type grantStoreAccessRequestBody struct {
	AccountID      string `json:"account_id" validate:"required,min=2,max=64"`
	NewAccessLevel int32  `json:"new_access_level" validate:"required,min=1,max=5"`
}

type grantStoreAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// addNewOwner maps to endpoint "POST /users/{user_id}/store/{store_id}/access"
func (s *StoreHub) grantStoreAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar grantStoreAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody grantStoreAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check access is granted to an existing user.
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "can't grant access to non-existent user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	arg := db.AddCoOwnerAccessParams{
		AccessLevels: []int32{reqBody.NewAccessLevel},
		StoreID:      pathVar.StoreID,
		UserID:       user.ID,
		IsPrimary:    false,
	}
	newOwner, err := s.dbStore.AddCoOwnerAccess(r.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				s.errorResponse(w, r, http.StatusConflict, "incorrect create store owner details")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to add owner")
			}
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to add owner")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "added a new owner",
			"result": envelop{
				"store_owner": newOwner,
			},
		},
	}, nil)
}

type revokeAllUserAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

type revokeAllUserAccessRequestBody struct {
	AccountID string `json:"account_id" validate:"required,min=2,max=64"`
}

// deleteOwner maps to endpoint "DELETE /users/{user_id}/store/{store_id}/access"
func (s *StoreHub) revokeAllUserAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar revokeAllUserAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody revokeAllUserAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check access is granted to an existing user.
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "can't revoke access for a non-existent user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	err = s.dbStore.RevokeAllAccess(r.Context(), db.RevokeAllAccessParams{
		UserID:  user.ID,
		StoreID: pathVar.StoreID,
	})
	if err != nil { 
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "no access exists for user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete owner")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusNoContent, nil, nil)
}

type revokeUserAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

type revokeUserAccessRequestBody struct {
	AccountID   string `json:"account_id" validate:"required,min=2,max=64"`
	AccessLevel int16  `json:"access_level" validate:"required,min=1,max=5"`
}

// revokeUserAccess maps to endpoint "PATCH /users/{user_id}/store/{store_id}/revoke-access"
func (s *StoreHub) revokeUserAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar revokeUserAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody revokeUserAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check access is granted to an existing user.
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "can't revoke access for a non-existent user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	storeOwner, err := s.dbStore.RevokeAccess(r.Context(), db.RevokeAccessParams{
		UserID:              user.ID,
		StoreID:             pathVar.StoreID,
		AccessLevelToRevoke: reqBody.AccessLevel,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "no access exists for user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete owner")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "user access updated",
			"result": envelop{
				"store_owner": storeOwner,
			},
		},
	}, nil)
}

type addUserAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

type addUserAccessRequestBody struct {
	AccountID      string `json:"account_id" validate:"required,min=2,max=64"`
	NewAccessLevel int16  `json:"new_access_level" validate:"required,min=1,max=5"`
}

// revokeUserAccess maps to endpoint "PATCH /users/{user_id}/store/{store_id}/add-access"
func (s *StoreHub) addUserAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar addUserAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody addUserAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetToken(r) // authorize
	if pathVar.UserID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatch user")
		return
	}

	// check access is granted to an existing user.
	user, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "can't grant access to non-existent user")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	storeOwner, err := s.dbStore.AddToCoOwnerAccess(r.Context(), db.AddToCoOwnerAccessParams{
		UserID:         user.ID,
		StoreID:        pathVar.StoreID,
		NewAccessLevel: reqBody.NewAccessLevel,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "no access exists for user")
			return
		}
		
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete owner")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "user access updated",
			"result": envelop{
				"store_owner": storeOwner,
			},
		},
	}, nil)
}
