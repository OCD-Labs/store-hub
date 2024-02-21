// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: order.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createOrder = `-- name: CreateOrder :one
SELECT id, delivery_status, delivered_on, expected_delivery_date, item_id, order_quantity, buyer_id, seller_id, store_id, delivery_fee, payment_channel, payment_method, is_reviewed, created_at FROM create_order(
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
)
`

type CreateOrderParams struct {
	ItemID         int64  `json:"item_id"`
	OrderQuantity  int32  `json:"order_quantity"`
	BuyerID        int64  `json:"buyer_id"`
	SellerID       int64  `json:"seller_id"`
	StoreID        int64  `json:"store_id"`
	DeliveryFee    string `json:"delivery_fee"`
	PaymentChannel string `json:"payment_channel"`
	PaymentMethod  string `json:"payment_method"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.ItemID,
		arg.OrderQuantity,
		arg.BuyerID,
		arg.SellerID,
		arg.StoreID,
		arg.DeliveryFee,
		arg.PaymentChannel,
		arg.PaymentMethod,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.BuyerID,
		&i.SellerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.IsReviewed,
		&i.CreatedAt,
	)
	return i, err
}

const createOrderFn = `-- name: CreateOrderFn :one
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
) RETURNING id, delivery_status, delivered_on, expected_delivery_date, item_id, order_quantity, buyer_id, seller_id, store_id, delivery_fee, payment_channel, payment_method, is_reviewed, created_at
`

type CreateOrderFnParams struct {
	ItemID         int64  `json:"item_id"`
	OrderQuantity  int32  `json:"order_quantity"`
	BuyerID        int64  `json:"buyer_id"`
	SellerID       int64  `json:"seller_id"`
	StoreID        int64  `json:"store_id"`
	DeliveryFee    string `json:"delivery_fee"`
	PaymentChannel string `json:"payment_channel"`
	PaymentMethod  string `json:"payment_method"`
}

func (q *Queries) CreateOrderFn(ctx context.Context, arg CreateOrderFnParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrderFn,
		arg.ItemID,
		arg.OrderQuantity,
		arg.BuyerID,
		arg.SellerID,
		arg.StoreID,
		arg.DeliveryFee,
		arg.PaymentChannel,
		arg.PaymentMethod,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.BuyerID,
		&i.SellerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.IsReviewed,
		&i.CreatedAt,
	)
	return i, err
}

const getOrderForBuyer = `-- name: GetOrderForBuyer :one
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
  o.id = $1 AND o.buyer_id = $2 AND o.store_id = $3
`

type GetOrderForBuyerParams struct {
	OrderID int64 `json:"order_id"`
	BuyerID int64 `json:"buyer_id"`
	StoreID int64 `json:"store_id"`
}

type GetOrderForBuyerRow struct {
	OrderID              int64     `json:"order_id"`
	DeliveryStatus       string    `json:"delivery_status"`
	DeliveredOn          time.Time `json:"delivered_on"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date"`
	ItemID               int64     `json:"item_id"`
	OrderQuantity        int32     `json:"order_quantity"`
	SellerID             int64     `json:"seller_id"`
	StoreID              int64     `json:"store_id"`
	DeliveryFee          string    `json:"delivery_fee"`
	PaymentChannel       string    `json:"payment_channel"`
	PaymentMethod        string    `json:"payment_method"`
	IsReviewed           bool      `json:"is_reviewed"`
	CreatedAt            time.Time `json:"created_at"`
	ItemName             string    `json:"item_name"`
	ItemDescription      string    `json:"item_description"`
	Price                string    `json:"price"`
	CoverImgUrl          string    `json:"cover_img_url"`
	DiscountPercentage   string    `json:"discount_percentage"`
}

func (q *Queries) GetOrderForBuyer(ctx context.Context, arg GetOrderForBuyerParams) (GetOrderForBuyerRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderForBuyer, arg.OrderID, arg.BuyerID, arg.StoreID)
	var i GetOrderForBuyerRow
	err := row.Scan(
		&i.OrderID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.SellerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.IsReviewed,
		&i.CreatedAt,
		&i.ItemName,
		&i.ItemDescription,
		&i.Price,
		&i.CoverImgUrl,
		&i.DiscountPercentage,
	)
	return i, err
}

const getOrderForSeller = `-- name: GetOrderForSeller :one
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
  o.id = $1 AND o.seller_id = $2 AND o.store_id = $3
`

type GetOrderForSellerParams struct {
	OrderID  int64 `json:"order_id"`
	SellerID int64 `json:"seller_id"`
	StoreID  int64 `json:"store_id"`
}

