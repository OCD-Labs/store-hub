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
  i.cover_img_url AS item_cover_img_url,
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

-- name: ReduceSalesOverview :exec
SELECT reduce_sale(sqlc.arg(store_id), sqlc.arg(item_id), sqlc.arg(order_id));

-- name: GetStoreMetrics :one
WITH TodaySales AS (
    SELECT 
        SUM(o.order_quantity) AS today_sales_count,
        SUM(o.order_quantity * i.price) AS today_sales_revenue    FROM sales s
    JOIN orders o ON s.order_id = o.id
    JOIN items i ON s.item_id = i.id
    WHERE s.store_id = sqlc.arg(store_id) AND DATE(s.created_at) = CURRENT_DATE
)

SELECT 
    COALESCE(CAST((SELECT today_sales_count FROM TodaySales) AS TEXT), '0') AS sales_today,
    COALESCE(CAST(SUM(o.order_quantity * i.price) AS TEXT), '0') AS total_sales_revenue,
    COUNT(DISTINCT s.customer_id) AS total_customers,
    COALESCE(SUM(o.order_quantity), 0) AS total_items_sold
FROM sales s
LEFT JOIN orders o ON s.order_id = o.id
LEFT JOIN items i ON s.item_id = i.id
WHERE s.store_id = sqlc.arg(store_id);

-- name: HasMadePurchase :one
SELECT EXISTS(
    SELECT 1 
    FROM sales 
    WHERE customer_id = sqlc.arg(customer_id) AND item_id = sqlc.arg(item_id) AND store_id = sqlc.arg(store_id)
) AS has_made_purchase;
