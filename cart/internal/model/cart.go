package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Product_id  string             `json:"product_id" bson:"product_id"`
	Price       float64            `json:"price" bson:"price"`
	Quantity    int                `json:"quantity" bson:"quantity"`
	Total       float64            `json:"total" bson:"total"`
	CartItem_id string             `json:"cartItem_id" bson:"cartItem_id"`
}

type Cart struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	User_id     string             `json:"user_id" bson:"user_id"`
	Items       []CartItem         `json:"items" bson:"items"`
	TotalAmount float64            `json:"total_amount" bson:"total_amount"`
	Created_at  time.Time          `json:"created_at" bson:"created_at"`
	Updated_at  time.Time          `json:"updated_at" bson:"updated_at"`
	Cart_id     string             `json:"cart_id" bson:"cart_id"`
}
