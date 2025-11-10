package repository

import (
	"context"
	"time"

	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	CreateProduct(p *model.Product) (*model.Product, error)
	GetAllProducts() ([]*model.Product, error)
	GetProductById(pro *model.Product) (*model.Product, error)
	// UpdateProduct(p *model.Product) (*model.Product, error)
	// DeleteProduct(id string) error
	CheckIfPoductExist(pro *model.Product) (bool, error)
}

type productRepo struct {
	productColl *mongo.Collection
}

func NewProductRepository(database *db.Database) ProductRepository {
	return &productRepo{
		productColl: database.ProductCollection,
	}
}

func (p *productRepo) CreateProduct(pro *model.Product) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := p.productColl.InsertOne(ctx, pro)
	if err != nil {
		return nil, err
	}
	return pro, nil
}

func (p *productRepo) GetAllProducts() ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cursor, err := p.productColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	for cursor.Next(ctx) {
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productRepo) GetProductById(pro *model.Product) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var product model.Product
	filter := bson.M{
		"product_id": pro.Product_id,
	}
	err := p.productColl.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &product, nil
}

// func (p *productRepo) UpdateProduct(pro *model.Product) (*model.Product, error) {
// 	return nil, nil
// }

// func (p *productRepo) DeleteProduct(id string) error {
// 	return nil
// }

func (p *productRepo) CheckIfPoductExist(pro *model.Product) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filter := bson.M{
		"product_id": pro.Product_id,
	}
	var porduct model.Product

	err := p.productColl.FindOne(ctx, filter).Decode(&porduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
