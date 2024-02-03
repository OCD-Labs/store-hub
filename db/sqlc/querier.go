// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	AddCoOwnerAccess(ctx context.Context, arg AddCoOwnerAccessParams) (StoreOwner, error)
	AddToCoOwnerAccess(ctx context.Context, arg AddToCoOwnerAccessParams) (StoreOwner, error)
	CheckItemStoreMatch(ctx context.Context, arg CheckItemStoreMatchParams) (int64, error)
	CheckSessionExists(ctx context.Context, arg CheckSessionExistsParams) (bool, error)
	CreateCartForUser(ctx context.Context, userID int64) error
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderFn(ctx context.Context, arg CreateOrderFnParams) (Order, error)
	CreateReview(ctx context.Context, arg CreateReviewParams) error
	CreateReviewFn(ctx context.Context, arg CreateReviewFnParams) error
	CreateSale(ctx context.Context, arg CreateSaleParams) (Sale, error)
	CreateSaleFn(ctx context.Context, arg CreateSaleFnParams) (Sale, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error)
	CreateStoreItem(ctx context.Context, arg CreateStoreItemParams) (Item, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DecreaseCartItemQuantity(ctx context.Context, arg DecreaseCartItemQuantityParams) (CartItem, error)
	DeductItemSupply(ctx context.Context, arg DeductItemSupplyParams) error
	DeleteExpiredSession(ctx context.Context) error
	DeleteItem(ctx context.Context, arg DeleteItemParams) error
	DeleteReview(ctx context.Context, arg DeleteReviewParams) error
	DeleteStore(ctx context.Context, storeID int64) error
	GetCartByUserID(ctx context.Context, userID int64) ([]GetCartByUserIDRow, error)
	GetCartID(ctx context.Context, userID int64) (int64, error)
	GetItem(ctx context.Context, itemID int64) (Item, error)
	GetOrderForBuyer(ctx context.Context, arg GetOrderForBuyerParams) (GetOrderForBuyerRow, error)
	GetOrderForSeller(ctx context.Context, arg GetOrderForSellerParams) (GetOrderForSellerRow, error)
	GetSale(ctx context.Context, arg GetSaleParams) (GetSaleRow, error)
	GetSession(ctx context.Context, id uuid.UUID) (Session, error)
	GetStoreByID(ctx context.Context, storeID int64) (GetStoreByIDRow, error)
	GetStoreMetrics(ctx context.Context, storeID int64) (GetStoreMetricsRow, error)
	GetStoreOwnersByStoreID(ctx context.Context, storeID int64) ([]StoreOwner, error)
	GetUserAccessLevelsForStore(ctx context.Context, arg GetUserAccessLevelsForStoreParams) ([]int32, error)
	GetUserByAccountID(ctx context.Context, accountID string) (User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (User, error)
	GetUserByID(ctx context.Context, userID int64) (User, error)
	HasMadePurchase(ctx context.Context, arg HasMadePurchaseParams) (bool, error)
	IncreaseCartItemQuantity(ctx context.Context, arg IncreaseCartItemQuantityParams) (CartItem, error)
	LogAction(ctx context.Context, arg LogActionParams) error
	RatingOverview(ctx context.Context, storeID int64) (RatingOverviewRow, error)
	ReduceSalesOverview(ctx context.Context, arg ReduceSalesOverviewParams) error
	RemoveItemFromCart(ctx context.Context, arg RemoveItemFromCartParams) error
	RevokeAccess(ctx context.Context, arg RevokeAccessParams) error
	RevokeAllAccess(ctx context.Context, arg RevokeAllAccessParams) error
	UpdateBuyerOrder(ctx context.Context, arg UpdateBuyerOrderParams) (Order, error)
	UpdateItem(ctx context.Context, arg UpdateItemParams) (Item, error)
	UpdateSellerOrder(ctx context.Context, arg UpdateSellerOrderParams) (Order, error)
	UpdateStore(ctx context.Context, arg UpdateStoreParams) (Store, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUserReview(ctx context.Context, arg UpdateUserReviewParams) (Review, error)
	UpsertCartItem(ctx context.Context, arg UpsertCartItemParams) (CartItem, error)
}

var _ Querier = (*Queries)(nil)
