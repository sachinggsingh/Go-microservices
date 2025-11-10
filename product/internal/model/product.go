package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        *string            `json:"name" validate:"required"`
	Description *string            `json:"description" validate:"required"`
	Price       *float64           `json:"price" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Product_id  string             `json:"product_id"`
}
