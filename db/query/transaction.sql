-- name: CreateTransaction :one
SELECT * FROM initialize_transaction(
  sqlc.arg(customer_id)::bigint,
  sqlc.arg(amount)::NUMERIC(18, 2),
  sqlc.arg(payment_provider)::varchar,
  sqlc.arg(provider_tx_ref_id)::varchar,
  sqlc.arg(provider_tx_access_code)::varchar
);

-- name: ProcessTransaction :one
SELECT * FROM process_transaction_completion(
  sqlc.arg(provider_tx_ref_id)::varchar,
  sqlc.arg(status)::varchar,
  sqlc.arg(provider_tx_fee)::NUMERIC(10, 2),
  sqlc.arg(cart_items)::jsonb
);

-- name: ReleaseFunds :exec
SELECT release_pending_funds(sqlc.arg(order_id)::bigint);

-- name: GetTransactionByRefID :one
SELECT * FROM transactions 
WHERE provider_tx_ref_id = sqlc.arg(provider_tx_ref_id);

-- name: GetPendingFunds :one
SELECT * FROM pending_transaction_funds
WHERE store_id = sqlc.arg(store_id)
AND account_type = sqlc.arg(account_type);

-- name: ListStorePendingFunds :many
SELECT 
  ptf.*,
  s.name as store_name
FROM pending_transaction_funds ptf
JOIN stores s ON s.id = ptf.store_id
WHERE 
  CASE 
    WHEN sqlc.narg(store_id)::bigint IS NOT NULL THEN ptf.store_id = sqlc.narg(store_id)
    ELSE true
  END
ORDER BY ptf.amount DESC
LIMIT sqlc.arg(rw_limit)
OFFSET sqlc.arg(rw_offset);

-- name: GetTransactionOrders :many
SELECT 
  o.*,
  i.name as item_name,
  i.description as item_description,
  s.name as store_name,
  u.account_id as buyer_account_id
FROM transactions t
JOIN unnest(t.order_ids) AS order_id ON true
JOIN orders o ON o.id = order_id
JOIN items i ON i.id = o.item_id
JOIN stores s ON s.id = o.store_id
JOIN users u ON u.id = o.buyer_id
WHERE t.provider_tx_ref_id = sqlc.arg(provider_tx_ref_id);
