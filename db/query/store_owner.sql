-- name: CreateStoreOwner :one
INSERT INTO store_owners (
  user_id,
  store_id
) VALUES (
  $1, $2
) RETURNING *;