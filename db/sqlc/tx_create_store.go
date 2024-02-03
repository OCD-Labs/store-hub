package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/OCD-Labs/store-hub/util"
)

// A CreateStoreTxParams contains the input parameters for
// the create store transaction.
type CreateStoreTxParams struct {
	CreateStoreParams
	OwnerID int64 `json:"user_id"`
	AfterCreate func(context.Context, Store) error
}

type StoreOwnerDetails struct {
	AccountID       string    `json:"account_id"`
	ProfileImgURL   string    `json:"profile_img_url"`
	AccessLevels    []int32   `json:"access_levels"`
	IsOriginalOwner bool      `json:"is_original_owner"`
	AddedAt         time.Time `json:"added_at"`
}

// A CreateStoreTxResult contains the result of the create store transaction.
type CreateStoreTxResult struct {
	Store       Store               `json:"store"`
	StoreOwners []StoreOwnerDetails `json:"store_owners"`
}

// CreateStoreTx creates a store and its ownership data.
func (dbTx *SQLTx) CreateStoreTx(ctx context.Context, arg CreateStoreTxParams) (CreateStoreTxResult, error) {
	var result CreateStoreTxResult

	err := dbTx.execTx(ctx, func(q *Queries) error {
		// Create the store
		store, err := q.CreateStore(ctx, arg.CreateStoreParams)
		if err != nil {
			return err
		}
		result.Store = store

		err = arg.AfterCreate(ctx, store)
		if err != nil {
			return err
		}

		// Add the owner to the store
		_, err = q.AddCoOwnerAccess(ctx, AddCoOwnerAccessParams{
			UserID:       arg.OwnerID,
			StoreID:      store.ID,
			AccessLevels: []int32{1},
			IsPrimary:    true,
		})
		if err != nil {
			return err
		}

		// Update the user's status to STOREOWNER
		argUpdate := UpdateUserParams{
			ID: sql.NullInt64{
				Int64: arg.OwnerID,
				Valid: true,
			},
			Status: sql.NullString{
				String: util.STOREOWNER,
				Valid:  true,
			},
		}
		_, err = q.UpdateUser(ctx, argUpdate)
		if err != nil {
			return err
		}

		// Retrieve the store and its owners
		storeWithOwners, err := q.GetStoreByID(ctx, store.ID)
		if err != nil {
			return err
		}

		// Unmarshal the store owners data
		var owners []StoreOwnerDetails
		if err := json.Unmarshal(storeWithOwners.StoreOwners, &owners); err != nil {
			return err
		}
		result.StoreOwners = owners

		return nil
	})

	return result, err
}
