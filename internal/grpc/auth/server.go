package auth

import (
	"context"
	"errors"

	auth1 "github.com/Goga211/go-auth/internal/services/auth"
	"github.com/Goga211/go-auth/storage"

	auth "github.com/Goga211/proto/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type serverAPI struct {
	auth.UnimplementedAuthServer
	authentication Auth
}

func Register(gRPC *grpc.Server, authentication Auth) {
	auth.RegisterAuthServer(gRPC, &serverAPI{authentication: authentication})
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := ValidateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.authentication.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppID()))

	if err != nil {
		if errors.Is(err, auth1.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &auth.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if err := ValidateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.authentication.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &auth.RegisterResponse{
		UserId: userID,
	}, nil
}

func ValidateLogin(req *auth.LoginRequest) error {
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	return nil
}

func ValidateRegister(req *auth.RegisterRequest) error {
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	return nil
}