type GetOrderForSellerRow struct {
	OrderID              int64     `json:"order_id"`
	DeliveryStatus       string    `json:"delivery_status"`
	DeliveredOn          time.Time `json:"delivered_on"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date"`
	ItemID               int64     `json:"item_id"`
	OrderQuantity        int32     `json:"order_quantity"`
	BuyerID              int64     `json:"buyer_id"`
	StoreID              int64     `json:"store_id"`
	DeliveryFee          string    `json:"delivery_fee"`
	PaymentChannel       string    `json:"payment_channel"`
	PaymentMethod        string    `json:"payment_method"`
	CreatedAt            time.Time `json:"created_at"`
	ItemName             string    `json:"item_name"`
	ItemDescription      string    `json:"item_description"`
	Price                string    `json:"price"`
	CoverImgUrl          string    `json:"cover_img_url"`
	DiscountPercentage   string    `json:"discount_percentage"`
	FirstName            string    `json:"first_name"`
	LastName             string    `json:"last_name"`
	Email                string    `json:"email"`
	AccountID            string    `json:"account_id"`
}

func (q *Queries) GetOrderForSeller(ctx context.Context, arg GetOrderForSellerParams) (GetOrderForSellerRow, error) {
	row := q.db.QueryRowContext(ctx, getOrderForSeller, arg.OrderID, arg.SellerID, arg.StoreID)
	var i GetOrderForSellerRow
	err := row.Scan(
		&i.OrderID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.BuyerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.CreatedAt,
		&i.ItemName,
		&i.ItemDescription,
		&i.Price,
		&i.CoverImgUrl,
		&i.DiscountPercentage,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.AccountID,
	)
	return i, err
}

const updateBuyerOrder = `-- name: UpdateBuyerOrder :one
UPDATE orders
SET
  is_reviewed = COALESCE($1, is_reviewed)
WHERE
  id = $2 AND buyer_id = $3 AND store_id = $4
RETURNING id, delivery_status, delivered_on, expected_delivery_date, item_id, order_quantity, buyer_id, seller_id, store_id, delivery_fee, payment_channel, payment_method, is_reviewed, created_at
`

type UpdateBuyerOrderParams struct {
	IsReviewed sql.NullBool `json:"is_reviewed"`
	OrderID    int64        `json:"order_id"`
	BuyerID    int64        `json:"buyer_id"`
	StoreID    int64        `json:"store_id"`
}

func (q *Queries) UpdateBuyerOrder(ctx context.Context, arg UpdateBuyerOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateBuyerOrder,
		arg.IsReviewed,
		arg.OrderID,
		arg.BuyerID,
		arg.StoreID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.BuyerID,
		&i.SellerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.IsReviewed,
		&i.CreatedAt,
	)
	return i, err
}

const updateSellerOrder = `-- name: UpdateSellerOrder :one
UPDATE orders
SET
  delivered_on = COALESCE($1, delivered_on),
  delivery_status = COALESCE($2, delivery_status),
  expected_delivery_date = COALESCE($3, expected_delivery_date)
WHERE
  id = $4 AND seller_id = $5 AND store_id = $6
RETURNING id, delivery_status, delivered_on, expected_delivery_date, item_id, order_quantity, buyer_id, seller_id, store_id, delivery_fee, payment_channel, payment_method, is_reviewed, created_at
`

type UpdateSellerOrderParams struct {
	DeliveredOn          sql.NullTime   `json:"delivered_on"`
	DeliveryStatus       sql.NullString `json:"delivery_status"`
	ExpectedDeliveryDate sql.NullTime   `json:"expected_delivery_date"`
	OrderID              int64          `json:"order_id"`
	SellerID             int64          `json:"seller_id"`
	StoreID              int64          `json:"store_id"`
}

func (q *Queries) UpdateSellerOrder(ctx context.Context, arg UpdateSellerOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateSellerOrder,
		arg.DeliveredOn,
		arg.DeliveryStatus,
		arg.ExpectedDeliveryDate,
		arg.OrderID,
		arg.SellerID,
		arg.StoreID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.DeliveryStatus,
		&i.DeliveredOn,
		&i.ExpectedDeliveryDate,
		&i.ItemID,
		&i.OrderQuantity,
		&i.BuyerID,
		&i.SellerID,
		&i.StoreID,
		&i.DeliveryFee,
		&i.PaymentChannel,
		&i.PaymentMethod,
		&i.IsReviewed,
		&i.CreatedAt,
	)
	return i, err
}
