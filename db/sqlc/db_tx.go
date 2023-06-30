package db

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	AnonymousUser = &User{}
)

// IsAnonymous checks if a User instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Store provides all functions to execute db queries
// and transactions.
type StoreTx interface {
	Querier

	// CreateStoreTx creates a store and its ownership data
	CreateStoreTx(ctx context.Context, arg CreateStoreTxParams) (CreateStoreTxResult, error)
}

// A SQLTx provides all functions to execute SQL queries and transactions.
type SQLTx struct {
	*Queries
	db *sql.DB
}

func NewSQLTx(db *sql.DB) StoreTx {
	return &SQLTx{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction,
func (store *SQLTx) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}