package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// Get product service URL from environment or use default
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "localhost:9091"
	}

	// Get product ID from command line argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <product_id>", os.Args[0])
	}
	productID := os.Args[1]

	fmt.Println(strings.Repeat("=", 62))
	fmt.Println("  Product gRPC Client Test")
	fmt.Println(strings.Repeat("=", 62))
	fmt.Printf("\nConnecting to product gRPC server at %s...\n", productServiceURL)

	conn, err := grpc.NewClient(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("✓ Successfully connected to gRPC server!")
	fmt.Printf("\nFetching product with ID: %s\n\n", productID)

	// Create context with timeout for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := proto.NewGetProductsClient(conn)
	req := &proto.GetProductRequest{
		ProductId: productID,
	}

	fmt.Println("Sending request to gRPC server...")
	res, err := client.GetProducts(ctx, req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("✗ Request timeout: gRPC server did not respond within 10 seconds\n")
			fmt.Printf("  Make sure the server is running and accessible\n")
		} else {
			st, ok := status.FromError(err)
			if ok {
				fmt.Printf("✗ gRPC Error: %s (Code: %s)\n", st.Message(), st.Code())
			} else {
				fmt.Printf("✗ Error: %v\n", err)
			}
		}
		os.Exit(1)
	}

	fmt.Println("✓ SUCCESS! Product retrieved successfully")
	fmt.Printf("\nProduct Details:\n")
	fmt.Printf("  ID:          %s\n", res.Id)
	fmt.Printf("  Name:        %s\n", res.Name)
	fmt.Printf("  Description: %s\n", res.Description)
	fmt.Printf("  Price:       $%.2f\n", res.Price)

	fmt.Println("\n" + strings.Repeat("=", 62))
	fmt.Println("gRPC client is working perfectly!")
	fmt.Println(strings.Repeat("=", 62))
}
