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
  description,
  price,
  store_id,
  image_urls,
  category,
  discount_percentage,
  supply_quantity,
  extra
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, created_at, updated_at
`

type CreateStoreItemParams struct {
	Description        string          `json:"description"`
	Price              interface{}     `json:"price"`
	StoreID            int64           `json:"store_id"`
	ImageUrls          []string        `json:"image_urls"`
	Category           string          `json:"category"`
	DiscountPercentage interface{}     `json:"discount_percentage"`
	SupplyQuantity     int64           `json:"supply_quantity"`
	Extra              json.RawMessage `json:"extra"`
}

func (q *Queries) CreateStoreItem(ctx context.Context, arg CreateStoreItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, createStoreItem,
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
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getItem = `-- name: GetItem :one
SELECT id, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, created_at, updated_at FROM items
WHERE id = $1
`

func (q *Queries) GetItem(ctx context.Context, itemID int64) (Item, error) {
	row := q.db.QueryRowContext(ctx, getItem, itemID)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateItem = `-- name: UpdateItem :one
UPDATE items 
SET 
  description = COALESCE($1, description),
  price = COALESCE($2, price),
  image_urls = COALESCE($3, image_urls),
  category = COALESCE($4, category),
  discount_percentage = COALESCE($5, discount_percentage),
  supply_quantity = COALESCE($6, supply_quantity),
  extra = COALESCE($7, extra)
WHERE
  id = $8
RETURNING id, description, price, store_id, image_urls, category, discount_percentage, supply_quantity, extra, created_at, updated_at
`

type UpdateItemParams struct {
	Description        sql.NullString        `json:"description"`
	Price              interface{}           `json:"price"`
	ImageUrls          []string              `json:"image_urls"`
	Category           sql.NullString        `json:"category"`
	DiscountPercentage interface{}           `json:"discount_percentage"`
	SupplyQuantity     sql.NullInt64         `json:"supply_quantity"`
	Extra              pqtype.NullRawMessage `json:"extra"`
	ItemID             int64                 `json:"item_id"`
}

func (q *Queries) UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error) {
	row := q.db.QueryRowContext(ctx, updateItem,
		arg.Description,
		arg.Price,
		pq.Array(arg.ImageUrls),
		arg.Category,
		arg.DiscountPercentage,
		arg.SupplyQuantity,
		arg.Extra,
		arg.ItemID,
	)
	var i Item
	err := row.Scan(
		&i.ID,
		&i.Description,
		&i.Price,
		&i.StoreID,
		pq.Array(&i.ImageUrls),
		&i.Category,
		&i.DiscountPercentage,
		&i.SupplyQuantity,
		&i.Extra,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}