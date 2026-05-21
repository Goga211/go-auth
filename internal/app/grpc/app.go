package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Goga211/go-auth/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	auth.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const operation = "grpcapp.Run"

	log := app.log.With(slog.String("operation", operation))

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	log.Info("grpc server successfully started", "port", app.port, slog.String("address", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (app *App) Stop() error {
	const operation = "grpcapp.Stop"

	app.log.With(slog.String("operation", operation)).Info("stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()

	return nil
}
