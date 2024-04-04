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

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	CreateApp(ctx context.Context, email string, password string, app_name string, app_secret string) (int64, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	// service layer
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	// service layer
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	// service layer
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func (s *serverAPI) CreateApp(ctx context.Context, req *authv1.CreateAppRequest) (*authv1.CreateAppResponse, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}
	// service layer
	appID, err := s.auth.CreateApp(ctx, req.GetEmail(), req.GetPassword(), req.GetAppName(), req.GetAppSecret())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authv1.CreateAppResponse{AppId: appID}, nil
}

// Validators
func validateCreate(req *authv1.CreateAppRequest) error {
	if err := validateRequest(req.GetAppName(), emptyStringValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetAppSecret(), emptyStringValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetEmail(), emptyStringValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetPassword(), emptyStringValue); err != nil {
		return err
	}
	return nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if err := validateRequest(req.GetAppId(), emptyIntValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetEmail(), emptyStringValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetPassword(), emptyStringValue); err != nil {
		return err
	}
	return nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if err := validateRequest(req.GetEmail(), emptyStringValue); err != nil {
		return err
	}
	if err := validateRequest(req.GetPassword(), emptyStringValue); err != nil {
		return err
	}
	return nil
}

func validateIsAdmin(req *authv1.IsAdminRequest) error {
	if err := validateRequest(req.GetUserId(), emptyIntValue); err != nil {
		return err
	}
	return nil
}

func validateRequest(input interface{}, emptyValue interface{}) error {
	if input == emptyValue {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s is required", input))
	}
	return nil
}
