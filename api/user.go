package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

type createUserRequest struct {
	FirstName       string  `json:"first_name" validate:"required,min=1"`
	LastName        string  `json:"last_name" validate:"required,min=1"`
	Password        string  `json:"password" validate:"required,min=8"`
	Email           string  `json:"email" validate:"required,email"`
	AccountID       string  `json:"account_id" validate:"required,min=2,max=64"`
	ProfileImageUrl *string `json:"profile_image_url"`
}

type userResponse struct {
	ID                int64           `json:"user_id"`
	FirstName         string          `json:"first_name"`
	LastName          string          `json:"last_name"`
	AccountID         string          `json:"account_id"`
	Status            string          `json:"status"`
	About             string          `json:"about"`
	Email             string          `json:"email"`
	Socials           json.RawMessage `json:"socials"`
	ProfileImageurl   string          `json:"profil_image_url"`
	CreatedAt         time.Time       `json:"created_at"`
	PasswordChangedAt time.Time       `json:"password_changed_at"`
	IsActive          bool            `json:"is_active"`
	IsEmailVerified   bool            `json:"is_email_verified"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		AccountID:         user.AccountID,
		Status:            user.Status,
		About:             user.About,
		Email:             user.Email,
		ProfileImageurl:   user.ProfileImageUrl.String,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		IsActive:          user.IsActive,
		IsEmailVerified:   user.IsEmailVerified,
	}
}

// createUser maps to endpoint "POST /users"
func (s *StoreHub) createUser(w http.ResponseWriter, r *http.Request) {
	var reqBody createUserRequest
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	// hash password
	hashedPassword, err := util.HashedPassword(reqBody.Password)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to hash password")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// db transaction
	createUserArg := db.CreateUserParams{
		FirstName:      reqBody.FirstName,
		LastName:       reqBody.LastName,
		AccountID:      reqBody.AccountID,
		Status:         util.NORMALUSER,
		HashedPassword: hashedPassword,
		Email:          reqBody.Email,
		Socials:        []byte("{}"),
	}

	if reqBody.ProfileImageUrl != nil {
		createUserArg.ProfileImageUrl.String = *reqBody.ProfileImageUrl
		createUserArg.ProfileImageUrl.Valid = true
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: createUserArg,
		AfterCreate: func(user db.User) (err error) {
			taskSendVerifyEmailPayload := &worker.PayloadSendVerifyEmail{
				UserID:    user.ID,
				ClientIp:  r.RemoteAddr,
				UserAgent: r.UserAgent(),
			}

			sendVerifyEmailopts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			err = s.taskDistributor.DistributeTaskSendVerifyEmail(r.Context(), taskSendVerifyEmailPayload, sendVerifyEmailopts...)
			if err != nil {
				return err
			}

			subaccount := fmt.Sprintf("%s.%s", util.SanitizeAccountID(reqBody.AccountID, user.ID), s.configs.NEARAccountID)

			taskNEARTxPayload := &worker.PayloadNEARTx{
				Args: []string{"create-account", subaccount, "--masterAccount", s.configs.NEARAccountID, "--initialBalance", "10"},
			}

			nearTxopts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			err = s.taskDistributor.DistributeTaskNEARTx(r.Context(), taskNEARTxPayload, nearTxopts...)

			return err
		},
	}

	user, err := s.dbStore.CreateUserTx(r.Context(), arg)
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
				"user": newUserResponse(user),
			},
		},
	}, nil)
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// login maps to endpoint "GET /auth/login"
func (s *StoreHub) login(w http.ResponseWriter, r *http.Request) {
	var reqBody loginRequest
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		log.Error().Err(err).Msg("error occurred")
		return
	}

	user, err := s.dbStore.GetUserByEmail(r.Context(), reqBody.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "Invalid login credentials")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch user's profile")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if !user.IsEmailVerified {
		taskSendVerifyEmailPayload := &worker.PayloadSendVerifyEmail{
			UserID:    user.ID,
			ClientIp:  r.RemoteAddr,
			UserAgent: r.UserAgent(),
		}

		sendVerifyEmailopts := []asynq.Option{
			asynq.MaxRetry(10),
			asynq.ProcessIn(5 * time.Second),
			asynq.Queue(worker.QueueCritical),
		}

		err = s.taskDistributor.DistributeTaskSendVerifyEmail(r.Context(), taskSendVerifyEmailPayload, sendVerifyEmailopts...)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to schedule a verification-email task")
			log.Error().Err(err).Msg("error occurred")
			return
		}

		s.writeJSON(w, http.StatusAccepted, envelop{
			"status":  "success",
			"message": "Please check your email to complete the verification process.",
		}, nil)
		return
	}

	if !user.IsActive {
		s.errorResponse(w, r, http.StatusNoContent, "user is not activated")
		return
	}

	err = util.CheckPassword(reqBody.Password, user.HashedPassword)
	if err != nil {
		s.errorResponse(w, r, http.StatusUnauthorized, "Invalid login credentials")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	token, _, err := s.tokenMaker.CreateToken(user.ID, user.AccountID, 24*time.Hour, nil)
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

// logout maps to endpoint "POST /auth/logout"
func (s *StoreHub) logout(w http.ResponseWriter, r *http.Request) {
	authPayload := s.contextGetMustToken(r)

	err := s.cache.BlacklistSession(r.Context(), authPayload.ID.String(), time.Until(authPayload.ExpiredAt))
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
	ID int64 `path:"user_id" validate:"required,min=1"`
}

// getUser maps to endpoint "GET /users/{id}"
func (s *StoreHub) getUser(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar getUserPathVariable
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	authPayload := s.contextGetMustToken(r)

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
			"message": "found user",
			"result": envelop{
				"user": newUserResponse(user),
			},
		},
	}, nil)
}

type verifyEmailQueryStr struct {
	Token string `querystr:"secret_code" validate:"required"`
	Email string `querystr:"email" validate:"required,email"`
}

// verifyEmail maps to endpoint "POST /users/verify-email"
func (s *StoreHub) verifyEmail(w http.ResponseWriter, r *http.Request) {
	var queryStr verifyEmailQueryStr
	if err := s.ShouldBindPathVars(w, r, &queryStr); err != nil {
		return
	}

	sessionExists, err := s.dbStore.CheckSessionExists(r.Context(), db.CheckSessionExistsParams{
		Token: queryStr.Token,
		Scope: "access_invitation_email",
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to verify token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if !sessionExists {
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
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to verify secret code")
		}

		log.Error().Err(err).Msg("error occurred")
		return
	}

	exists, err := s.cache.IsSessionBlacklisted(r.Context(), payload.ID.String())
	if err != nil || exists {
		s.errorResponse(w, r, http.StatusUnauthorized, "invalid token")
		return
	}

	arg := db.UpdateUserParams{
		IsEmailVerified: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Email: sql.NullString{
			String: queryStr.Email,
			Valid:  true,
		},
	}

	user, err := s.dbStore.UpdateUser(r.Context(), arg)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "Invalid login credentials")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to fetch user's profile")
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

	// access token
	token, _, err := s.tokenMaker.CreateToken(user.ID, user.AccountID, 24*time.Hour, nil)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to generate access token")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "user email verified",
			"result": envelop{
				"user":         newUserResponse(user),
				"access_token": token,
			},
		},
	}, nil)
}

// TODO: update user's profile
