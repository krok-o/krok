package service

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	authv1 "github.com/krok-o/krok/proto/auth/v1"
)

type AuthService struct {
	authv1.UnimplementedAuthServiceServer
}

func (a *AuthService) Login(ctx context.Context, request *authv1.LoginRequest) (*empty.Empty, error) {
	fmt.Println("login?")
	return &empty.Empty{}, nil
}
