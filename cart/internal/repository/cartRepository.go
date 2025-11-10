package repository

import (
	"context"
	"time"

	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	CreateCart(cart *model.Cart) (*model.Cart, error)
	FindCartByUserID(userID string, items *model.CartItem) (*model.Cart, error)
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := c.db.CartCollection.InsertOne(ctx, cart)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
func (c *cartRepository) FindCartByUserID(userID string, items *model.CartItem) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var cart model.Cart

	filter := bson.M{
		"user_id": userID,
	}
	err := c.db.CartCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &cart, nil
}
func (c *cartRepository) UpdateCart(cart *model.Cart) (*model.Cart, error) {
	return nil, nil
}

func (c *cartRepository) DeleteCart(cart *model.Cart) error {
	return nil
}
