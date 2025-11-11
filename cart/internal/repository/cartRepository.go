package repository

import (
	"context"
	"log"
	"time"

	"github.com/sachinggsingh/e-comm/internal/errors"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository interface {
	CreateCart(cart *model.Cart) (*model.Cart, error)
	FindCartByUserID(userID string) (*model.Cart, error)
	UpdateCart(cart *model.Cart) (*model.Cart, error)
	DeleteCart(userID string) error
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.db.CartCollection.InsertOne(ctx, cart)
	if err != nil {
		log.Printf("Error creating cart: %v", err)
		return nil, err
	}
	return cart, nil
}

func (c *cartRepository) FindCartByUserID(userID string) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cart model.Cart

	filter := bson.M{
		"user_id": userID,
	}
	err := c.db.CartCollection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Cart not found for user_id: %s", userID)
			return nil, errors.ErrCartNotFound
		}
		log.Printf("Error finding cart: %v", err)
		return nil, err
	}
	return &cart, nil
}

func (c *cartRepository) UpdateCart(cart *model.Cart) (*model.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": cart.User_id,
	}

	update := bson.M{
		"$set": bson.M{
			"items":        cart.Items,
			"total_amount": cart.TotalAmount,
			"updated_at":   cart.Updated_at,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedCart model.Cart
	err := c.db.CartCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Cart not found for user_id: %s", cart.User_id)
			return nil, errors.ErrCartNotFound
		}
		log.Printf("Error updating cart: %v", err)
		return nil, err
	}
	return &updatedCart, nil
}

func (c *cartRepository) DeleteCart(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
	}

	result, err := c.db.CartCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting cart: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("Cart not found for user_id: %s", userID)
		return errors.ErrCartNotFound
	}

	return nil
}
