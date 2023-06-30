// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: item.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/tabbed/pqtype"
)

const createStoreItem = `-- name: CreateStoreItem :one
INSERT INTO items (
  name,
  description,
  price,
  store_id,
  image_urls,
  category,
  discount_percentage,
  supply_quantity,
  extra
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, name, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, is_frozen, created_at, updated_at
`

type CreateStoreItemParams struct {
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Price              string          `json:"price"`
	StoreID            int64           `json:"store_id"`
	ImageUrls          []string        `json:"image_urls"`
	Category           string          `json:"category"`
	DiscountPercentage string          `json:"discount_percentage"`
	SupplyQuantity     int64           `json:"supply_quantity"`
	Extra              json.RawMessage `json:"extra"`
}

func (q *Queries) CreateStoreItem(ctx context.Context, arg CreateStoreItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createStoreItem,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.StoreID,
		pq.Array(arg.ImageUrls),
		arg.Category,
		arg.DiscountPercentage,
		arg.SupplyQuantity,
		arg.Extra,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.IsFrozen,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteItem = `-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1
`

func (q *Queries) DeleteItem(ctx context.Context, itemID int64) error {
	_, err := q.db.ExecContext(ctx, deleteItem, itemID)
	return err
}

const getItem = `-- name: GetItem :one
SELECT id, name, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, is_frozen, created_at, updated_at FROM items
WHERE id = $1 AND supply_quantity > 0
`

func (q *Queries) GetItem(ctx context.Context, itemID int64) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItem, itemID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.IsFrozen,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateItem = `-- name: UpdateItem :one
UPDATE items 
SET 
  description = COALESCE($1, description),
  name = COALESCE($2, name),
  price = COALESCE($3, price),
  image_urls = COALESCE($4, image_urls),
  category = COALESCE($5, category),
  discount_percentage = COALESCE($6, discount_percentage),
  supply_quantity = COALESCE($7, supply_quantity),
  extra = COALESCE($8, extra),
  is_frozen = COALESCE($9, is_frozen),
  updated_at = COALESCE($10, updated_at)
WHERE
  id = $11
RETURNING id, name, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, is_frozen, created_at, updated_at
`

type UpdateItemParams struct {
	Description        sql.NullString        `json:"description"`
	Name               sql.NullString        `json:"name"`
	Price              sql.NullString        `json:"price"`
	ImageUrls          []string              `json:"image_urls"`
	Category           sql.NullString        `json:"category"`
	DiscountPercentage sql.NullString        `json:"discount_percentage"`
	SupplyQuantity     sql.NullInt64         `json:"supply_quantity"`
	Extra              pqtype.NullRawMessage `json:"extra"`
	IsFrozen           sql.NullBool          `json:"is_frozen"`
	UpdatedAt          sql.NullTime          `json:"updated_at"`
	ItemID             int64                 `json:"item_id"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, updateItem,
		arg.Description,
		arg.Name,
		arg.Price,
		pq.Array(arg.ImageUrls),
		arg.Category,
		arg.DiscountPercentage,
		arg.SupplyQuantity,
		arg.Extra,
		arg.IsFrozen,
		arg.UpdatedAt,
		arg.ItemID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.IsFrozen,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
