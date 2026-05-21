package main

import (
	"github.com/Goga211/go-auth/internal/app"
	"github.com/Goga211/go-auth/internal/config"
	"github.com/Goga211/go-auth/internal/logger"
)

func main() {
	cfg := config.Load()
	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting app...")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	application.GrpcServer.MustRun()
}
