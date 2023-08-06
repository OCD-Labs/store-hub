// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CheckSessionExistence(ctx context.Context, token string) (bool, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error)
	CreateStoreItem(ctx context.Context, arg CreateStoreItemParams) (Item, error)
	CreateStoreOwner(ctx context.Context, arg CreateStoreOwnerParams) (StoreOwner, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteExpiredSession(ctx context.Context) error
	DeleteItem(ctx context.Context, arg DeleteItemParams) error
	DeleteStore(ctx context.Context, storeID int64) error
	DeleteStoreOwner(ctx context.Context, arg DeleteStoreOwnerParams) error
	GetItem(ctx context.Context, itemID int64) (Item, error)
	GetOrderForSeller(ctx context.Context, arg GetOrderForSellerParams) (GetOrderForSellerRow, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetStoreByID(ctx context.Context, storeID int64) (GetStoreByIDRow, error)
	GetStoreByOwner(ctx context.Context, userID int64) ([]Store, error)
	GetUserByAccountID(ctx context.Context, accountID string) (User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (User, error)
	GetUserByID(ctx context.Context, userID int64) (User, error)
	IsStoreOwner(ctx context.Context, arg IsStoreOwnerParams) (IsStoreOwnerRow, error)
	UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateStore(ctx context.Context, arg UpdateStoreParams) (Store, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
