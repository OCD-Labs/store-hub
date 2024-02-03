package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/lib/pq"
)

type ListAllStoresParams struct {
	Search  string
	Filters pagination.Filters
}

type StoreAndOwnersResult struct {
	Store       Store `json:"store"`
	StoreOwners []struct {
		AccountID     string `json:"account_id"`
		ProfileImgURL string `json:"profile_img_url"`
	} `json:"store_owners"`
}

// ListAllStores do a fulltext search to list stores, and paginates accordingly.
func (q *SQLTx) ListAllStores(ctx context.Context, arg ListAllStoresParams) ([]StoreAndOwnersResult, pagination.Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT
		count(*) OVER() AS total_count,
  	s.id, s.name, s.description, s.profile_image_url, s.store_account_id, s.is_verified, s.category, s.is_frozen, s.created_at,
  	json_agg(json_build_object(
      	'account_id', u.account_id,
      	'profile_img_url', u.profile_image_url
  	)) AS store_owners
	FROM 
  	stores AS s
	JOIN 
  	store_owners AS so ON s.id = so.store_id
	JOIN 
  	users AS u ON so.user_id = u.id
	WHERE 
		(name ILIKE '%%' || $1 || '%%' OR $1 = '')
	GROUP BY 
  	s.id
	ORDER BY %s %s, s.id ASC
	LIMIT $2 OFFSET $3`, arg.Filters.SortColumn(), arg.Filters.SortDirection())

	args := []interface{}{arg.Search, arg.Filters.Limit(), arg.Filters.Offset()}

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	results := []StoreAndOwnersResult{}

	for rows.Next() {
		var store Store
		var ownersJSON []byte
		if err := rows.Scan(
			&totalRecords,
			&store.ID,
			&store.Name,
			&store.Description,
			&store.ProfileImageUrl,
			&store.StoreAccountID,
			&store.IsVerified,
			&store.Category,
			&store.IsFrozen,
			&store.CreatedAt,
			&ownersJSON,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}

		var owners []struct {
			AccountID     string `json:"account_id"`
			ProfileImgURL string `json:"profile_img_url"`
		}
		if err := json.Unmarshal(ownersJSON, &owners); err != nil {
			return nil, pagination.Metadata{}, err
		}

		result := StoreAndOwnersResult{
			Store:       store,
			StoreOwners: owners,
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return results, metadata, nil
}

type ListStoreItemsParams struct {
	ItemName     string
	StoreID      int64
	Filters      pagination.Filters
	IsStorefront bool
}

// ListStoreItems do a fulltext search to list store items, and paginates accordingly.
func (q *SQLTx) ListStoreItems(ctx context.Context, arg ListStoreItemsParams) ([]Item, pagination.Metadata, error) {
	statusWhereClause := ""
	if arg.IsStorefront {
		// Only show VISIBLE items to buyers
		statusWhereClause = "AND status = 'VISIBLE'"
	}
	stmt := fmt.Sprintf(`
		SELECT count(*) OVER() AS total_count, items.*
		FROM items
		WHERE (name ILIKE '%%' || $1 || '%%' OR $1 = '') AND store_id = $4 %s
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, statusWhereClause, arg.Filters.SortColumn(), arg.Filters.SortDirection(),
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
			&i.Currency,
			&i.CoverImgUrl,
			&i.Status,
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
	CreatedAtEnd   time.Time
	PaymentChannel string
	DeliveryStatus string
	ItemName       string
	SellerID       int64
	StoreID        int64
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
				AND o.store_id = $%d
    ORDER by o.%s %s, o.id ASC
    LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, len(args)+2, arg.Filters.SortColumn(), arg.Filters.SortDirection(), len(args)+3, len(args)+4)

	args = append(args, arg.SellerID, arg.StoreID, arg.Filters.Limit(), arg.Filters.Offset())

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
	ItemPriceStart    string
	ItemPriceEnd      string
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

	if arg.ItemPriceStart != "" && arg.ItemPriceEnd != "" {
		// Convert strings to float64
		startPrice, startErr := strconv.ParseFloat(arg.ItemPriceStart, 64)
		endPrice, endErr := strconv.ParseFloat(arg.ItemPriceEnd, 64)

		// Check for valid float conversions
		if startErr != nil || endErr != nil {
			return nil, pagination.Metadata{}, fmt.Errorf("invalid price range values")
		} else {
			// Compare and swap if needed
			if startPrice > endPrice {
				startPrice, endPrice = endPrice, startPrice
			}

			whereClauses = append(whereClauses, "i.price BETWEEN $1 AND $2")
			args = append(args, startPrice, endPrice)
		}
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
			&s.ItemPrice,
			&s.ItemCoverImgUrl,
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

type SalesOverviewParams struct {
	ItemName     string
	RevenueStart string
	RevenueEnd   string
	StoreID      int64
	Filters      pagination.Filters
}

type SaleOverviewResult struct {
	SaleID          int64  `json:"sale_id"`
	NumberOfSales   int64  `json:"number_of_sales"`
	SalesPercentage string `json:"sales_percentage"`
	Revenue         string `json:"revenue"`
	ItemID          int64  `json:"item_id"`
	StoreID         int64  `json:"store_id"`
	ItemName        string `json:"item_name"`
	ItemPrice       string `json:"item_price"`
	ItemCoverImgUrl string `json:"item_cover_img_url"`
}

// ListSalesOverview do a full search to list a store's sales overview, and paginates accordingly.
func (dbTx SQLTx) ListSalesOverview(ctx context.Context, arg SalesOverviewParams) ([]SaleOverviewResult, pagination.Metadata, error) {
	var whereClauses []string
	var args []interface{}

	if arg.ItemName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("i.name ILIKE '%%' || $%d || '%%'", len(args)+1))
		args = append(args, arg.ItemName)
	}

	whereClauses, args = addFloatRangeFilter(arg.RevenueStart, arg.RevenueEnd, "so.revenue", whereClauses, args)

	whereClause := strings.Join(whereClauses, " OR ")

	if whereClause == "" {
		whereClause = "TRUE"
	}

	stmt := fmt.Sprintf(`
	SELECT 
		count(*) OVER() AS total_count,
		so.*,
		i.name AS item_name,
		i.cover_img_url AS item_cover_img_url,
		i.price AS item_price
	FROM
		sales_overview so
	JOIN
		items i ON so.item_id = i.id
	WHERE
		(%s)
		AND so.store_id = $%d
	ORDER by so.%s %s, so.id ASC
	LIMIT $%d OFFSET $%d`, whereClause, len(args)+1, arg.Filters.SortColumn(), arg.Filters.SortDirection(), len(args)+2, len(args)+3,
	)

	args = append(args, arg.StoreID, arg.Filters.Limit(), arg.Filters.Offset())

	rows, err := dbTx.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	sors := []SaleOverviewResult{}

	for rows.Next() {
		var sor SaleOverviewResult
		if err := rows.Scan(
			&totalRecords,
			&sor.SaleID,
			&sor.NumberOfSales,
			&sor.SalesPercentage,
			&sor.Revenue,
			&sor.ItemID,
			&sor.StoreID,
			&sor.ItemName,
			&sor.ItemCoverImgUrl,
			&sor.ItemPrice,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}

		sor.SalesPercentage = convertToPercentage(sor.SalesPercentage)
		sors = append(sors, sor)
	}
	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return sors, metadata, nil
}

