package db

import "context"

// UpdateOrderTx updates a order row, create a sale row if order is DELIVERED.
func (dbTx SQLTx) UpdateOrderTx(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	var order Order

	err := dbTx.execTx(ctx, func(q *Queries) error {
		var err error

		order, err = dbTx.UpdateOrder(ctx, arg)
		if err != nil {
			return err
		}

		if order.DeliveryStatus == "DELIVERED" {
			exist, err := dbTx.SaleExists(ctx, order.ID)
			if err != nil {
				return err
			}

			if !exist {
				sArg := CreateSaleParams{
					StoreID:    order.StoreID,
					ItemID:     order.ItemID,
					CustomerID: order.BuyerID,
					SellerID:   order.SellerID,
					OrderID:    order.ID,
				}
	
				_, err = dbTx.CreateSale(ctx, sArg)
				if err != nil {
					return err
				}
			}
		}

		if order.DeliveryStatus == "RETURNED" {
			err = dbTx.ReduceSaleCount(ctx, ReduceSaleCountParams{
				StoreID: order.StoreID,
				ItemID: order.ItemID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return order, err
}
