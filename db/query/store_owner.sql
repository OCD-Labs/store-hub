-- name: CreateStoreOwner :one
INSERT INTO store_owners (
  user_id,
  store_id,
  permission_level
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: DeleteStoreOwner :exec
DELETE FROM store_owners
WHERE user_id = $1 AND store_id = $2; 