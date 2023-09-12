package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/OCD-Labs/store-hub/pagination"
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

	// CreateStoreTx creates a store and its ownership data.
	CreateStoreTx(ctx context.Context, arg CreateStoreTxParams) (CreateStoreTxResult, error)

	// AddCoOwnerAccessTx creates/update a user's access for a store.
	AddCoOwnerAccessTx(ctx context.Context, arg AddCoOwnerAccessTxParams) (StoreOwner, error)

	// RevokeAccessTx deletes all user's access to a store.
	RevokeAccessTx(ctx context.Context, arg RevokeAccessTxParams) ([]StoreOwner, error)

	// AddCoOwner adds a co-owner to a store.
	AddCoOwnerAccess(ctx context.Context, arg AddCoOwnerAccessParams) (StoreOwner, error)

	// UpdateSellerOrderTx updates a order row, create a sale row if order is DELIVERED.
	UpdateSellerOrderTx(ctx context.Context, arg UpdateSellerOrderParams) (GetOrderForSellerRow, error)

	// CreateReviewTx create a review for an item under a store, updates an order.
	CreateReviewTx(ctx context.Context, arg CreateReviewTxParams) error

	// GetUserCart retrieves a user's cart items.
	GetUserCartTx(ctx context.Context, userID int64) (GetUserCartResult, error)

	// ListAllStores do a fulltext search to list stores, and paginates accordingly.
	ListAllStores(ctx context.Context, arg ListAllStoresParams) ([]StoreAndOwnersResult, pagination.Metadata, error)

	// ListStoreItems do a fulltext search to list store items, and paginates accordingly.
	ListStoreItems(ctx context.Context, arg ListStoreItemsParams) ([]Item, pagination.Metadata, error)

	// CreateUserTx creates a user row and schedules a verify email task on redis.
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (User, error)

	// ListSellerOrders do a fulltext search to list orders, and paginates accordingly.
	ListSellerOrders(ctx context.Context, arg ListSellerOrdersParams) ([]SellerOrder, pagination.Metadata, error)

	// ListAllSellerSales do a fulltext search to list a seller sales, and paginates accordingly.
	ListAllSellerSales(ctx context.Context, arg ListAllSellerSalesParams) ([]GetSaleRow, pagination.Metadata, error)

	// ListSalesOverview do a full search to list a store's sales overview, and paginates accordingly.
	ListSalesOverview(ctx context.Context, arg SalesOverviewParams) ([]SaleOverviewResult, pagination.Metadata, error)

	// ListReviews retrieves all the reviews for an item under a store.
	ListReviews(ctx context.Context, arg ListReviewsParams) ([]ListReviewsResult, pagination.Metadata, error)
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
