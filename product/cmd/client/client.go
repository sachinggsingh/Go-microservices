package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// Get token from environment variable, command line argument, or use default test token
	token := os.Getenv("TEST_TOKEN")
	if token == "" && len(os.Args) > 1 {
		token = os.Args[1]
	}
	if token == "" {
		fmt.Println("No token provided. Please set the TEST_TOKEN environment variable or provide a token as an argument.")
		return
	}

	fmt.Println("=" + strings.Repeat("=", 60) + "=")
	fmt.Println("  gRPC Token Validation Test Client")
	fmt.Println("=" + strings.Repeat("=", 60) + "=")
	fmt.Printf("\nConnecting to gRPC server at localhost:9090...\n")

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v\n", err)
	}
	defer conn.Close()

	fmt.Println(" Successfully connected to gRPC server!")
	fmt.Printf("\nTesting token validation...\n")
	tokenPreview := token
	if len(token) > 50 {
		tokenPreview = token[:50] + "..."
	}
	fmt.Printf("Token: %s\n\n", tokenPreview)

	// Create context with timeout for the RPC call
	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer rpcCancel()

	client := pb.NewValidateTokenClient(conn)
	req := pb.ValidateTokenRequest{
		Token: token,
	}

	fmt.Println("Sending request to gRPC server...")
	res, err := client.ValidateToken(rpcCtx, &req)
	if err != nil {
		if rpcCtx.Err() == context.DeadlineExceeded {
			fmt.Printf(" Request timeout: gRPC server did not respond within 10 seconds\n")
			fmt.Printf("   Make sure the server is running and accessible\n")
		} else {
			st, ok := status.FromError(err)
			if ok {
				fmt.Printf("gRPC Error: %s (Code: %s)\n", st.Message(), st.Code())
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}
		fmt.Printf("Token validation failed!\n")
		os.Exit(1)
	}

	if res.Valid {
		fmt.Println(" SUCCESS! Token is VALID")
		fmt.Printf("Response: %+v\n", res)
	} else {
		fmt.Println(" Token is INVALID")
		fmt.Printf("Response: %+v\n", res)
		os.Exit(1)
	}

	fmt.Println("\n" + strings.Repeat("=", 62))
	fmt.Println("gRPC is working perfectly!")
	fmt.Println(strings.Repeat("=", 62))
}
