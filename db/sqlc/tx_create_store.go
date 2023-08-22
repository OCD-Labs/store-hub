package db

import (
	"context"
	"database/sql"

	"github.com/OCD-Labs/store-hub/util"
)

// A CreateStoreTxParams contains the input parameters for
// the create store transaction.
type CreateStoreTxParams struct {
	CreateStoreParams
	OwnerID int64 `json:"user_id"`
}

// A CreateStoreTxResult contains the result of the create store transaction.
type CreateStoreTxResult struct {
	Store  Store        `json:"store"`
	Owners []StoreOwner `json:"store_owners"`
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

		owner, err := q.AddCoOwnerAccess(ctx, AddCoOwnerAccessParams{
			UserID:       arg.OwnerID,
			StoreID:      result.Store.ID,
			AccessLevels: []int32{1},
			IsPrimary:    true,
		})
		if err != nil {
			return err
		}

		result.Owners = append(result.Owners, owner)

		arg := UpdateUserParams{
			ID: sql.NullInt64{
				Int64: arg.OwnerID,
				Valid: true,
			},
			Status: sql.NullString{
				String: util.STOREOWNER,
				Valid:  true,
			},
		}
		_, err = q.UpdateUser(ctx, arg)
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
