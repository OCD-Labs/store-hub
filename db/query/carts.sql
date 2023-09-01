-- name: CreateCartForUser :exec
INSERT INTO carts (
  user_id
) VALUES (
  $1
) RETURNING *;

-- name: GetCartByUserID :many
SELECT 
  c.id AS cart_id,
  ci.item_id,
  i.name AS item_name,
  i.description AS item_description,
  i.price,
  ci.quantity,
  i.cover_img_url AS item_image
FROM 
  carts c
JOIN 
  cart_items ci ON c.id = ci.cart_id
JOIN 
  items i ON ci.item_id = i.id
WHERE 
  c.user_id = sqlc.arg(user_id);

-- name: UpsertCartItem :one
SELECT * FROM upsert_cart_item(
  sqlc.arg(cart_id)::bigint,
  sqlc.arg(item_id)::bigint,
  sqlc.arg(store_id)::bigint
);

-- name: RemoveItemFromCart :exec
DELETE FROM cart_items
WHERE cart_id = sqlc.arg(cart_id) AND item_id = sqlc.arg(item_id);

-- name: IncreaseCartItemQuantity :one
WITH item_supply AS (
  SELECT supply_quantity 
  FROM items 
  WHERE id = sqlc.arg(item_id)
)
UPDATE cart_items 
SET 
  quantity = LEAST((SELECT supply_quantity FROM item_supply), quantity + sqlc.arg(increase_amount))
WHERE 
  cart_id = sqlc.arg(cart_id) 
  AND item_id = sqlc.arg(item_id)
RETURNING *;

-- name: DecreaseCartItemQuantity :one
UPDATE cart_items 
SET 
  quantity = GREATEST(1, cart_items.quantity - sqlc.arg(decrease_amount))
WHERE 
  cart_id = sqlc.arg(cart_id) AND item_id = sqlc.arg(item_id)
RETURNING *;
