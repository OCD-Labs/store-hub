package api

import (
	"errors"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type addReviewRequestBody struct {
	StoreID int64  `json:"store_id" validate:"required,min=1"`
	ItemID  int64  `json:"item_id" validatw:"required,min=1"`
	Rating  string `json:"rating" validate:"required,oneof=1 2 3 4 5"`
	Comment string `json:"comment" validate:"required"`
}

type addReviewPathVars struct {
	OrderID   int64  `path:"order_id" validatw:"required,min=1"`
	AccountID string `path:"account_id" validate:"required"`
}

// addReview maps to endpoint "PUT /accounts/{account_id}/reviews/{order_id}"
func (s *StoreHub) addReview(w http.ResponseWriter, r *http.Request) {
	var reqBody addReviewRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	var pathVars addReviewPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	authPayload := s.contextGetToken(r)

	arg := db.CreateReviewTxParams{
		StoreID: reqBody.StoreID,
		UserID:  authPayload.UserID,
		ItemID:  reqBody.ItemID,
		Rating:  reqBody.Rating,
		Comment: reqBody.Comment,
		OrderID: pathVars.OrderID,
	}

	err := s.dbStore.CreateReviewTx(r.Context(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				s.errorResponse(w, r, http.StatusConflict, "incorrect create review params")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to create item review")
			}
		} else if errors.Is(err, db.ErrNoPurchase) {
			s.errorResponse(w, r, http.StatusForbidden, err.Error())
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
			"message": "added a new review to store",
		},
	}, nil)
}
