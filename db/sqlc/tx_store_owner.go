package db

import (
	"context"
	"database/sql"
)

// A AddCoOwnerAccessTxParams contains the input parameters of
// the add store access transaction.
type AddCoOwnerAccessTxParams struct {
	AccessLevel int32
	StoreID     int64
	InviteeID   int64
}

// AddCoOwnerAccessTx creates/update a user's access for a store.
func (store *SQLTx) AddCoOwnerAccessTx(ctx context.Context, arg AddCoOwnerAccessTxParams) (StoreOwner, error) {
	var coOwnerAccess StoreOwner

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// check access is granted to an existing user.
		_, err = q.GetUserAccessLevelsForStore(ctx, GetUserAccessLevelsForStoreParams{
			StoreID: arg.StoreID,
			UserID:  arg.InviteeID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				coOwnerAccess, err = q.AddCoOwnerAccess(ctx, AddCoOwnerAccessParams{
					AccessLevels: []int32{arg.AccessLevel},
					StoreID:      arg.StoreID,
					UserID:       arg.InviteeID,
					IsPrimary:    false,
				})
				if err != nil {
					return err
				}
			} else {
				coOwnerAccess, err = q.AddToCoOwnerAccess(ctx, AddToCoOwnerAccessParams{
					StoreID:        arg.StoreID,
					UserID:         arg.InviteeID,
					NewAccessLevel: arg.AccessLevel,
				})
				if err != nil {
					return err
				}
			}
		}

		return err
	})

	return coOwnerAccess, err
}

// A RevokeAccessTxParams contains the input parameters of
// the delete one or all user's store access transaction.
type RevokeAccessTxParams struct {
	AccountID           string
	StoreID             int64
	AccessLevelToRevoke int32
	DeleteAll           bool
}

// RevokeAccessTx deletes all user's access to a store.
func (store *SQLTx) RevokeAccessTx(ctx context.Context, arg RevokeAccessTxParams) ([]StoreOwner, error) {
	var accesses []StoreOwner

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// check access is granted to an existing user.
		user, err := q.GetUserByAccountID(ctx, arg.AccountID)
		if err != nil {
			return err
		}

		if arg.DeleteAll {
			err = q.RevokeAllAccess(ctx, RevokeAllAccessParams{
				UserID:  user.ID,
				StoreID: arg.StoreID,
			})
		} else {
			err = q.RevokeAccess(ctx, RevokeAccessParams{
				UserID:              user.ID,
				StoreID:             arg.StoreID,
				AccessLevelToRevoke: arg.AccessLevelToRevoke,
			})
		}

		if err != nil {
			return err
		}

		accesses, err = q.GetStoreOwnersByStoreID(ctx, arg.StoreID)
		if err != nil {
			return err
		}

		return err
	})

	return accesses, err
}
