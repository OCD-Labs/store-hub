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
SELECT 
  s.id AS sale_id,
  s.store_id,
  s.created_at,
  s.item_id,
  i.name AS item_name,
  i.price AS item_price,
  s.customer_id,
  u.account_id AS customer_account_id,
  s.order_id,
  o.created_at AS order_date,
  o.delivered_on AS delivery_date
FROM 
  sales s
JOIN
  users u ON s.customer_id = u.id
JOIN
  items i ON s.item_id = i.id
JOIN 
  orders o ON s.order_id = o.id
WHERE 
  s.id = sqlc.arg(sale_id)
  AND s.store_id = sqlc.arg(store_id)
  AND s.seller_id = sqlc.arg(seller_id);

-- name: SaleExists :one
SELECT EXISTS (
    SELECT 1
    FROM sales
    WHERE order_id = sqlc.arg(order_id)
);
