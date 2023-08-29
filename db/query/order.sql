-- name: CreateOrder :one
INSERT INTO orders (
  item_id,
  order_quantity,
  buyer_id,
  seller_id,
  store_id,
  delivery_fee,
  payment_channel,
  payment_method
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetOrderForSeller :one
SELECT
  o.id AS order_id,
  o.delivery_status,
  o.delivered_on,
  o.expected_delivery_date,
  o.item_id,
  o.order_quantity,
  o.buyer_id,
  o.store_id,
  o.delivery_fee,
  o.payment_channel,
  o.payment_method,
  o.created_at,
  i.name AS item_name,
  i.description AS item_description,
  i.price,
  i.cover_img_url,
  i.discount_percentage,
  u.first_name,
  u.last_name,
  u.email,
  u.account_id
FROM
  orders o
JOIN
  items i ON o.item_id = i.id
JOIN
  users u ON o.buyer_id = u.id
WHERE 
  o.id = sqlc.arg(order_id) AND o.seller_id = sqlc.arg(seller_id) AND o.store_id = sqlc.arg(store_id);

-- name: UpdateSellerOrder :one
UPDATE orders
SET
  delivered_on = COALESCE(sqlc.narg(delivered_on), delivered_on),
  delivery_status = COALESCE(sqlc.narg(delivery_status), delivery_status),
  expected_delivery_date = COALESCE(sqlc.narg(expected_delivery_date), expected_delivery_date)
WHERE
  id = sqlc.arg(order_id) AND seller_id = sqlc.arg(seller_id) AND store_id = sqlc.arg(store_id)
RETURNING *;

-- name: UpdateBuyerOrder :one
UPDATE orders
SET
  is_reviewed = COALESCE(sqlc.narg(is_reviewed), is_reviewed)
WHERE
  id = sqlc.arg(order_id) AND buyer_id = sqlc.arg(buyer_id) AND store_id = sqlc.arg(store_id)
RETURNING *;


-- name: GetOrderForBuyer :one
SELECT
  o.id AS order_id,
  o.delivery_status,
  o.delivered_on,
  o.expected_delivery_date,
  o.item_id,
  o.order_quantity,
  o.seller_id,
  o.store_id,
  o.delivery_fee,
  o.payment_channel,
  o.payment_method,
  o.is_reviewed,
  o.created_at,
  i.name AS item_name,
  i.description AS item_description,
  i.price,
  i.cover_img_url,
  i.discount_percentage
FROM
  orders o
JOIN
  items i ON o.item_id = i.id
WHERE 
  o.id = sqlc.arg(order_id) AND o.buyer_id = sqlc.arg(buyer_id) AND o.store_id = sqlc.arg(store_id);