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
	proto.UnimplementedShowProductServer
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

func (p *ProductServer) ShowProduct(req *proto.ShowProductRequest, stream proto.ShowProduct_ShowProductServer) error {
	// Validate request
	if req.ProductId == "" {
		return status.Errorf(codes.InvalidArgument, "product_id is required")
	}

	ctx := stream.Context()
	cursor, err := p.database.ProductCollection.Find(ctx, bson.M{
		"product_id": req.ProductId,
	})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to query products: %v", err)
	}
	defer cursor.Close(ctx)

	found := false
	for cursor.Next(ctx) {
		found = true
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			return err
		}
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
		res := &proto.ShowProductResponse{
			Id:          product.ID.String(),
			Name:        name,
			Description: description,
			Price:       price,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	if !found {
		return status.Errorf(codes.NotFound, "product with id %s not found", req.ProductId)
	}
	return nil
}
