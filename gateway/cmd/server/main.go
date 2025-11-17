package main

import (
	"log"
	"net/http"

	grpc_client "github.com/sachinggsingh/e-comm/internal/gRPC"
	"github.com/sachinggsingh/e-comm/internal/handler"
)

func main() {
	client, err := grpc_client.NewClient(
		"localhost:9090", // auth
		"localhost:9091", // product
	)
	if err != nil {
		log.Fatal(err)
	}

	productHandler := handler.NewProductHandler(client.AuthClient, client.ProductClient)
	mux := http.NewServeMux()
	mux.HandleFunc("/gateway/product/", productHandler.GetProductGateway)

	log.Println("Gateway running at :8085")
	http.ListenAndServe(":8085", mux)
}
