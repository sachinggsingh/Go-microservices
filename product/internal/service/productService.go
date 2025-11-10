package service

import (
	"errors"
	"time"

	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Productservice struct {
	proRepo repository.ProductRepository
}

var (
	CantCreateProduct = errors.New("cannot create product")
	ProductExist      = errors.New("product already exist")
)

func NewProductService(proRepo repository.ProductRepository) *Productservice {
	return &Productservice{
		proRepo: proRepo,
	}
}

func (p *Productservice) CreateProduct(pro *model.Product) (*model.Product, error) {
	if pro.Name == nil || pro.Description == nil || pro.Price == nil {
		return nil, CantCreateProduct
	}

	// check if the product already exist or not
	exist, err := p.proRepo.CheckIfPoductExist(pro)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ProductExist
	}
	now := time.Now()
	pro.ID = primitive.NewObjectID()
	pro.Created_at = now
	pro.Updated_at = now
	pro.Product_id = pro.ID.Hex()

	return p.proRepo.CreateProduct(pro)
}

func (p *Productservice) GetAllProducts() ([]*model.Product, error) {
	return p.proRepo.GetAllProducts()
}

func (p *Productservice) GetProductById(pro *model.Product) (*model.Product, error) {
	return p.proRepo.GetProductById(pro)
}

// func (p *Productservice) UpdateProduct(pro *model.Product) (*model.Product, error) {
// 	return p.proRepo.UpdateProduct(pro)
// }

// func (p *Productservice) DeleteProduct(id string) error {
// 	return p.proRepo.DeleteProduct(id)
// }
