// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: review.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createReview = `-- name: CreateReview :exec
INSERT INTO reviews (
  store_id,
  user_id,
  item_id,
  rating,
  review_type,
  comment,
  is_verified_purchase
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
`

type CreateReviewParams struct {
	StoreID            int64  `json:"store_id"`
	UserID             int64  `json:"user_id"`
	ItemID             int64  `json:"item_id"`
	Rating             string `json:"rating"`
	ReviewType         string `json:"review_type"`
	Comment            string `json:"comment"`
	IsVerifiedPurchase bool   `json:"is_verified_purchase"`
}

func (q *Queries) CreateReview(ctx context.Context, arg CreateReviewParams) error {
	_, err := q.db.ExecContext(ctx, createReview,
		arg.StoreID,
		arg.UserID,
		arg.ItemID,
		arg.Rating,
		arg.ReviewType,
		arg.Comment,
		arg.IsVerifiedPurchase,
	)
	return err
}

const deleteReview = `-- name: DeleteReview :exec
DELETE FROM reviews
WHERE 
  id = $1 
  AND user_id = $2
`

type DeleteReviewParams struct {
	ReviewID int64 `json:"review_id"`
	UserID   int64 `json:"user_id"`
}

func (q *Queries) DeleteReview(ctx context.Context, arg DeleteReviewParams) error {
	_, err := q.db.ExecContext(ctx, deleteReview, arg.ReviewID, arg.UserID)
	return err
}

const listReviews = `-- name: ListReviews :many
SELECT 
  r.id, r.store_id, r.user_id, r.item_id, r.rating, r.review_type, r.comment, r.is_verified_purchase, r.created_at, r.updated_at,
  u.first_name,
  u.last_name,
  u.account_id,
  u.profile_image_url
FROM 
  reviews r
JOIN 
  users u ON r.user_id = u.id
WHERE 
  r.store_id = $1 
  OR r.item_id = $2
`

type ListReviewsParams struct {
	StoreID int64 `json:"store_id"`
	ItemID  int64 `json:"item_id"`
}

type ListReviewsRow struct {
	ID                 int64          `json:"id"`
	StoreID            int64          `json:"store_id"`
	UserID             int64          `json:"user_id"`
	ItemID             int64          `json:"item_id"`
	Rating             string         `json:"rating"`
	ReviewType         string         `json:"review_type"`
	Comment            string         `json:"comment"`
	IsVerifiedPurchase bool           `json:"is_verified_purchase"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
	AccountID          string         `json:"account_id"`
	ProfileImageUrl    sql.NullString `json:"profile_image_url"`
}

func (q *Queries) ListReviews(ctx context.Context, arg ListReviewsParams) ([]ListReviewsRow, error) {
	rows, err := q.db.QueryContext(ctx, listReviews, arg.StoreID, arg.ItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListReviewsRow{}
	for rows.Next() {
		var i ListReviewsRow
		if err := rows.Scan(
			&i.ID,
			&i.StoreID,
			&i.UserID,
			&i.ItemID,
			&i.Rating,
			&i.ReviewType,
			&i.Comment,
			&i.IsVerifiedPurchase,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FirstName,
			&i.LastName,
			&i.AccountID,
			&i.ProfileImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserReview = `-- name: UpdateUserReview :one
UPDATE reviews
SET 
  rating = COALESCE($1, rating),
  comment = COALESCE($2, comment)
WHERE 
  id = $3 
  AND user_id = $4
RETURNING id, store_id, user_id, item_id, rating, review_type, comment, is_verified_purchase, created_at, updated_at
`

type UpdateUserReviewParams struct {
	Rating   sql.NullString `json:"rating"`
	Comment  sql.NullString `json:"comment"`
	ReviewID int64          `json:"review_id"`
	UserID   int64          `json:"user_id"`
}

func (q *Queries) UpdateUserReview(ctx context.Context, arg UpdateUserReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, updateUserReview,
		arg.Rating,
		arg.Comment,
		arg.ReviewID,
		arg.UserID,
	)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.UserID,
		&i.ItemID,
		&i.Rating,
		&i.ReviewType,
		&i.Comment,
		&i.IsVerifiedPurchase,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const ratingOverview = `-- name: ratingOverview :one
WITH RatingCounts AS (
    SELECT rating, COUNT(rating) as rate_count
    FROM reviews
    WHERE store_id = $1
    GROUP BY rating
)

SELECT 
    COUNT(*) as total_reviews,
    SUM(CASE WHEN DATE(r.created_at) = CURRENT_DATE THEN 1 ELSE 0 END) as total_reviews_today,
    ROUND(AVG(r.rating), 2) as average_rating,
    COALESCE((SELECT rate_count FROM RatingCounts WHERE rating = 1), 0) as rate_1_count,
    COALESCE((SELECT rate_count FROM RatingCounts WHERE rating = 2), 0) as rate_2_count,
    COALESCE((SELECT rate_count FROM RatingCounts WHERE rating = 3), 0) as rate_3_count,
    COALESCE((SELECT rate_count FROM RatingCounts WHERE rating = 4), 0) as rate_4_count,
    COALESCE((SELECT rate_count FROM RatingCounts WHERE rating = 5), 0) as rate_5_count
FROM reviews r
WHERE r.store_id = $1
`

type ratingOverviewRow struct {
	TotalReviews      int64       `json:"total_reviews"`
	TotalReviewsToday int64       `json:"total_reviews_today"`
	AverageRating     string      `json:"average_rating"`
	Rate1Count        interface{} `json:"rate_1_count"`
	Rate2Count        interface{} `json:"rate_2_count"`
	Rate3Count        interface{} `json:"rate_3_count"`
	Rate4Count        interface{} `json:"rate_4_count"`
	Rate5Count        interface{} `json:"rate_5_count"`
}

func (q *Queries) ratingOverview(ctx context.Context, storeID int64) (ratingOverviewRow, error) {
	row := q.db.QueryRowContext(ctx, ratingOverview, storeID)
	var i ratingOverviewRow
	err := row.Scan(
		&i.TotalReviews,
		&i.TotalReviewsToday,
		&i.AverageRating,
		&i.Rate1Count,
		&i.Rate2Count,
		&i.Rate3Count,
		&i.Rate4Count,
		&i.Rate5Count,
	)
	return i, err
}