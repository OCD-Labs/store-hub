package db

import (
	"context"
	"fmt"

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
		SELECT count(*) OVER() AS total_count, id, name, description, "profile_image_url", is_verified, category, is_frozen, created_at
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