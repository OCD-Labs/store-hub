-- name: AddCoOwner :one
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
SELECT access_level
FROM store_owners
WHERE user_id = $1
  AND store_id = $2;

-- name: UpdateCoOwnerAccess :one
UPDATE store_owners 
SET 
  access_level = sqlc.arg(access_level)
WHERE 
  store_id = sqlc.arg(store_id) AND user_id = sqlc.arg(user_id)
RETURNING *;