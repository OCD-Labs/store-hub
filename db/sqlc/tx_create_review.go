package db

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNoPurchase = errors.New("can't review what you didn't buy")
)

type CreateReviewTxParams struct {
	UserID  int64
	ItemID  int64
	StoreID int64
	OrderID int64
	Rating  string
	Comment string
}

// CreateReviewTx create a review for an item under a store, updates an order.
func (dbTx SQLTx) CreateReviewTx(ctx context.Context, arg CreateReviewTxParams) error {
	hasMade, err := dbTx.HasMadePurchase(ctx, HasMadePurchaseParams{
		CustomerID: arg.UserID,
		ItemID:     arg.ItemID,
		StoreID:    arg.StoreID,
	})
	if err != nil {
		return err
	}

	if !hasMade {
		return ErrNoPurchase
	}

	order, err := dbTx.GetOrderForBuyer(ctx, GetOrderForBuyerParams{
		OrderID: arg.OrderID,
		BuyerID: arg.UserID,
		StoreID: arg.StoreID,
	})
	if err != nil {
		return err
	}

	if order.IsReviewed {
		return nil
	}

	err = dbTx.execTx(ctx, func(q *Queries) error {
		err = q.CreateReview(ctx, CreateReviewParams{
			StoreID:            arg.StoreID,
			UserID:             arg.UserID,
			ItemID:             arg.ItemID,
			Rating:             arg.Rating,
			ReviewType:         "type:item_review",
			Comment:            arg.Comment,
			IsVerifiedPurchase: hasMade,
		})
		if err != nil {
			return err
		}

		_, err = q.UpdateBuyerOrder(ctx, UpdateBuyerOrderParams{
			OrderID: order.OrderID,
			BuyerID: arg.UserID,
			StoreID: arg.StoreID,
			IsReviewed: sql.NullBool{
				Bool:  hasMade,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}

		return err
	})

	return err
}
