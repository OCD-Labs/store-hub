package db

import (
	"context"
	"fmt"
	"strings"
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
			&i.Currency,
			&i.CoverImgUrl,
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
	CreatedAtEnd   time.Time
	PaymentChannel string
	DeliveryStatus string
	ItemName       string
	SellerID       int64
	Filters        pagination.Filters
}

type SellerOrder struct {
	OrderID         int64     `json:"order_id"`
	DeliveryStatus  string    `json:"delivery_status"`
	PaymentChannel  string    `json:"payment_channel"`
	CreatedAt       time.Time `json:"created_at"`
	ItemName        string    `json:"item_name"`
	ItemPrice       string    `json:"item_price"`
	ItemCoverImgUrl string    `json:"item_cover_img_url"`
	BuyerFirstName  string    `json:"buyer_first_name"`
	BuyerLastName   string    `json:"buyer_last_name"`
}

// ListSellerOrders do a fulltext search to list a seller orders, and paginates accordingly.
func (q *SQLTx) ListSellerOrders(ctx context.Context, arg ListSellerOrdersParams) ([]SellerOrder, pagination.Metadata, error) {
	var whereClauses []string
	var args []interface{}

	if arg.ItemName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("i.name ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.ItemName)
	}
	if arg.DeliveryStatus != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.delivery_status ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.DeliveryStatus)
	}
	if arg.PaymentChannel != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.payment_channel ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.PaymentChannel)
	}

	whereClauses, args = addDateRangeFilter(arg.CreatedAtStart, arg.CreatedAtEnd, "o.created_at", whereClauses, args)

	whereClause := strings.Join(whereClauses, " OR ")

	if whereClause == "" {
		whereClause = "TRUE"
	}

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
        u.last_name AS buyer_last_name
    FROM
        orders o
    JOIN
        items i ON o.item_id = i.id
    JOIN
        users u ON o.buyer_id = u.id
    WHERE
        (%s)
        AND o.seller_id = $%d
    ORDER by o.%s %s, o.id ASC
    LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, arg.Filters.SortColumn(), arg.Filters.SortDirection(), len(args)+2, len(args)+3)

	args = append(args, arg.SellerID, arg.Filters.Limit(), arg.Filters.Offset())

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

type ListAllSellerSalesParams struct {
	ItemPrice         string
	ItemName          string
	CustomerAccountID string
	DeliveryDateStart time.Time
	DeliveryDateEnd   time.Time
	OrderDateStart    time.Time
	OrderDateEnd      time.Time
	StoreID           int64
	SellerID          int64
	Filters           pagination.Filters
}

// ListAllSellerSales do a fulltext search to list a seller sales, and paginates accordingly.
func (q SQLTx) ListAllSellerSales(ctx context.Context, arg ListAllSellerSalesParams) ([]GetSaleRow, pagination.Metadata, error) {
	var whereClauses []string
	var args []interface{}

	if arg.ItemPrice != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("i.price ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.ItemPrice)
	}
	if arg.ItemName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("i.name ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.ItemName)
	}
	if arg.CustomerAccountID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.account_id ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.CustomerAccountID)
	}
	whereClauses, args = addDateRangeFilter(arg.DeliveryDateStart, arg.DeliveryDateEnd, "o.delivered_on", whereClauses, args)
	whereClauses, args = addDateRangeFilter(arg.OrderDateStart, arg.OrderDateEnd, "o.created_at", whereClauses, args)

	whereClause := strings.Join(whereClauses, " OR ")

	if whereClause == "" {
		whereClause = "TRUE"
	}

	stmt := fmt.Sprintf(`
		SELECT 
			count(*) OVER() AS total_count,
			s.id AS sale_id,
			s.store_id,
			s.created_at,
			s.item_id,
			i.name AS item_name,
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
			(%s)
			AND s.store_id = $%d
  		AND s.seller_id = $%d
		ORDER by s.%s %s, s.id ASC
		LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2, arg.Filters.SortColumn(), arg.Filters.SortDirection(), len(args)+3, len(args)+4,
	)

	args = append(args, arg.StoreID, arg.SellerID, arg.Filters.Limit(), arg.Filters.Offset())

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	ss := []GetSaleRow{}

	for rows.Next() {
		var s GetSaleRow
		if err := rows.Scan(
			&totalRecords,
			&s.SaleID,
			&s.StoreID,
			&s.CreatedAt,
			&s.ItemID,
			&s.ItemName,
			&s.CustomerID,
			&s.CustomerAccountID,
			&s.OrderID,
			&s.OrderDate,
			&s.DeliveryDate,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}
		ss = append(ss, s)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return ss, metadata, nil
}

func addDateRangeFilter(startDate, endDate time.Time, columnName string, whereClauses []string, args []interface{}) ([]string, []interface{}) {
	if !startDate.IsZero() && !endDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN $%d AND $%d", columnName, len(args)+1, len(args)+2))

		if startDate.After(endDate) {
			startDate, endDate = endDate, startDate
		}

		endDate = endDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)

		args = append(args, startDate, endDate)
	}
	return whereClauses, args
}
