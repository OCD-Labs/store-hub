-- name: CreateOrder :one
INSERT INTO orders (
  delivery_status,
  item_id,
  order_quantity,
  store_id,
  buyer_id, -- update by checking parameters that are missing.
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
  o.delivered_date,
  o.item_id,
  o.order_quantity,
  o.buyer_id,
  o.store_id,
  o.delivery_fee,
  o.payment_channel,
  o.payment_method,
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
    o.id = $1;

-- name: UpdateOrder :one
UPDATE orders
SET
  delivered_date = COALESCE(sqlc.narg(delivered_date), delivered_date),
  delivery_status = COALESCE(sqlc.narg(delivery_status), delivery_status),
  expected_delivery_date = COALESCE(sqlc.narg(expected_delivery_date), expected_delivery_date)
WHERE
  id = sqlc.arg(order_id)
RETURNING *;
