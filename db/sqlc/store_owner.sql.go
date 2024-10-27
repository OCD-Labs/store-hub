// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: store_owner.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const addCoOwnerAccess = `-- name: AddCoOwnerAccess :one
INSERT INTO store_owners (
  store_id,
  user_id,
  access_levels,
  is_primary
) VALUES (
  $1, $2, $3, $4
) RETURNING user_id, store_id, access_levels, is_primary, added_at
`

type AddCoOwnerAccessParams struct {
	StoreID      int64   `json:"store_id"`
	UserID       int64   `json:"user_id"`
	AccessLevels []int32 `json:"access_levels"`
	IsPrimary    bool    `json:"is_primary"`
}

func (q *Queries) AddCoOwnerAccess(ctx context.Context, arg AddCoOwnerAccessParams) (StoreOwner, error) {
	row := q.db.QueryRowContext(ctx, addCoOwnerAccess,
		arg.StoreID,
		arg.UserID,
		pq.Array(arg.AccessLevels),
		arg.IsPrimary,
	)
	var i StoreOwner
	err := row.Scan(
		&i.UserID,
		&i.StoreID,
		pq.Array(&i.AccessLevels),
		&i.IsPrimary,
		&i.AddedAt,
	)
	return i, err
}

const addToCoOwnerAccess = `-- name: AddToCoOwnerAccess :one
UPDATE store_owners 
SET 
  access_levels = array_append(access_levels, $1)
WHERE 
  store_id = $2 AND user_id = $3
RETURNING user_id, store_id, access_levels, is_primary, added_at
`

type AddToCoOwnerAccessParams struct {
	NewAccessLevel interface{} `json:"new_access_level"`
	StoreID        int64       `json:"store_id"`
	UserID         int64       `json:"user_id"`
}

func (q *Queries) AddToCoOwnerAccess(ctx context.Context, arg AddToCoOwnerAccessParams) (StoreOwner, error) {
	row := q.db.QueryRowContext(ctx, addToCoOwnerAccess, arg.NewAccessLevel, arg.StoreID, arg.UserID)
	var i StoreOwner
	err := row.Scan(
		&i.UserID,
		&i.StoreID,
		pq.Array(&i.AccessLevels),
		&i.IsPrimary,
		&i.AddedAt,
	)
	return i, err
}

const getStoreOwnersByStoreID = `-- name: GetStoreOwnersByStoreID :many
SELECT user_id, store_id, access_levels, is_primary, added_at
FROM store_owners
WHERE store_id = $1
`

func (q *Queries) GetStoreOwnersByStoreID(ctx context.Context, storeID int64) ([]StoreOwner, error) {
	rows, err := q.db.QueryContext(ctx, getStoreOwnersByStoreID, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreOwner{}
	for rows.Next() {
		var i StoreOwner
		if err := rows.Scan(
			&i.UserID,
			&i.StoreID,
			pq.Array(&i.AccessLevels),
			&i.IsPrimary,
			&i.AddedAt,
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

const getUserAccessLevelsForStore = `-- name: GetUserAccessLevelsForStore :one
SELECT access_levels
FROM store_owners
WHERE user_id = $1
  AND store_id = $2
`

type GetUserAccessLevelsForStoreParams struct {
	UserID  int64 `json:"user_id"`
	StoreID int64 `json:"store_id"`
}

func (q *Queries) GetUserAccessLevelsForStore(ctx context.Context, arg GetUserAccessLevelsForStoreParams) ([]int32, error) {
	row := q.db.QueryRowContext(ctx, getUserAccessLevelsForStore, arg.UserID, arg.StoreID)
	var access_levels []int32
	err := row.Scan(pq.Array(&access_levels))
	return access_levels, err
}

const revokeAccess = `-- name: RevokeAccess :exec
UPDATE store_owners 
SET access_levels = ARRAY_REMOVE(access_levels, $1::int)
WHERE 
  user_id = $2 AND store_id = $3
`

type RevokeAccessParams struct {
	AccessLevelToRevoke int32 `json:"access_level_to_revoke"`
	UserID              int64 `json:"user_id"`
	StoreID             int64 `json:"store_id"`
}

func (q *Queries) RevokeAccess(ctx context.Context, arg RevokeAccessParams) error {
	_, err := q.db.ExecContext(ctx, revokeAccess, arg.AccessLevelToRevoke, arg.UserID, arg.StoreID)
	return err
}

const revokeAllAccess = `-- name: RevokeAllAccess :exec
DELETE FROM store_owners
WHERE user_id = $1 AND store_id = $2
`

type RevokeAllAccessParams struct {
	UserID  int64 `json:"user_id"`
	StoreID int64 `json:"store_id"`
}

func (q *Queries) RevokeAllAccess(ctx context.Context, arg RevokeAllAccessParams) error {
	_, err := q.db.ExecContext(ctx, revokeAllAccess, arg.UserID, arg.StoreID)
	return err
}
