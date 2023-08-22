package db

import (
	"context"

	"github.com/OCD-Labs/store-hub/util"
)

// UpdateSellerOrderTx updates a order row, create a sale row if order is DELIVERED.
func (dbTx SQLTx) UpdateSellerOrderTx(ctx context.Context, arg UpdateSellerOrderParams) (GetOrderForSellerRow, error) {
	var sellerOrder GetOrderForSellerRow
	var err error

	sellerOrder, err = dbTx.GetOrderForSeller(ctx, GetOrderForSellerParams{
		OrderID:  arg.OrderID,
		SellerID: arg.SellerID,
		StoreID: arg.StoreID,
	})
	if err != nil {
		return sellerOrder, err
	}

	if sellerOrder.DeliveryStatus != arg.DeliveryStatus.String && util.CanChangeStatus(sellerOrder.DeliveryStatus, arg.DeliveryStatus.String) {

		err = dbTx.execTx(ctx, func(q *Queries) error {
			var o Order

			o, err = dbTx.UpdateSellerOrder(ctx, arg)
			if err != nil {
				return err
			}

			if o.DeliveryStatus == "DELIVERED" {
				sArg := CreateSaleParams{
					StoreID:    o.StoreID,
					ItemID:     o.ItemID,
					CustomerID: o.BuyerID,
					SellerID:   o.SellerID,
					OrderID:    o.ID,
				}

				_, err = dbTx.CreateSale(ctx, sArg)
				if err != nil {
					return err
				}
			}

			if o.DeliveryStatus == "RETURNED" {
				err = dbTx.ReduceSalesOverview(ctx, ReduceSalesOverviewParams{
					StoreID: o.StoreID,
					ItemID:  o.ItemID,
					OrderID: o.ID,
				})
				if err != nil {
					return err
				}
			}

			sellerOrder.DeliveredOn = o.DeliveredOn
			sellerOrder.ExpectedDeliveryDate = o.ExpectedDeliveryDate
			sellerOrder.DeliveryStatus = o.DeliveryStatus

			return nil
		})
	} else {
		return sellerOrder, nil
	}

	return sellerOrder, err
}
