package db

import (
	"context"
	"fmt"

	"github.com/OCD-Labs/store-hub/pagination"
)

type ListRemindersParamsX struct {
	StoreName string
	Filters   pagination.Filters
}

// ListStoresX do a fulltext search to list stores, and paginates accordingly.
func (q *SQLTx) ListStoresX(ctx context.Context, arg ListRemindersParamsX) ([]Store, pagination.Metadata, error) {
	stmt := fmt.Sprintf(`
		SELECT count(*) OVER() AS total_count, id, name, description, "profile_image_url", is_verified, category, is_frozen, created_at
		FROM stores
		WHERE (name ILIKE '%%' || $1 || '%%' OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, arg.Filters.SortColumn(), arg.Filters.SortDirection(),
	)

	args := []interface{}{arg.StoreName, arg.Filters.Limit(), arg.Filters.Offset()}

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