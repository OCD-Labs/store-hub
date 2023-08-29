-- name: CreateReview :exec
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
);

-- name: ListReviews :many
SELECT 
  r.*,
  u.first_name,
  u.last_name,
  u.account_id,
  u.profile_image_url
FROM 
  reviews r
JOIN 
  users u ON r.user_id = u.id
WHERE 
  r.store_id = sqlc.arg(store_id) 
  OR r.item_id = sqlc.arg(item_id);

-- name: UpdateUserReview :one
UPDATE reviews
SET 
  rating = COALESCE(sqlc.narg(rating), rating),
  comment = COALESCE(sqlc.narg(comment), comment)
WHERE 
  id = sqlc.arg(review_id) 
  AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: DeleteReview :exec
DELETE FROM reviews
WHERE 
  id = sqlc.arg(review_id) 
  AND user_id = sqlc.arg(user_id);

-- name: ratingOverview :one
WITH RatingCounts AS (
    SELECT rating, COUNT(rating) as rate_count
    FROM reviews
    WHERE store_id = sqlc.arg(store_id)
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
WHERE r.store_id = sqlc.arg(store_id);
