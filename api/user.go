package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/OCD-Labs/store-hub/worker"
	"github.com/go-playground/validator"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
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
			"result": envelop{
				"stores": stores,
			},
		},
	}, nil)
}

type createUserRequest struct {
	FirstName       string  `json:"first_name" validate:"required,min=1"`
	LastName        string  `json:"last_name" validate:"required,min=1"`
	Password        string  `json:"password" validate:"required,min=8"`
	Email           string  `json:"email" validate:"required,email"`
	AccountID       string  `json:"account_id" validate:"required,len=2|len=64"`
	ProfileImageUrl *string `json:"profile_image_url"`
}

type userResponse struct {
	ID                int64     `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	AccountID         string    `json:"account_id"`
	Email             string    `json:"email"`
	ProfileImageurl   string    `json:"profil_image_url"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	IsActive          bool      `json:"is_active"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		AccountID:         user.AccountID,
		Email:             user.Email,
		ProfileImageurl:   user.ProfileImageUrl.String,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		IsActive:          user.IsActive,
		IsEmailVerified:   user.IsEmailVerified,
	}
}

// createUser maps to endpoint "POST /api/v1/users"
func (s *StoreHub) createUser(w http.ResponseWriter, r *http.Request) {
	// parse request
	var req createUserRequest
	if err := s.readJSON(w, r, &req); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, "failed to parse request body")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if err := s.bindJSONWithValidation(w, r, req, validator.New()); err != nil {
		return
	}

	// hash password
	hashedPassword, err := util.HashedPassword(req.Password)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to hash password")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db transaction
	createUserArg := db.CreateUserParams{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		AccountID:      req.AccountID,
		Status:         util.NORMALUSER,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	if req.ProfileImageUrl != nil {
		createUserArg.ProfileImageUrl.String = *req.ProfileImageUrl
		createUserArg.ProfileImageUrl.Valid = true
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: createUserArg,
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				UserID:    user.ID,
				ClientIp:  r.RemoteAddr,
				UserAgent: r.UserAgent(),
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			err := s.taskDistributor.DistributeTaskSendVerifyEmail(r.Context(), taskPayload, opts...)
			return err
		},
	}

	result, err := s.dbStore.CreateUserTx(r.Context(), arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				s.errorResponse(w, r, http.StatusForbidden, "user already exist")
				log.Error().Err(err).Msg("error occurred")
				return
			}
		}
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to create user")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "new user created",
			"result": envelop{
				"user": newUserResponse(result.User),
			},
		},
	}, nil)
}

type loginRequest struct {
	// Email    string `json:"email" validate:"required,email"`
	// Password string `json:"password" validate:"required,min=8"`
	AccountID string `json:"account_id" validate:"required,len=2|len=64"`
}

func (s *StoreHub) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := s.readJSON(w, r, &req); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, "failed to parse request body")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if err := s.bindJSONWithValidation(w, r, &req, validator.New()); err != nil {
		return
	}

	user, err := s.dbStore.GetUserByAccountID(r.Context(), req.AccountID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "user not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch user's profile")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// if !user.IsEmailVerified {
	// 	newReq, err := http.NewRequest(http.MethodPost, "/api/v1/resend_email_verification", r.Body)
	// 	if err != nil {
	// 		s.errorResponse(w, r, http.StatusInternalServerError, "failed to resend email verification mail")
	// 		log.Error().Err(err).Msg("error occurred")
	// 		return
	// 	}

	// 	for key, value := range r.Header {
	// 		newReq.Header.Set(key, value[0])
	// 	}

	// 	http.Redirect(w, r, fmt.Sprintf("/api/v1/resend_email_verification?user_id=%d", user.ID), http.StatusTemporaryRedirect)
	// }

	// if !user.IsActive {
	// 	s.errorResponse(w, r, http.StatusNoContent, "user is not activated")
	// 	return
	// }

	// err = util.CheckPassword(user.HashedPassword, req.Password)
	// if err != nil {
	// 	s.errorResponse(w, r, http.StatusUnauthorized, "Invalid login credentials")
	// 	log.Error().Err(err).Msg("error occurred")
	// 	return
	// }

	token, _, err := s.tokenMaker.CreateToken(user.ID, req.AccountID, 24*time.Hour)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to generate access token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "logged user in",
			"result": envelop{
				"user":         newUserResponse(user),
				"access_token": token,
			},
		},
	}, nil)
}

func (s *StoreHub) logout(w http.ResponseWriter, r *http.Request) {
	authPayload := s.contextGetToken(r)

	expiredAt := authPayload.ExpiredAt
	duration := time.Until(expiredAt)

	err := s.cache.BlacklistSession(r.Context(), authPayload.ID.String(), duration)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to blacklist access token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "logged out user",
		},
	}, nil)
}

type getUserPathVariable struct {
	ID int64 `json:"id" validate:"required,min=1"`
}

func (s *StoreHub) getUser(w http.ResponseWriter, r *http.Request) {
	authPayload := s.contextGetToken(r)
	var pathVar getUserPathVariable
	var err error

	pathVar.ID, err = s.retrieveIDParam(r, "id")
	if err != nil || pathVar.ID == 0 {
		s.errorResponse(w, r, http.StatusBadRequest, "invalid user id")
		return
	}

	if pathVar.ID != authPayload.UserID {
		s.errorResponse(w, r, http.StatusUnauthorized, "mismatched user")
		return
	}

	user, err := s.dbStore.GetUserByID(r.Context(), authPayload.UserID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "user not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch user's profile")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "logged out user",
			"result": envelop{
				"user": newUserResponse(user),
			},
		},
	}, nil)
}