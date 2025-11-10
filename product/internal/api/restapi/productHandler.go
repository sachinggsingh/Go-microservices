package restapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/service"
)

type ProductHandler struct {
	productService *service.Productservice
}

func NewProductHandler(productService *service.Productservice) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}
func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	pro, err := ph.productService.CreateProduct(&product)
	if err != nil {
		if errors.Is(err, service.CantCreateProduct) {
			http.Error(w, "Missing required product fields", http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ProductExist) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pro)
}

func (ph *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	pro, err := ph.productService.GetAllProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pro)
}

func (ph *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	productId := vars["product_id"]
	if productId == "" {
		http.Error(w, "Product id is required", http.StatusBadRequest)
		return
	}

	pro, err := ph.productService.GetProductById(&model.Product{Product_id: productId})
	log.Println(pro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if pro == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pro)
}
