package server

import (
	authv1 "auth/protos/gen/go"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyStringValue = ""
	emptyIntValue    = 0
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

// Заглушки
func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateRequest(req.GetAppId(), emptyIntValue); err != nil {
		return nil, err
	}
	if err := validateRequest(req.GetEmail(), emptyStringValue); err != nil {
		return nil, err
	}
	if err := validateRequest(req.GetPassword(), emptyStringValue); err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{Token: req.GetEmail()}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("implement me Register()")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("implement me IsAdmin()")
}

func validateRequest(input interface{}, emptyValue interface{}) error {
	if input == emptyValue {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s is required", input))
	}
	return nil
}
