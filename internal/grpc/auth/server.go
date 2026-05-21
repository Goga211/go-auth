package auth

import (
	"context"

	auth "github.com/Goga211/proto/gen/go/auth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	auth.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	auth.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	panic("implement me")
}
