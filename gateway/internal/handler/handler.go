package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	proto "github.com/sachinggsingh/e-comm/pb"
)

type ProductHandler struct {
	auth    proto.ValidateTokenClient
	product proto.GetProductsClient
}

func NewProductHandler(auth proto.ValidateTokenClient, product proto.GetProductsClient) *ProductHandler {
	return &ProductHandler{auth, product}
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
