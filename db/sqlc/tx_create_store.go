package db

import (
	"context"
)

// A CreateStoreTxParams contains the input parameters for
// the create store transaction.
type CreateStoreTxParams struct {
	CreateStoreParams
	OwnerID     int64 `json:"user_id"`
	AccessLevel int16 `json:"access_level"`
}

// A CreateStoreTxResult contains the result of the create store transaction.
type CreateStoreTxResult struct {
	Store  Store
	Owners []StoreOwner
}

// CreateStoreTx creates a store and its ownership data.
func (store *SQLTx) CreateStoreTx(ctx context.Context, arg CreateStoreTxParams) (CreateStoreTxResult, error) {
	var result CreateStoreTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Store, err = q.CreateStore(ctx, arg.CreateStoreParams)
		if err != nil {
			return err
		}

		owner, err := q.CreateStoreOwner(ctx, CreateStoreOwnerParams{
			UserID:      arg.OwnerID,
			StoreID:     result.Store.ID,
			AccessLevel: arg.AccessLevel,
		})
		if err != nil {
			return err
		}

		result.Owners = append(result.Owners, owner)

		return nil
	})

	return result, err
}