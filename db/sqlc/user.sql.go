// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  first_name,
  last_name,
  account_id,
  status,
  hashed_password,
  about,
  email,
  socials,
  profile_image_url
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, first_name, last_name, account_id, status, about, email, socials, profile_image_url, hashed_password, password_changed_at, created_at, is_active, is_email_verified
`

type CreateUserParams struct {
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	AccountID       string          `json:"account_id"`
	Status          string          `json:"status"`
	HashedPassword  string          `json:"hashed_password"`
	About           string          `json:"about"`
	Email           string          `json:"email"`
	Socials         json.RawMessage `json:"socials"`
	ProfileImageUrl sql.NullString  `json:"profile_image_url"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.AccountID,
		arg.Status,
		arg.HashedPassword,
		arg.About,
		arg.Email,
		arg.Socials,
		arg.ProfileImageUrl,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AccountID,
		&i.Status,
		&i.About,
		&i.Email,
		&i.Socials,
		&i.ProfileImageUrl,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsActive,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserByAccountID = `-- name: GetUserByAccountID :one
SELECT id, first_name, last_name, account_id, status, about, email, socials, profile_image_url, hashed_password, password_changed_at, created_at, is_active, is_email_verified FROM users
WHERE account_id = $1 LIMIT 1
`

func (q *Queries) GetUserByAccountID(ctx context.Context, accountID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByAccountID, accountID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AccountID,
		&i.Status,
		&i.About,
		&i.Email,
		&i.Socials,
		&i.ProfileImageUrl,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsActive,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, account_id, status, about, email, socials, profile_image_url, hashed_password, password_changed_at, created_at, is_active, is_email_verified FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, userEmail string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, userEmail)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AccountID,
		&i.Status,
		&i.About,
		&i.Email,
		&i.Socials,
		&i.ProfileImageUrl,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsActive,
		&i.IsEmailVerified,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, first_name, last_name, account_id, status, about, email, socials, profile_image_url, hashed_password, password_changed_at, created_at, is_active, is_email_verified FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, userID int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AccountID,
		&i.Status,
		&i.About,
		&i.Email,
		&i.Socials,
		&i.ProfileImageUrl,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsActive,
		&i.IsEmailVerified,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = COALESCE($1, hashed_password),
  password_changed_at = COALESCE($2, password_changed_at),
  first_name = COALESCE($3, first_name),
  last_name = COALESCE($4, last_name),
  email = COALESCE($5, email),
  is_email_verified = COALESCE($6, is_email_verified),
  is_active = COALESCE($7, is_active),
  profile_image_url = COALESCE($8, profile_image_url),
  socials = COALESCE($9, socials),
  status = COALESCE($10, status),
  about = COALESCE($11, about),
  account_id = COALESCE($12, account_id)
WHERE 
  id = $13 OR email = $5
RETURNING id, first_name, last_name, account_id, status, about, email, socials, profile_image_url, hashed_password, password_changed_at, created_at, is_active, is_email_verified
`

type UpdateUserParams struct {
	HashedPassword    sql.NullString        `json:"hashed_password"`
	PasswordChangedAt sql.NullTime          `json:"password_changed_at"`
	FirstName         sql.NullString        `json:"first_name"`
	LastName          sql.NullString        `json:"last_name"`
	Email             sql.NullString        `json:"email"`
	IsEmailVerified   sql.NullBool          `json:"is_email_verified"`
	IsActive          sql.NullBool          `json:"is_active"`
	ProfileImageUrl   sql.NullString        `json:"profile_image_url"`
	Socials           pqtype.NullRawMessage `json:"socials"`
	Status            sql.NullString        `json:"status"`
	About             sql.NullString        `json:"about"`
	AccountID         sql.NullString        `json:"account_id"`
	ID                sql.NullInt64         `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.IsEmailVerified,
		arg.IsActive,
		arg.ProfileImageUrl,
		arg.Socials,
		arg.Status,
		arg.About,
		arg.AccountID,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AccountID,
		&i.Status,
		&i.About,
		&i.Email,
		&i.Socials,
		&i.ProfileImageUrl,
		&i.HashedPassword,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsActive,
		&i.IsEmailVerified,
	)
	return i, err
}
