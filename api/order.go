package api

import (
	"fmt"
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type createOrderRequestBody struct {
	ItemID int64 `json:"item_id" validate:"required,min=1"`
	OrderQuantity  int32  `json:"order_quantity" validate:"required,min=1"`
	BuyerID        int64  `json:"buyer_id" validate:"required,min=1"`
	SellerID       int64  `json:"seller_id" validate:"required,min=1"`
	StoreID        int64  `json:"store_id" validate:"required,min=1"`
	DeliveryFee    string `json:"delivery_fee" validate:"required"`
	PaymentChannel string `json:"payment_channel" validate:"required"`
	PaymentMethod  string `json:"payment_method" validate:"required"`
}

func (s *StoreHub) createOrder(w http.ResponseWriter, r *http.Request) {
	var reqBody createOrderRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	store, err := s.dbStore.CreateOrder(r.Context(), db.CreateOrderParams{
		ItemID: reqBody.ItemID,
		OrderQuantity: reqBody.OrderQuantity,
		BuyerID: reqBody.BuyerID,
		SellerID: reqBody.SellerID,
		StoreID: reqBody.StoreID,
		DeliveryFee: reqBody.DeliveryFee,
		PaymentChannel: reqBody.PaymentChannel,
		PaymentMethod: reqBody.PaymentMethod,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				s.errorResponse(w, r, http.StatusConflict, "a referenced key was not found")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item")
			}
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	fmt.Printf("store: %v\n", store)
}

func (s *StoreHub) listSellerOrders(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) getSellerOrder(w http.ResponseWriter, r *http.Request) {

}

func (s *StoreHub) updateSellerOrder(w http.ResponseWriter, r *http.Request) {

}
