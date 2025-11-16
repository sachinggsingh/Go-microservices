package grpc_handler

import (
	"context"
	"fmt"

	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
	proto "github.com/sachinggsingh/e-comm/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServer struct {
	proto.UnimplementedGetProductsServer
	database *db.Database
}

func NewProductServer(database *db.Database) *ProductServer {
	return &ProductServer{
		database: database,
	}
}

func (p *ProductServer) GetProducts(ctx context.Context, req *proto.GetProductRequest) (*proto.GetPRoductResponse, error) {
	// Validate request
	if req.ProductId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "product_id is required")
	}

	var product model.Product
	filter := bson.M{
		"product_id": req.ProductId,
	}
	err := p.database.ProductCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "product with id %s not found", req.ProductId)
		}
		fmt.Printf("Error fetching product from database: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to fetch product: %v", err)
	}

	// Handle nil pointers
	var name, description string
	var price float64
	if product.Name != nil {
		name = *product.Name
	}
	if product.Description != nil {
		description = *product.Description
	}
	if product.Price != nil {
		price = *product.Price
	}

	return &proto.GetPRoductResponse{
		Id:          product.ID.String(),
		Name:        name,
		Description: description,
		Price:       price,
	}, nil
}
