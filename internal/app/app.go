package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Goga211/go-auth/internal/app/grpc"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: инициализировать хранилище
	// TODO: init auth service
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
