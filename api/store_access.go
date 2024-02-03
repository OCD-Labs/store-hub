package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/token"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// TODO: Revoke an access invitation email.
// TODO: Preventing send the same access invitation multiple times.

type grantStoreAccessQueryStr struct {
	Token string `querystr:"sth_code"`
}

// addNewOwner maps to endpoint "POST /inventory/stores/{store_id}/accept-access-invitation"
func (s *StoreHub) grantStoreAccess(w http.ResponseWriter, r *http.Request) {
	var queryStr grantStoreAccessQueryStr
	if err := s.shouldBindQuery(w, r, &queryStr); err != nil {
		return
	}

	exists, err := s.dbStore.CheckSessionExists(r.Context(), db.CheckSessionExistsParams{
		Token: queryStr.Token,
		Scope: "access_invitation_email",
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to verify token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if !exists {
		s.errorResponse(w, r, http.StatusBadRequest, token.ErrInvalidToken.Error())
		return
	}

	payload, err := s.tokenMaker.VerifyToken(util.Concat(queryStr.Token))
	if err != nil {
		switch {
		case errors.Is(err, token.ErrExpiredToken):
			s.errorResponse(w, r, http.StatusBadRequest, token.ErrExpiredToken.Error())
		case errors.Is(err, token.ErrInvalidToken):
			s.errorResponse(w, r, http.StatusBadRequest, token.ErrInvalidToken.Error())
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to verify token")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	if exists, err := s.cache.IsSessionBlacklisted(r.Context(), payload.ID.String()); err != nil || exists {
		s.errorResponse(w, r, http.StatusUnauthorized, "token has been used before")
		return
	}

	var tokenExtra worker.TokenExtra
	err = token.ExtractExtra(payload, &tokenExtra)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to verify token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	coOwnerAccess, err := s.dbStore.AddCoOwnerAccessTx(r.Context(), db.AddCoOwnerAccessTxParams{
		StoreID:     tokenExtra.StoreID,
		InviteeID:   tokenExtra.InviteeID,
		AccessLevel: tokenExtra.AccessLevel,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
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

	err = s.cache.BlacklistSession(r.Context(), payload.ID.String(), time.Until(payload.ExpiredAt))
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to blacklist access token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "added a new access to store",
			"result": envelop{
				"co_owner_access": coOwnerAccess,
			},
		},
	}, nil)
}

type revokeAllUserAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

type revokeAllUserAccessRequestBody struct {
	AccountID string `json:"account_id" validate:"required,min=2,max=64"`
}

// deleteOwner maps to endpoint "DELETE /inventory/stores/:store_id/revoke-all-access"
func (s *StoreHub) revokeAllUserAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar revokeAllUserAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody revokeAllUserAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	storeOwners, err := s.dbStore.RevokeAccessTx(r.Context(), db.RevokeAccessTxParams{
		AccountID: reqBody.AccountID,
		StoreID:   pathVar.StoreID,
		DeleteAll: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "failed to retrieve user details or no access exists for user")
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete all access")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message":      "user access updated",
			"store_owners": storeOwners,
		},
	}, nil)
}

type revokeUserAccessPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

type revokeUserAccessRequestBody struct {
	AccountID   string `json:"account_id" validate:"required,min=2,max=64"`
	AccessLevel int32  `json:"access_level" validate:"required,min=1,max=5"`
}

// revokeUserAccess maps to endpoint "PATCH /inventory/stores/:store_id/revoke-access"
func (s *StoreHub) revokeUserAccess(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar revokeUserAccessPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody revokeUserAccessRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	storeOwners, err := s.dbStore.RevokeAccessTx(r.Context(), db.RevokeAccessTxParams{
		AccountID:           reqBody.AccountID,
		StoreID:             pathVar.StoreID,
		AccessLevelToRevoke: reqBody.AccessLevel,
		DeleteAll:           false,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "failed to retrieve user details or no access exists for user")
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to delete all access")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message":      "user access updated",
			"store_owners": storeOwners,
		},
	}, nil)
}

type sendAccessInvitationRequestBody struct {
	AccountID      string `json:"account_id" validate:"required,min=2,max=64"`
	NewAccessLevel int32  `json:"new_access_level" validate:"required,min=1,max=5"`
}

type sendAccessInvitationPathVar struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

// sendAccessInvitation maps to endpoint "POST /inventory/stores/{store_id}/send-access-invitation"
func (s *StoreHub) sendAccessInvitation(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar sendAccessInvitationPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody sendAccessInvitationRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	// check access is granted to an existing user.
	invitee, err := s.dbStore.GetUserByAccountID(r.Context(), reqBody.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "can't grant access to non-existent user")
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve user details")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	authPayload := s.contextGetMustToken(r) // authorize

	taskPayload := &worker.PayloadSendAccessInvitation{
		InviterID:        authPayload.UserID,
		InviteeAccountID: invitee.AccountID,
		InviteeID:        invitee.ID,
		InviteeEmail:     invitee.Email,
		AccessLevel:      reqBody.NewAccessLevel,
		StoreID:          pathVar.StoreID,
		ClientIp:         r.RemoteAddr,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = s.taskDistributor.DistributeTaskSendAccessInvitation(r.Context(), taskPayload, opts...)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to send invitation email")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "Access invitation sent",
		},
	}, nil)
}
