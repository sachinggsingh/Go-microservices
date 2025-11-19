package grpc_client

import (
	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	AuthClient    proto.ValidateTokenClient
	ProductClient proto.GetProductsClient
	ShowProduct   proto.ShowProductClient
}

func NewClientWithConn(authConn *grpc.ClientConn, productConn *grpc.ClientConn, showProductConn *grpc.ClientConn) *Clients {
	return &Clients{
		AuthClient:    proto.NewValidateTokenClient(authConn),
		ProductClient: proto.NewGetProductsClient(productConn),
		ShowProduct:   proto.NewShowProductClient(showProductConn),
	}
}

func NewClient(authAddr string, productAddr string, showProductAddr string) (*Clients, error) {
	authConn, err := Dial(authAddr)
	if err != nil {
		return nil, err
	}

	productConn, err := Dial(productAddr)
	if err != nil {
		return nil, err
	}

	// ShowProduct is registered on the same gRPC server as GetProducts
	// So we reuse the productConn instead of creating a separate connection
	return NewClientWithConn(authConn, productConn, productConn), nil
}

func Dial(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
