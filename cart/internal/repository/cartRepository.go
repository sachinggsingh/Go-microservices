package repository

import (
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
)

type CartRepository interface {
	CreateCart(cart *model.Cart) (*model.Cart, error)
	GetAllCarts() ([]*model.Cart, error)
	GetCartById(cart *model.Cart) (*model.Cart, error)
	UpdateCart(cart *model.Cart) (*model.Cart, error)
	DeleteCart(cart *model.Cart) error
}

type cartRepository struct {
	db *db.Database
}

func NewCartRepository(db *db.Database) CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (c *cartRepository) CreateCart(cart *model.Cart) (*model.Cart, error) {
	return nil, nil
}

func (c *cartRepository) GetAllCarts() ([]*model.Cart, error) {
	return nil, nil
}

func (c *cartRepository) GetCartById(cart *model.Cart) (*model.Cart, error) {
	return nil, nil
}

func (c *cartRepository) UpdateCart(cart *model.Cart) (*model.Cart, error) {
	return nil, nil
}

func (c *cartRepository) DeleteCart(cart *model.Cart) error {
	return nil
}
