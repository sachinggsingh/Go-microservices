package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/service"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (c *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var newCart model.Cart
	if err := json.NewDecoder(r.Body).Decode(&newCart); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	cart, err := c.cartService.CreateCart(&newCart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func (c *CartHandler) FindCartByUserID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "user_id claim missing or invalid", http.StatusUnauthorized)
		return
	}

	cart, err := c.cartService.FindCartByUserID(userID, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (c *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {

}

func (c *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {

}
