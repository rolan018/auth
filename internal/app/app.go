package app

import (
	grpcapp "auth/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

// create New gRPC server app.
func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: инициализировать хранилище

	// TODO: init auth service

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCApp: grpcApp}
}
