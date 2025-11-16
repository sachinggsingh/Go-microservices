package grpc_handler

import (
	"context"
	"fmt"

	"github.com/sachinggsingh/e-comm/internal/helper"
	proto "github.com/sachinggsingh/e-comm/pb"
)

type AuthServer struct {
	proto.UnimplementedValidateTokenServer
}

func (a *AuthServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (res *proto.ValidateTokenResponse, err error) {
	fmt.Printf("[gRPC] Received ValidateToken request\n")

	claims, err := helper.ValidateToken(req.Token)
	if err != nil {
		fmt.Printf("[gRPC] Token validation failed: %v\n", err)
		return &proto.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	fmt.Printf("[gRPC] Token validated successfully - Authenticated User: %s (Email: %s)\n", claims.Uid, claims.Email)
	return &proto.ValidateTokenResponse{
		Valid: true,
	}, nil
}
