-- name: CreateStoreItem :one
INSERT INTO items (
  name,
  description,
  price,
  store_id,
  image_urls,
  category,
  cover_img_url,
  discount_percentage,
  supply_quantity,
  extra,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = sqlc.arg(item_id) AND supply_quantity > 0;

-- name: UpdateItem :one
UPDATE items 
SET 
  description = COALESCE(sqlc.narg(description), description),
  name = COALESCE(sqlc.narg(name), name),
  price = COALESCE(sqlc.narg(price), price),
  image_urls = COALESCE(sqlc.narg(image_urls), image_urls),
  cover_img_url = COALESCE(sqlc.narg(cover_img_url), cover_img_url),
  category = COALESCE(sqlc.narg(category), category),
  discount_percentage = COALESCE(sqlc.narg(discount_percentage), discount_percentage),
  supply_quantity = COALESCE(sqlc.narg(supply_quantity), supply_quantity),
  extra = COALESCE(sqlc.narg(extra), extra),
  is_frozen = COALESCE(sqlc.narg(is_frozen), is_frozen),
  status = COALESCE(sqlc.narg(status), status),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE
  id = sqlc.arg(item_id)
RETURNING *;

-- name: DeductItemSupply :exec
UPDATE items 
SET 
  supply_quantity = supply_quantity - sqlc.arg(order_quantity)
WHERE
  id = sqlc.arg(item_id) AND supply_quantity >= sqlc.arg(order_quantity);

-- name: DeleteItem :exec
DELETE FROM items
WHERE store_id = sqlc.arg(store_id) AND id = sqlc.arg(item_id);

-- name: CheckItemStoreMatch :one
SELECT supply_quantity from items
WHERE id = sqlc.arg(item_id)
  AND store_id = sqlc.arg(store_id);