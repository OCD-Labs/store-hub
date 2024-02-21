// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: carts.sql

package db

import (
	"context"
)

const createCartForUser = `-- name: CreateCartForUser :exec
INSERT INTO carts (
  user_id
) VALUES (
  $1
) RETURNING id, user_id, created_at
`

func (q *Queries) CreateCartForUser(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, createCartForUser, userID)
	return err
}

const decreaseCartItemQuantity = `-- name: DecreaseCartItemQuantity :one
UPDATE cart_items 
SET 
  quantity = GREATEST(1, cart_items.quantity - $1)
WHERE 
  cart_id = $2 AND item_id = $3
RETURNING id, cart_id, item_id, store_id, quantity, added_at, updated_at
`

type DecreaseCartItemQuantityParams struct {
	DecreaseAmount int32 `json:"decrease_amount"`
	CartID         int64 `json:"cart_id"`
	ItemID         int64 `json:"item_id"`
}

func (q *Queries) DecreaseCartItemQuantity(ctx context.Context, arg DecreaseCartItemQuantityParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, decreaseCartItemQuantity, arg.DecreaseAmount, arg.CartID, arg.ItemID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.ItemID,
		&i.StoreID,
		&i.Quantity,
		&i.AddedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCartByUserID = `-- name: GetCartByUserID :many
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
  c.user_id = $1
`

type GetCartByUserIDRow struct {
	CartID          int64  `json:"cart_id"`
	ItemID          int64  `json:"item_id"`
	ItemName        string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	Price           string `json:"price"`
	Quantity        int32  `json:"quantity"`
	ItemImage       string `json:"item_image"`
}

func (q *Queries) GetCartByUserID(ctx context.Context, userID int64) ([]GetCartByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCartByUserIDRow{}
	for rows.Next() {
		var i GetCartByUserIDRow
		if err := rows.Scan(
			&i.CartID,
			&i.ItemID,
			&i.ItemName,
			&i.ItemDescription,
			&i.Price,
			&i.Quantity,
			&i.ItemImage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCartID = `-- name: GetCartID :one
SELECT 
  id
FROM 
  carts 
WHERE 
  user_id = $1 LIMIT 1
`

func (q *Queries) GetCartID(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCartID, userID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const increaseCartItemQuantity = `-- name: IncreaseCartItemQuantity :one
WITH item_supply AS (
  SELECT supply_quantity 
  FROM items 
  WHERE id = $3
)
UPDATE cart_items 
SET 
  quantity = LEAST((SELECT supply_quantity FROM item_supply), quantity + $1)
WHERE 
  cart_id = $2 
  AND item_id = $3
RETURNING id, cart_id, item_id, store_id, quantity, added_at, updated_at
`

type IncreaseCartItemQuantityParams struct {
	IncreaseAmount int32 `json:"increase_amount"`
	CartID         int64 `json:"cart_id"`
	ItemID         int64 `json:"item_id"`
}

func (q *Queries) IncreaseCartItemQuantity(ctx context.Context, arg IncreaseCartItemQuantityParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, increaseCartItemQuantity, arg.IncreaseAmount, arg.CartID, arg.ItemID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.ItemID,
		&i.StoreID,
		&i.Quantity,
		&i.AddedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeItemFromCart = `-- name: RemoveItemFromCart :exec
DELETE FROM cart_items
WHERE cart_id = $1 AND item_id = $2
`

type RemoveItemFromCartParams struct {
	CartID int64 `json:"cart_id"`
	ItemID int64 `json:"item_id"`
}

func (q *Queries) RemoveItemFromCart(ctx context.Context, arg RemoveItemFromCartParams) error {
	_, err := q.db.ExecContext(ctx, removeItemFromCart, arg.CartID, arg.ItemID)
	return err
}

const upsertCartItem = `-- name: UpsertCartItem :one
SELECT id, cart_id, item_id, store_id, quantity, added_at, updated_at FROM upsert_cart_item(
  $1::bigint,
  $2::bigint,
  $3::bigint
)
`

type UpsertCartItemParams struct {
	CartID  int64 `json:"cart_id"`
	ItemID  int64 `json:"item_id"`
	StoreID int64 `json:"store_id"`
}

func (q *Queries) UpsertCartItem(ctx context.Context, arg UpsertCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, upsertCartItem, arg.CartID, arg.ItemID, arg.StoreID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.CartID,
		&i.ItemID,
		&i.StoreID,
		&i.Quantity,
		&i.AddedAt,
		&i.UpdatedAt,
	)
	return i, err
}
