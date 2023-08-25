-- name: AddCoOwnerAccess :one
INSERT INTO store_owners (
  store_id,
  user_id,
  access_levels,
  is_primary
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: RevokeAllAccess :exec
DELETE FROM store_owners
WHERE user_id = $1 AND store_id = $2;

-- name: RevokeAccess :exec
UPDATE store_owners 
SET access_levels = ARRAY_REMOVE(access_levels, sqlc.arg(access_level_to_revoke)::int)
WHERE 
  user_id = sqlc.arg(user_id) AND store_id = sqlc.arg(store_id)
;

-- name: GetUserAccessLevelsForStore :one
SELECT access_levels
FROM store_owners
WHERE user_id = sqlc.arg(user_id)
  AND store_id = sqlc.arg(store_id);

-- name: AddToCoOwnerAccess :one
UPDATE store_owners 
SET 
  access_levels = array_append(access_levels, sqlc.arg(new_access_level))
WHERE 
  store_id = sqlc.arg(store_id) AND user_id = sqlc.arg(user_id)
RETURNING *;

-- name: GetStoreOwnersByStoreID :many
SELECT *
FROM store_owners
WHERE store_id = sqlc.arg(store_id);
