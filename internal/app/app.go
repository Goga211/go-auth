package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Goga211/go-auth/internal/app/grpc"
	"github.com/Goga211/go-auth/internal/services/auth"
	"github.com/Goga211/go-auth/storage/sqlite"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
