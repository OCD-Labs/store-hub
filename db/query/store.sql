-- name: CreateStore :one
INSERT INTO stores (
  name,
  description,
  profile_image_url,
  store_account_id,
  category
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetStoreByID :one
SELECT 
  s.*, 
  json_agg(json_build_object(
      'user', json_build_object('id', u.id, 'account_id', u.account_id, 'first_name', u.first_name, 'last_name', u.last_name, 'email', u.email),
      'store_owners', json_build_object('user_id', so.user_id, 'store_id', so.store_id, 'added_at', so.added_at)
  )) AS owners
FROM 
  stores AS s
JOIN 
  store_owners AS so ON s.id = so.store_id
JOIN 
  users AS u ON so.user_id = u.id
WHERE 
  s.id = sqlc.arg(store_id)
GROUP BY 
  s.id;

-- name: GetStoreByOwner :many
SELECT s.*
FROM stores s
JOIN store_owners so ON s.id = so.store_id
WHERE so.user_id = sqlc.arg(user_id);


-- name: UpdateStore :one
UPDATE stores
SET
  name = COALESCE(sqlc.narg(name), name),
  description = COALESCE(sqlc.narg(description), description),
  profile_image_url = COALESCE(sqlc.narg(profile_image_url), profile_image_url),
  is_verified = COALESCE(sqlc.narg(is_verified), is_verified),
  category = COALESCE(sqlc.narg(category), category),
  is_frozen = COALESCE(sqlc.narg(is_frozen), is_frozen)
WHERE 
  id = sqlc.arg(store_id)
RETURNING *;

-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = sqlc.arg(store_id);