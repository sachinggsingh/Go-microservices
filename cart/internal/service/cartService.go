package service

import "github.com/sachinggsingh/e-comm/internal/repository"

type CartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return CartService{
		cartRepo: cartRepo,
	}
}
