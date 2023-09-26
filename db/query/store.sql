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
      'account_id', u.account_id,
      'profile_img_url', u.profile_image_url,
      'access_levels', so.access_levels,
      'is_original_owner', so.is_primary,
      'added_at', so.added_at
  )) AS store_owners
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