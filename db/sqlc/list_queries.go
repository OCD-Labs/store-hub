package db

import (
	"context"
	"fmt"
	"time"

	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/lib/pq"
)

type ListAllStoresParams struct {
	Search  string
	Filters pagination.Filters
}

// ListAllStores do a fulltext search to list stores, and paginates accordingly.
func (q *SQLTx) ListAllStores(ctx context.Context, arg ListAllStoresParams) ([]Store, pagination.Metadata, error) {
	stmt := fmt.Sprintf(`
		SELECT count(*) OVER() AS total_count, id, name, description, "profile_image_url", store_account_id, is_verified, category, is_frozen, created_at
		FROM stores
		WHERE (name ILIKE '%%' || $1 || '%%' OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, arg.Filters.SortColumn(), arg.Filters.SortDirection(),
	)

	args := []interface{}{arg.Search, arg.Filters.Limit(), arg.Filters.Offset()}

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	stores := []Store{}

	for rows.Next() {
		var i Store
		if err := rows.Scan(
			&totalRecords,
			&i.ID,
			&i.Name,
			&i.Description,
			&i.ProfileImageUrl,
			&i.StoreAccountID,
			&i.IsVerified,
			&i.Category,
			&i.IsFrozen,
			&i.CreatedAt,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}
		stores = append(stores, i)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return stores, metadata, nil
}

type ListStoreItemsParams struct {
	ItemName string
	StoreID  int64
	Filters  pagination.Filters
}

// ListStoreItems do a fulltext search to list store items, and paginates accordingly.
func (q *SQLTx) ListStoreItems(ctx context.Context, arg ListStoreItemsParams) ([]Item, pagination.Metadata, error) {
	stmt := fmt.Sprintf(`
		SELECT count(*) OVER() AS total_count, items.*
		FROM items
		WHERE (name ILIKE '%%' || $1 || '%%' OR $1 = '') AND store_id = $4
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, arg.Filters.SortColumn(), arg.Filters.SortDirection(),
	)

	args := []interface{}{arg.ItemName, arg.Filters.Limit(), arg.Filters.Offset(), arg.StoreID}

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	items := []Item{}

	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&totalRecords,
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
		); err != nil {
			return nil, pagination.Metadata{}, err
		}
		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return items, metadata, nil
}

type ListSellerOrdersParams struct {
	CreatedAtStart time.Time
	CreatedAtEnd time.Time
	PaymentChannel string
	Price string
	DeliveryStatus string
	ItemName string
	SellerID int64
	Filters pagination.Filters
}

type SellerOrder struct {
	OrderID int64 `json:"order_id"`
	DeliveryStatus string `json:"delivery_status"`
	PaymentChannel string `json:"payment_channel"`
	CreatedAt time.Time `json:"created_at"`
	ItemName string `json:"item_name"`
	ItemPrice string `json:"item_price"`
	ItemCoverImgUrl string `json:"item_cover_img_url"`
	BuyerFirstName string `json:"buyer_first_name"`
	BuyerLastName string `json:"buyer_last_name"`
}

// ListSellerOrders do a fulltext search to list a seller orders, and paginates accordingly.
func (q *SQLTx) ListSellerOrders(ctx context.Context, arg ListSellerOrdersParams) ([]SellerOrder, pagination.Metadata, error) {
	stmt := fmt.Sprintf(`
		SELECT
			count(*) OVER() AS total_count,
			o.id AS order_id,
  		o.delivery_status,
  		o.payment_channel,
			o.created_at,
  		i.name AS item_name,
  		i.price AS item_price,
  		i.cover_img_url AS item_cover_img_url,
  		u.first_name AS buyer_first_name,
  		u.last_name AS buyer_last_name,
		FROM
			orders o
		JOIN
			items i ON o.item_id = i.id
		JOIN
			users u ON o.buyer_id = u.id
		WHERE
			(
				($1 <> '' AND i.name ILIKE '%%' || $1 || '%%') OR
				($2 <> '' AND o.status ILIKE '%%' || $2 || '%%') OR
				($3 <> '' AND o.payment_channel ILIKE '%%' || $3 || '%%') OR
				(o.created_at BETWEEN $4 AND $5) OR
			)
			AND o.seller_id = $6
		ORDER by %s %s, id ASC
		LIMIT $7 OFFSET $8`, arg.Filters.SortColumn(), arg.Filters.SortDirection(),
	)

	args := []interface{}{
		arg.ItemName,
		arg.DeliveryStatus,
		arg.PaymentChannel,
		arg.CreatedAtStart,
		arg.CreatedAtEnd,
		arg.SellerID,
		arg.Filters.Limit(),
		arg.Filters.Offset(), 
	}

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	sos := []SellerOrder{}

	for rows.Next() {
		var so SellerOrder
		if err := rows.Scan(
			&totalRecords,
			&so.OrderID,
			&so.DeliveryStatus,
			&so.PaymentChannel,
			&so.CreatedAt,
			&so.ItemName,
			&so.ItemPrice,
			&so.ItemCoverImgUrl,
			&so.BuyerFirstName,
			&so.BuyerLastName,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}
		sos = append(sos, so)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)


	return sos, metadata, nil
}