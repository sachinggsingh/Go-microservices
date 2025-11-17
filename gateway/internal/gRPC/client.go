package grpc_client

import (
	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	AuthClient    proto.ValidateTokenClient
	ProductClient proto.GetProductsClient
}

func NewClientWithConn(authConn *grpc.ClientConn, productConn *grpc.ClientConn) *Clients {
	return &Clients{
		AuthClient:    proto.NewValidateTokenClient(authConn),
		ProductClient: proto.NewGetProductsClient(productConn),
	}
}

func NewClient(authAddr string, productAddr string) (*Clients, error) {
	authConn, err := Dial(authAddr)
	if err != nil {
		return nil, err
	}

	productConn, err := Dial(productAddr)
	if err != nil {
		return nil, err
	}

	return NewClientWithConn(authConn, productConn), nil
}

func Dial(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
