-- name: CreateSale :one
INSERT INTO sales (
  store_id,
  item_id,
  customer_id,
  seller_id,
  order_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSale :one
SELECT * FROM sales
WHERE id = sqlc.arg(sale_id) AND store_id = sqlc.arg(store_id);