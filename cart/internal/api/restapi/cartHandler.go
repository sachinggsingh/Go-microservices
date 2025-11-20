package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	carterrors "github.com/sachinggsingh/e-comm/internal/errors"
	"github.com/sachinggsingh/e-comm/internal/middleware"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/pkg/payment"
	"github.com/sachinggsingh/e-comm/internal/service"
)

type CartHandler struct {
	cartService   *service.CartService
	paymentClient payment.PaymentClient
}

type CartItemRequest struct {
	Product_id string  `json:"product_id"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
}

type CreateCartRequest struct {
	Items []CartItemRequest `json:"items"`
}

type UpdateCartRequest struct {
	Items []CartItemRequest `json:"items"`
}

func NewCartHandler(cartService *service.CartService, paymentClient payment.PaymentClient) *CartHandler {
	return &CartHandler{
		cartService:   cartService,
		paymentClient: paymentClient,
	}
}

func (c *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	fmt.Println(userID)
	fmt.Println(ok, "ok")
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	var req CreateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert request items to model items
	items := make([]model.CartItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, model.CartItem{
			Product_id: item.Product_id,
			Price:      item.Price,
			Quantity:   item.Quantity,
		})
	}

	cart, err := c.cartService.CreateCart(userID, items)
	if err != nil {
		if errors.Is(err, carterrors.ErrInvalidUserID) ||
			errors.Is(err, carterrors.ErrEmptyCart) ||
			errors.Is(err, carterrors.ErrInvalidItem) ||
			errors.Is(err, carterrors.ErrInvalidQuantity) ||
			errors.Is(err, carterrors.ErrInvalidPrice) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func (c *CartHandler) FindCartByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	cart, err := c.cartService.FindCartByUserID(userID)
	if err != nil {
		if errors.Is(err, carterrors.ErrCartNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, carterrors.ErrInvalidUserID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (c *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	var req UpdateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert request items to model items
	items := make([]model.CartItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, model.CartItem{
			Product_id: item.Product_id,
			Price:      item.Price,
			Quantity:   item.Quantity,
		})
	}

	cart, err := c.cartService.UpdateCart(userID, items)
	if err != nil {
		if errors.Is(err, carterrors.ErrCartNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, carterrors.ErrInvalidUserID) ||
			errors.Is(err, carterrors.ErrEmptyCart) ||
			errors.Is(err, carterrors.ErrInvalidItem) ||
			errors.Is(err, carterrors.ErrInvalidQuantity) ||
			errors.Is(err, carterrors.ErrInvalidPrice) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (c *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	err := c.cartService.DeleteCart(userID)
	if err != nil {
		if errors.Is(err, carterrors.ErrCartNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, carterrors.ErrInvalidUserID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (c *CartHandler) CheckoutCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	cart, err := c.cartService.FindCartByUserID(userID)
	if err != nil {
		if errors.Is(err, carterrors.ErrCartNotFound) {
			http.Error(w, "cart not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, carterrors.ErrInvalidUserID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(cart.Items) == 0 {
		http.Error(w, "cart is empty", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	paymentItems, err := c.cartService.PreparePaymentItems(ctx, cart)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to prepare payment items: %v", err), http.StatusInternalServerError)
		return
	}

	session, err := c.paymentClient.CreatePayment(paymentItems, userID, cart.Cart_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create payment session: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"checkout_url": session.URL,
		"session_id":   session.ID,
		"cart_id":      cart.Cart_id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
