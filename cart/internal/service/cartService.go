package service

import (
	"errors"
	"time"

	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService struct {
	cartRepo repository.CartRepository
}

func NewCartService(cartRepo repository.CartRepository) CartService {
	return CartService{
		cartRepo: cartRepo,
	}
}

func (c *CartService) CreateCart(cart *model.Cart) (*model.Cart, error) {
	now := time.Now().UTC()

	cart.ID = primitive.NewObjectID()
	cart.Created_at = now
	cart.Updated_at = now
	cart.Cart_id = cart.ID.Hex()

	if cart.User_id == "" {
		return nil, errors.New("user_id is required to create cart")
	}
	createdCart, err := c.cartRepo.CreateCart(cart)
	if err != nil {
		return nil, err
	}
	return createdCart, nil
}

func (c *CartService) FindCartByUserID(userID string, items *model.CartItem) (*model.Cart, error) {
	return c.cartRepo.FindCartByUserID(userID, items)
}
