package db

import "context"

type GetUserCartResult struct {
	Cart []GetCartByUserIDRow `json:"cart"`
	CartID int64 `json:"cart_id"`
}

// GetUserCart retrieves a user's cart items.
func (dbTx *SQLTx) GetUserCartTx(ctx context.Context, userID int64 ) (GetUserCartResult, error) {
	var result GetUserCartResult

	err := dbTx.execTx(ctx, func(q *Queries) error {
		var err error

		result.Cart, err = q.GetCartByUserID(ctx, userID)
		if err != nil {
			return err
		}
		result.CartID, err = q.GetCartID(ctx, userID)

		return err
	})

	return result, err
}
