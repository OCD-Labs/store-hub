// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: store.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const createStore = `-- name: CreateStore :one
INSERT INTO stores (
  name,
  description,
  profile_image_url,
  store_account_id,
  category
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, name, description, store_account_id, profile_image_url, is_verified, category, is_frozen, created_at
`

type CreateStoreParams struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	StoreAccountID  string `json:"store_account_id"`
	Category        string `json:"category"`
}

func (q *Queries) CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, createStore,
		arg.Name,
		arg.Description,
		arg.ProfileImageUrl,
		arg.StoreAccountID,
		arg.Category,
	)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StoreAccountID,
		&i.ProfileImageUrl,
		&i.IsVerified,
		&i.Category,
		&i.IsFrozen,
		&i.CreatedAt,
	)
	return i, err
}

const deleteStore = `-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = $1
`

func (q *Queries) DeleteStore(ctx context.Context, storeID int64) error {
	_, err := q.db.ExecContext(ctx, deleteStore, storeID)
	return err
}

const getStoreByID = `-- name: GetStoreByID :one
SELECT 
  s.id, s.name, s.description, s.store_account_id, s.profile_image_url, s.is_verified, s.category, s.is_frozen, s.created_at, 
  json_agg(json_build_object(
      'user', json_build_object('id', u.id, 'account_id', u.account_id, 'first_name', u.first_name, 'last_name', u.last_name, 'email', u.email),
      'store_owners', json_build_object('user_id', so.user_id, 'store_id', so.store_id, 'added_at', so.added_at)
  )) AS owners
FROM 
  stores AS s
JOIN 
  store_owners AS so ON s.id = so.store_id
JOIN 
  users AS u ON so.user_id = u.id
WHERE 
  s.id = $1
GROUP BY 
  s.id
`

type GetStoreByIDRow struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	StoreAccountID  string          `json:"store_account_id"`
	ProfileImageUrl string          `json:"profile_image_url"`
	IsVerified      bool            `json:"is_verified"`
	Category        string          `json:"category"`
	IsFrozen        bool            `json:"is_frozen"`
	CreatedAt       time.Time       `json:"created_at"`
	Owners          json.RawMessage `json:"owners"`
}

func (q *Queries) GetStoreByID(ctx context.Context, storeID int64) (GetStoreByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getStoreByID, storeID)
	var i GetStoreByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StoreAccountID,
		&i.ProfileImageUrl,
		&i.IsVerified,
		&i.Category,
		&i.IsFrozen,
		&i.CreatedAt,
		&i.Owners,
	)
	return i, err
}

const getStoreByOwner = `-- name: GetStoreByOwner :many
SELECT s.id, s.name, s.description, s.store_account_id, s.profile_image_url, s.is_verified, s.category, s.is_frozen, s.created_at
FROM stores s
JOIN store_owners so ON s.id = so.store_id
WHERE so.user_id = $1
`

func (q *Queries) GetStoreByOwner(ctx context.Context, userID int64) ([]Store, error) {
	rows, err := q.db.QueryContext(ctx, getStoreByOwner, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Store{}
	for rows.Next() {
		var i Store
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.StoreAccountID,
			&i.ProfileImageUrl,
			&i.IsVerified,
			&i.Category,
			&i.IsFrozen,
			&i.CreatedAt,
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

const updateStore = `-- name: UpdateStore :one
UPDATE stores
SET
  name = COALESCE($1, name),
  description = COALESCE($2, description),
  profile_image_url = COALESCE($3, profile_image_url),
  is_verified = COALESCE($4, is_verified),
  category = COALESCE($5, category),
  is_frozen = COALESCE($6, is_frozen)
WHERE 
  id = $7
RETURNING id, name, description, store_account_id, profile_image_url, is_verified, category, is_frozen, created_at
`

type UpdateStoreParams struct {
	Name            sql.NullString `json:"name"`
	Description     sql.NullString `json:"description"`
	ProfileImageUrl sql.NullString `json:"profile_image_url"`
	IsVerified      sql.NullBool   `json:"is_verified"`
	Category        sql.NullString `json:"category"`
	IsFrozen        sql.NullBool   `json:"is_frozen"`
	StoreID         int64          `json:"store_id"`
}

func (q *Queries) UpdateStore(ctx context.Context, arg UpdateStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, updateStore,
		arg.Name,
		arg.Description,
		arg.ProfileImageUrl,
		arg.IsVerified,
		arg.Category,
		arg.IsFrozen,
		arg.StoreID,
	)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.StoreAccountID,
		&i.ProfileImageUrl,
		&i.IsVerified,
		&i.Category,
		&i.IsFrozen,
		&i.CreatedAt,
	)
	return i, err
}
