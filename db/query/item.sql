-- name: CreateStoreItem :one
INSERT INTO items (
  name,
  description,
  price,
  store_id,
  image_urls,
  category,
  discount_percentage,
  supply_quantity,
  extra
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = sqlc.arg(item_id) AND supply_quantity > 0;

-- name: UpdateItem :one
UPDATE items 
SET 
  description = COALESCE(sqlc.narg(description), description),
  price = COALESCE(sqlc.narg(price), price),
  image_urls = COALESCE(sqlc.narg(image_urls), image_urls),
  category = COALESCE(sqlc.narg(category), category),
  discount_percentage = COALESCE(sqlc.narg(discount_percentage), discount_percentage),
  supply_quantity = COALESCE(sqlc.narg(supply_quantity), supply_quantity),
  extra = COALESCE(sqlc.narg(extra), extra),
  is_frozen = COALESCE(sqlc.narg(is_frozen), is_frozen),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE
  id = sqlc.arg(item_id)
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = sqlc.arg(item_id);