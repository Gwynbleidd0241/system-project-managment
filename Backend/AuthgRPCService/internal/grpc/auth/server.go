package auth

import (
	customErrors "authService/internal/errors"
	"context"
	"errors"
	authv1 "github.com/Dragzet/gRPCProtosv2/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string) (string, error)
	RegisterNewUser(ctx context.Context, email string, password string) (string, error)
}

type ServerAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func RegisterServer(s *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(s, &ServerAPI{
		auth: auth,
	})
}

func (s *ServerAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "login error")
	}
	return &authv1.LoginResponse{Token: token}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	token, err := s.auth.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, customErrors.InvalidInputData) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.AlreadyExists, "user already exisrs")
	}
	return &authv1.RegisterResponse{Token: token}, nil
}
