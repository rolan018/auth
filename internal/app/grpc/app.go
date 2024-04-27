package grpcapp

import (
	server "auth/internal/grpc"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New create New gRPC server app.
func New(log *slog.Logger, port int, auth server.Auth) *App {
	gRPCServer := grpc.NewServer()

	server.Register(gRPCServer, auth)

	return &App{log: log, gRPCServer: gRPCServer, port: port}
}

func (a *App) run() error {
	const op = "grpcapp.run()"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running ...", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Must Run gRPC server
func (a *App) MustRun() {
	err := a.run()
	if err != nil {
		panic(err)
	}
}

// For graceful shutdown
func (a *App) Stop() {
	const op = "grpcapp.Stop()"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