type ListReviewsParams struct {
	StoreID      int64 `json:"store_id"`
	ItemID       int64 `json:"item_id"`
	IsStorefront bool
	Filters      pagination.Filters
}

type ListReviewsResult struct {
	ID                 int64     `json:"id"`
	StoreID            int64     `json:"store_id"`
	UserID             int64     `json:"user_id"`
	ItemID             int64     `json:"item_id"`
	Rating             string    `json:"rating"`
	ReviewType         string    `json:"review_type"`
	Comment            string    `json:"comment"`
	IsVerifiedPurchase bool      `json:"is_verified_purchase"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	AccountID          string    `json:"account_id"`
	ProfileImageUrl    string    `json:"profile_image_url"`
}

// ListReviews retrieves all the reviews for an item under a store.
func (q *Queries) ListReviews(ctx context.Context, arg ListReviewsParams) ([]ListReviewsResult, pagination.Metadata, error) {
	itemIDClause := ""
	args := []interface{}{arg.StoreID}

	if arg.IsStorefront {
		itemIDClause = "AND r.item_id = $2"
		args = append(args, arg.ItemID)
	}

	stmt := fmt.Sprintf(`
	SELECT
		count(*) OVER() AS total_count,
  	r.*,
  	u.first_name,
  	u.last_name,
  	u.account_id,
  	u.profile_image_url
	FROM 
  	reviews r
	JOIN 
  	users u ON r.user_id = u.id
	WHERE 
  	r.store_id = $1 
  	%s
	ORDER BY %s %s, r.id ASC
	LIMIT %d OFFSET %d`, itemIDClause, arg.Filters.SortColumn(), arg.Filters.SortDirection(), arg.Filters.Limit(), arg.Filters.Offset())

	rows, err := q.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	reviews := []ListReviewsResult{}

	for rows.Next() {
		var r ListReviewsResult
		if err := rows.Scan(
			&totalRecords,
			&r.ID,
			&r.StoreID,
			&r.UserID,
			&r.ItemID,
			&r.Rating,
			&r.ReviewType,
			&r.Comment,
			&r.IsVerifiedPurchase,
			&r.CreatedAt,
			&r.UpdatedAt,
			&r.FirstName,
			&r.LastName,
			&r.AccountID,
			&r.ProfileImageUrl,
		); err != nil {
			return nil, pagination.Metadata{}, err
		}
		reviews = append(reviews, r)
	}

	if err := rows.Err(); err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalcMetadata(totalRecords, arg.Filters.Page, arg.Filters.PageSize)

	return reviews, metadata, nil
}

type ListUserStoresWithAccessRow struct {
	StoreID          int64           `json:"store_id"`
	StoreName        string          `json:"store_name"`
	StoreDescription string          `json:"store_description"`
	StoreImage       string          `json:"store_image"`
	StoreAccountID   string          `json:"store_account_id"`
	IsVerified       bool            `json:"is_verified"`
	Category         string          `json:"category"`
	IsFrozen         bool            `json:"is_frozen"`
	StoreCreatedAt   time.Time       `json:"store_created_at"`
	StoreOwners      json.RawMessage `json:"store_owners"`
}

// ListUserStoresWithAccess retrieves all the stores & its owners for a user
func (q *Queries) ListUserStoresWithAccess(ctx context.Context, userID int64) ([]ListUserStoresWithAccessRow, error) {

	const listUserStoresWithAccess = `-- name: ListUserStoresWithAccess :many
		SELECT 
				store_id,
				store_name,
				store_description,
				store_image,
				store_account_id,
				is_verified,
				category,
				is_frozen,
				store_created_at,
				store_owners
		FROM get_stores_by_user($1)
		`
	rows, err := q.db.QueryContext(ctx, listUserStoresWithAccess, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stores := []ListUserStoresWithAccessRow{}
	for rows.Next() {
		var i ListUserStoresWithAccessRow
		if err := rows.Scan(
			&i.StoreID,
			&i.StoreName,
			&i.StoreDescription,
			&i.StoreImage,
			&i.StoreAccountID,
			&i.IsVerified,
			&i.Category,
			&i.IsFrozen,
			&i.StoreCreatedAt,
			&i.StoreOwners,
		); err != nil {
			return nil, err
		}
		stores = append(stores, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stores, nil
}

func addFloatRangeFilter(startVal, endVal, columnName string, whereClauses []string, args []interface{}) ([]string, []interface{}) {
	if startVal != "" && endVal != "" {
		// Convert strings to float64
		startFloat, startErr := strconv.ParseFloat(startVal, 64)
		endFloat, endErr := strconv.ParseFloat(endVal, 64)

		// Check for valid float conversions
		if startErr != nil || endErr != nil {
			return whereClauses, args
		} else {
			// Compare and swap if needed
			if startFloat > endFloat {
				startFloat, endFloat = endFloat, startFloat
			}

			whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN $%d AND $%d", columnName, len(args)+1, len(args)+2))
			args = append(args, startFloat, endFloat)
		}
	}

	return whereClauses, args
}

func addDateRangeFilter(startDate, endDate time.Time, columnName string, whereClauses []string, args []interface{}) ([]string, []interface{}) {
	if !startDate.IsZero() && !endDate.IsZero() {
		whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN $%d AND $%d", columnName, len(args)+1, len(args)+2))

		if startDate.After(endDate) {
			startDate, endDate = endDate, startDate
		}

		// TODO: this assumes that the end date will always be the start of the day, but nothing guarantees that.
		endDate = endDate.Add(time.Hour*23 + time.Minute*59 + time.Second*59)

		args = append(args, startDate, endDate)
	}
	return whereClauses, args
}

// convertToPercentage converts a string representation of a float64 value to a formatted percentage string.
func convertToPercentage(valueStr string) string {
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return valueStr
	}
	return fmt.Sprintf("%.1f%%", value)
}
