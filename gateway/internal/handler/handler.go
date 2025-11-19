package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	proto "github.com/sachinggsingh/e-comm/pb"
)

type ProductHandler struct {
	auth        proto.ValidateTokenClient
	product     proto.GetProductsClient
	showProduct proto.ShowProductClient
}

func NewProductHandler(auth proto.ValidateTokenClient, product proto.GetProductsClient, showProduct proto.ShowProductClient) *ProductHandler {
	return &ProductHandler{auth, product, showProduct}
}

func (h *ProductHandler) GetProductGateway(w http.ResponseWriter, r *http.Request) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		http.Error(w, "missing or malformed bearer token", http.StatusUnauthorized)
		return
	}

	res, err := h.auth.ValidateToken(r.Context(), &proto.ValidateTokenRequest{Token: strings.TrimSpace(parts[1])})
	if err != nil {
		http.Error(w, "token validation failed", http.StatusUnauthorized)
		fmt.Println()
		return
	}
	if !res.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	productID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/gateway/product/"), "/")
	if productID == "" {
		productID = r.URL.Query().Get("product_id")
	}
	if productID == "" {
		http.Error(w, "missing product_id", http.StatusBadRequest)
		return
	}

	productData, err := h.product.GetProducts(r.Context(), &proto.GetProductRequest{ProductId: productID})
	if err != nil {
		http.Error(w, "unable to fetch product", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"authorized": true,
		"product":    productData,
	})
}

// Client remaining for the ShowProduct gRPC service server streaming need to be implemented here in handler.go file
func (h *ProductHandler) ShowProductGateway(w http.ResponseWriter, r *http.Request) {
	// Implementation for ShowProduct gRPC service
	productId := strings.Trim(strings.TrimPrefix(r.URL.Path, "/gateway/showproduct/"), "/")
	if productId == "" {
		productId = r.URL.Query().Get("product_id")
	}
	if productId == "" {
		http.Error(w, "missing product_id", http.StatusBadRequest)
		return
	}

	stream, err := h.showProduct.ShowProduct(r.Context(), &proto.ShowProductRequest{ProductId: productId})
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to fetch product: %v", err), http.StatusBadGateway)
		return
	}
	var products []map[string]any

	for {
		res, err := stream.Recv()
		if err != nil {
			// Check if error is EOF (end of stream) - this is normal
			if errors.Is(err, io.EOF) {
				break
			}
			// For other errors, return error response
			http.Error(w, fmt.Sprintf("error receiving product stream: %v", err), http.StatusBadGateway)
			return
		}
		product := map[string]any{
			"id":          res.Id,
			"name":        res.Name,
			"description": res.Description,
			"price":       res.Price,
		}
		products = append(products, product)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"products": products,
	})

}
