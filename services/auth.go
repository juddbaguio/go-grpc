package services

import (
	"context"
	"go-grpc/grpc/auth"
)

type Auth struct {
	auth.UnimplementedAuthServiceServer
}

func (a *Auth) HandleLogin(ctx context.Context, login *auth.Login) (*auth.LoginResponse, error) {
	return &auth.LoginResponse{
		AuthToken:    login.Username,
		RefreshToken: login.Password,
	}, nil
}
