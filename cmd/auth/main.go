package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Goga211/go-auth/internal/app"
	"github.com/Goga211/go-auth/internal/config"
	"github.com/Goga211/go-auth/internal/lib/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	_ = godotenv.Load()
	viper.AutomaticEnv()

	cfg := config.Load()
	log := logger.SetupLogger()

	log.Info("Starting app...")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GrpcServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcServer.Stop()
	log.Info("Shutting down...")

}
