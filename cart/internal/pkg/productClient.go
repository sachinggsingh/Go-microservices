package pkg

import (
	"context"
	"fmt"

	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ProductClient struct {
	client proto.GetProductsClient
	conn   *grpc.ClientConn
}

func NewProductClient(productServiceURL string) (*ProductClient, error) {
	conn, err := grpc.NewClient(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := proto.NewGetProductsClient(conn)

	return &ProductClient{
		client: client,
		conn:   conn,
	}, nil
}

func (pc *ProductClient) GetProduct(ctx context.Context, productID string) (*proto.GetPRoductResponse, error) {
	if productID == "" {
		return nil, fmt.Errorf("product_id is required")
	}

	req := &proto.GetProductRequest{
		ProductId: productID,
	}

	resp, err := pc.client.GetProducts(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, fmt.Errorf("gRPC error: %s (code: %s)", st.Message(), st.Code())
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp, nil
}

func (pc *ProductClient) Close() error {
	if pc.conn != nil {
		return pc.conn.Close()
	}
	return nil
}

// ValidateProduct validates if a product exists and returns its price
func (pc *ProductClient) ValidateProduct(ctx context.Context, productID string) (float64, error) {
	product, err := pc.GetProduct(ctx, productID)
	if err != nil {
		return 0, err
	}

	if product.Price <= 0 {
		return 0, fmt.Errorf("invalid product price: %f", product.Price)
	}

	return product.Price, nil
}
