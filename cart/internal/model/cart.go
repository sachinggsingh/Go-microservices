package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ID         primitive.ObjectID ` bson:"_id"`
	Product_id string             `json:"product_id" `
	Price      float64            `json:"price" `
	Quantity   int                `json:"quantity"`
	Total      float64            `json:"total" `
	Cart_id    string             `json:"cart_id" `
}

type Cart struct {
	ID          primitive.ObjectID `bson:"_id"`
	User_id     string             `json:"user_id"`
	Cart_id     string             `json:"cart_id"`
	Items       []CartItem         `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
