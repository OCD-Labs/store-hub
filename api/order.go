package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/OCD-Labs/store-hub/pagination"
	"github.com/OCD-Labs/store-hub/util"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type createOrderRequestBody struct {
	ItemID         int64  `json:"item_id" validate:"required,min=1"`
	OrderQuantity  int32  `json:"order_quantity" validate:"required,min=1"`
	SellerID       int64  `json:"seller_id" validate:"required,min=1"`
	StoreID        int64  `json:"store_id" validate:"required,min=1"`
	DeliveryFee    string `json:"delivery_fee" validate:"required"`
	PaymentChannel string `json:"payment_channel" validate:"required,oneof=NEAR 'Debit Card' PayPal 'Credit Card'"`
	PaymentMethod  string `json:"payment_method" validate:"required,oneof='Instant Pay' 'Pay on Delivery'"`
}

// createOrder maps to endpoint "POST /seller/orders"
func (s *StoreHub) createOrder(w http.ResponseWriter, r *http.Request) {
	var reqBody createOrderRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	supply_quantity, err := s.dbStore.CheckItemStoreMatch(r.Context(), db.CheckItemStoreMatchParams{
		ItemID:  reqBody.ItemID,
		StoreID: reqBody.StoreID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusForbidden, "failed to create orders")
			log.Error().Err(err).Msg("error occurred")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to create order")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	if supply_quantity < int64(reqBody.OrderQuantity) {
		s.errorResponse(w, r, http.StatusForbidden, "failed to create orders")
		err = fmt.Errorf("order quantity %d is greater than supply %d", reqBody.OrderQuantity, supply_quantity)
		log.Error().Err(err).Msg("error occurred")
		return
	}

	_, err = s.dbStore.IsStoreOwner(r.Context(), db.IsStoreOwnerParams{
		UserID:  reqBody.SellerID,
		StoreID: reqBody.StoreID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			s.errorResponse(w, r, http.StatusNotFound, "user is not the owner of the store")
			log.Error().Err(err).Msg("error occurred")
			return
		}

		s.errorResponse(w, r, http.StatusInternalServerError, "failed to create order")
		return
	}

	authPayload := s.contextGetToken(r)

	order, err := s.dbStore.CreateOrder(r.Context(), db.CreateOrderParams{
		ItemID:         reqBody.ItemID,
		OrderQuantity:  reqBody.OrderQuantity,
		BuyerID:        authPayload.UserID,
		SellerID:       reqBody.SellerID,
		StoreID:        reqBody.StoreID,
		DeliveryFee:    reqBody.DeliveryFee,
		PaymentChannel: reqBody.PaymentChannel,
		PaymentMethod:  reqBody.PaymentMethod,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				s.errorResponse(w, r, http.StatusConflict, "a referenced key was not found")
			default:
				s.errorResponse(w, r, http.StatusInternalServerError, "failed to create order")
			}
		} else {
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to create order")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "created a new order",
			"result": envelop{
				"order": order,
			},
		},
	}, nil)
}

type listSellerOrdersQueryStr struct {
	ItemName       string    `querystr:"item_name"`
	CreatedAtStart time.Time `querystr:"created_at_start"`
	CreatedAtEnd   time.Time `querystr:"created_at_end"`
	PaymentChannel string    `querystr:"payment_channel"`
	DeliveryStatus string    `querystr:"delivery_status"`
	Page           int       `querystr:"page" validate:"max=10000000"`
	PageSize       int       `querystr:"page_size" validate:"max=20"`
	Sort           string    `querystr:"sort"`
}

