package restapi

import (
	"net/http"

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

}

func (c *CartHandler) GetAllCarts(w http.ResponseWriter, r *http.Request) {

}

func (c *CartHandler) GetCartById(w http.ResponseWriter, r *http.Request) {

}

func (c *CartHandler) UpdateCart(w http.ResponseWriter, r *http.Request) {

}

func (c *CartHandler) DeleteCart(w http.ResponseWriter, r *http.Request) {

}
