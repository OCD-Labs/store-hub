-- name: CreateStoreOwner :one
INSERT INTO store_owners (
  user_id,
  store_id,
  access_level
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: DeleteStoreOwner :exec
DELETE FROM store_owners
WHERE user_id = $1 AND store_id = $2;

-- name: IsStoreOwner :one
SELECT COUNT(*) AS ownership_count
FROM store_owners
WHERE user_id = sqlc.arg(user_id)
  AND store_id = sqlc.arg(store_id);