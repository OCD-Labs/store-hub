package api

import (
	"net/http"

	db "github.com/OCD-Labs/store-hub/db/sqlc"
	"github.com/rs/zerolog/log"
)

type getUserCartPathVar struct {
	UserID int64 `path:"user_id" validate:"required,min=1"`
}

// getUserCart maps to endpoint "GET /carts/{user_id}"
func (s *StoreHub) getUserCart(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar getUserCartPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// db query
	cart, err := s.dbStore.GetCartByUserID(r.Context(), pathVar.UserID)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to retrieve cart")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "retrieved user cart",
			"result": envelop{
				"cart": cart,
			},
		},
	}, nil)
}

type addItemToCartRequestBody struct {
	ItemID  int64 `json:"item_id" validate:"required,min=1"`
	StoreID int64 `json:"store_id" validate:"required,min=1"`
}

type addItemToCartPathVar struct {
	CartID int64 `path:"cart_id" validate:"required,min=1"`
}

// addItemToCart maps to endpoint "POST /carts/{cart_id}/items"
func (s *StoreHub) addItemToCart(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar addItemToCartPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody addItemToCartRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	// db query
	arg := db.UpsertCartItemParams{
		CartID:  pathVar.CartID,
		ItemID:  reqBody.ItemID,
		StoreID: reqBody.StoreID,
	}
	cartItem, err := s.dbStore.UpsertCartItem(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to add item to cart")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusCreated, envelop{
		"status": "success",
		"data": envelop{
			"message": "item added to cart",
			"result": envelop{
				"cart_item": cartItem,
			},
		},
	}, nil)
}

type removeItemFromCartPathVar struct {
	CartID int64 `path:"cart_id" validate:"required,min=1"`
	ItemID int64 `path:"item_id" validate:"required,min=1"`
}

// removeItemFromCart maps to endpoint "DELETE /carts/{cart_id}/items/{item_id}"
func (s *StoreHub) removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar removeItemFromCartPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// db query
	err := s.dbStore.RemoveItemFromCart(r.Context(), db.RemoveItemFromCartParams{
		CartID: pathVar.CartID,
		ItemID: pathVar.ItemID,
	})
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to remove item from cart")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "item removed from cart",
		},
	}, nil)
}

type increaseCartItemRequestBody struct {
	IncreaseAmount int32 `json:"increase_amount" validate:"required,min=1"`
}

type increaseCartItemPathVar struct {
	CartID int64 `path:"cart_id" validate:"required,min=1"`
	ItemID int64 `path:"item_id" validate:"required,min=1"`
}

// increaseCartItem maps to endpoint "PUT /carts/{cart_id}/items/{item_id}/increase"
func (s *StoreHub) increaseCartItem(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar increaseCartItemPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody increaseCartItemRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	// db query
	arg := db.IncreaseCartItemQuantityParams{
		CartID:         pathVar.CartID,
		ItemID:         pathVar.ItemID,
		IncreaseAmount: reqBody.IncreaseAmount,
	}
	cartItem, err := s.dbStore.IncreaseCartItemQuantity(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to increase cart item quantity")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "cart item quantity increased",
			"result": envelop{
				"updated_cart_item": cartItem,
			},
		},
	}, nil)
}

type decreaseCartItemRequestBody struct {
	DecreaseAmount int32 `json:"decrease_amount" validate:"required,min=1"`
}

type decreaseCartItemPathVar struct {
	CartID int64 `path:"cart_id" validate:"required,min=1"`
	ItemID int64 `path:"item_id" validate:"required,min=1"`
}

// decreaseCartItem maps to endpoint "PUT /carts/{cart_id}/items/{item_id}/decrease"
func (s *StoreHub) decreaseCartItem(w http.ResponseWriter, r *http.Request) {
	// parse path variables
	var pathVar decreaseCartItemPathVar
	if err := s.ShouldBindPathVars(w, r, &pathVar); err != nil {
		return
	}

	// parse request body
	var reqBody decreaseCartItemRequestBody
	if err := s.shouldBindBody(w, r, &reqBody); err != nil {
		return
	}

	// db query
	arg := db.DecreaseCartItemQuantityParams{
		CartID:         pathVar.CartID,
		ItemID:         pathVar.ItemID,
		DecreaseAmount: reqBody.DecreaseAmount,
	}
	cartItem, err := s.dbStore.DecreaseCartItemQuantity(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, r, http.StatusInternalServerError, "failed to decrease cart item quantity")
		log.Error().Err(err).Msg("error occurred")
		return
	}

	// return response
	s.writeJSON(w, http.StatusOK, envelop{
		"status": "success",
		"data": envelop{
			"message": "cart item quantity decreased",
			"result": envelop{
				"updated_cart_item": cartItem,
			},
		},
	}, nil)
}