// listSellerOrders maps to endpoint "GET /seller/orders"
func (s *StoreHub) listSellerOrders(w http.ResponseWriter, r *http.Request) {
	var reqQueryStr listSellerOrdersQueryStr
	if err := s.shouldBindQuery(w, r, &reqQueryStr); err != nil {
		return
	}

	if reqQueryStr.Page < 1 {
		reqQueryStr.Page = 1
	}
	if reqQueryStr.PageSize < 1 {
		reqQueryStr.PageSize = 15
	}
	if reqQueryStr.Sort == "" {
		reqQueryStr.Sort = "-id"
	}

	authPayload := s.contextGetToken(r)

	// TODO: Add a proper range logic for createdAt search params
	arg := db.ListSellerOrdersParams{
		ItemName:       reqQueryStr.ItemName,
		SellerID:       authPayload.UserID,
		CreatedAtStart: reqQueryStr.CreatedAtStart,
		CreatedAtEnd:   reqQueryStr.CreatedAtEnd,
		PaymentChannel: reqQueryStr.PaymentChannel,
		DeliveryStatus: reqQueryStr.DeliveryStatus,
		Filters: pagination.Filters{
			Page:         reqQueryStr.Page,
			PageSize:     reqQueryStr.PageSize,
			Sort:         reqQueryStr.Sort,
			SortSafelist: []string{"-id", "-item_name", "-created_at", "order_id", "item_name", "created_at"},
		},
	}
	orders, pagination, err := s.dbStore.ListSellerOrders(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve orders")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found some orders",
			"result": envelop{
				"order":    orders,
				"metadata": pagination,
			},
		},
	}, nil)
}

type getSellerOrderPathVars struct {
	OrderID int64 `path:"order_id" validate:"required,min=1"`
}

// getSellerOrder maps to endpoint "GET /seller/orders/{order_id}"
func (s *StoreHub) getSellerOrder(w http.ResponseWriter, r *http.Request) {
	var pathVar getSellerOrderPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	authPayload := s.contextGetToken(r)

	order, err := s.dbStore.GetOrderForSeller(r.Context(), db.GetOrderForSellerParams{
		SellerID: authPayload.UserID,
		OrderID:  pathVar.OrderID,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "order not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve order details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "found an order",
			"result": envelop{
				"order": order,
			},
		},
	}, nil)
}

type updateSellerOrderRequest struct {
	DeliveredOn          time.Time `json:"delivered_on"`
	DeliveryStatus       *string   `json:"delivery_status"`
	ExpectedDeliveryDate time.Time `json:"expected_delivery_date"`
}

type updateSellerOrderPathVars struct {
	OrderID int64 `path:"order_id" validate:"required,min=1"`
}

// getSellerOrder maps to endpoint "PATCH /seller/orders/{order_id}"
func (s *StoreHub) updateSellerOrder(w http.ResponseWriter, r *http.Request) {
	var pathVars updateSellerOrderPathVars
	if err := s.ShouldBindPathVars(w, r, &pathVars); err != nil {
		return
	}

	var reqBody updateSellerOrderRequest
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	authPayload := s.contextGetToken(r)

	arg := db.UpdateSellerOrderParams{
		OrderID:  pathVars.OrderID,
		SellerID: authPayload.UserID,
	}

	if reqBody.DeliveryStatus != nil && *reqBody.DeliveryStatus != "" {
		if !util.IsValidStatus(*reqBody.DeliveryStatus) {
			s.errorResponse(w, r, http.StatusForbidden, "unsupported status")
			return
		}
		arg.DeliveryStatus = sql.NullString{
			String: *reqBody.DeliveryStatus,
			Valid:  true,
		}
	}

	if reqBody.DeliveryStatus != nil && *reqBody.DeliveryStatus == "DELIVERED" {
		if reqBody.DeliveredOn.IsZero() {
			s.errorResponse(w, r, http.StatusBadRequest, "can't change order status to DELIVERED without its 'delivered_on' date")
			return
		}
		arg.DeliveredOn = sql.NullTime{
			Time:  reqBody.DeliveredOn,
			Valid: true,
		}
	} else if !reqBody.DeliveredOn.IsZero() {
		// If DeliveredOn is set but DeliveryStatus is not "DELIVERED"
		s.errorResponse(w, r, http.StatusBadRequest, "can't set 'delivered_on' date if order status is not DELIVERED")
		return
	}

	if !reqBody.ExpectedDeliveryDate.IsZero() {
		arg.ExpectedDeliveryDate = sql.NullTime{
			Time:  reqBody.ExpectedDeliveryDate,
			Valid: true,
		}
	}

	updatedOrder, err := s.dbStore.UpdateSellerOrderTx(r.Context(), arg) // TODO: Ensure that the DeliveredOn date is not before created_at value of an order.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.errorResponse(w, r, http.StatusNotFound, "order not found")
		default:
			s.errorResponse(w, r, http.StatusInternalServerError, "failed to update order details")
		}
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "updated order's details",
			"result": envelop{
				"order": updatedOrder,
			},
		},
	}, nil)
}
