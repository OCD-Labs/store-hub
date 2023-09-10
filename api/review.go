package api

import (
	"errors"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
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
	OrderID int64 `path:"order_id" validatw:"required,min=1"`
	UserID  int64 `path:"user_id" validate:"required,min=1"`
}

// addReview maps to endpoint "PUT /users/{user_id}/reviews/{order_id}"
func (s *StoreHub) addReview(w http.ResponseWriter, r *http.Request) {
	var reqBody addReviewRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	var pathVars addReviewPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	authPayload := s.contextGetMustToken(r)

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

type listItemReviewStorefrontQueryStr struct {
	Page     int    `querystr:"page" validate:"max=10000000"`
	PageSize int    `querystr:"page_size" validate:"max=20"`
	Sort     string `querystr:"sort"`
}

type listItemReviewStorefrontPathVars struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
	ItemID  int64 `path:"item_id" validatw:"required,min=1"`
}

// listItemReviewStorefront map to endpoint "GET /stores/{store_id}/items/{item_id}/reviews"
func (s *StoreHub) listItemReviewStorefront(w http.ResponseWriter, r *http.Request) {
	var reqQueryStr listItemReviewStorefrontQueryStr
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

	var pathVars listItemReviewStorefrontPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	reviews, pagination, err := s.dbStore.ListReviews(r.Context(), db.ListReviewsParams{
		StoreID:      pathVars.StoreID,
		ItemID:       pathVars.ItemID,
		IsStorefront: true,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"-id", "id"},
		},
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve reviews for item")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some reviews",
			"result": envelop{
				"reviews":  reviews,
				"metadata": pagination,
			},
		},
	}, nil)
}

type listItemReviewInventoryQueryStr struct {
	Page     int    `querystr:"page" validate:"max=10000000"`
	PageSize int    `querystr:"page_size" validate:"max=20"`
	Sort     string `querystr:"sort"`
}

type listItemReviewInventoryPathVars struct {
	StoreID int64 `path:"store_id" validate:"required,min=1"`
}

// listItemReviewStorefront map to endpoint "GET /inventory/stores/{store_id}/reviews"
func (s *StoreHub) listItemReviewInventory(w http.ResponseWriter, r *http.Request) {
	var reqQueryStr listItemReviewInventoryQueryStr
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

	var pathVars listItemReviewInventoryPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	reviews, pagination, err := s.dbStore.ListReviews(r.Context(), db.ListReviewsParams{
		StoreID: pathVars.StoreID,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"-id", "id"},
		},
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve reviews for item")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	aggregatedReviews, err := s.dbStore.RatingOverview(r.Context(), pathVars.StoreID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve reviews for item")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some reviews",
			"result": envelop{
				"reviews":            reviews,
				"aggregated_reviews": aggregatedReviews,
				"metadata":           pagination,
			},
		},
	}, nil)
}
