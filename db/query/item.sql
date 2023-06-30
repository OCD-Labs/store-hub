-- name: CreateStoreItem :one
INSERT INTO items (
  description,
  price,
  store_id,
  image_urls,
  category,
  discount_percentage,
  supply_quantity,
  extra
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = sqlc.arg(item_id);

-- name: UpdateItem :one
UPDATE items 
SET 
  description = COALESCE(sqlc.narg(description), description),
  price = COALESCE(sqlc.narg(price), price),
  image_urls = COALESCE(sqlc.narg(image_urls), image_urls),
  category = COALESCE(sqlc.narg(category), category),
  discount_percentage = COALESCE(sqlc.narg(discount_percentage), discount_percentage),
  supply_quantity = COALESCE(sqlc.narg(supply_quantity), supply_quantity),
  extra = COALESCE(sqlc.narg(extra), extra)
WHERE
  id = sqlc.arg(item_id)
RETURNING *;