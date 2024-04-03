package app

import (
	grpcapp "auth/internal/app/grpc"
	auth "auth/internal/services"
	"auth/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcapp.App
}

// New create New server app.
func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: инициализировать хранилище
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	// TODO: init auth service
	authService := auth.New(log, tokenTTL, storage)
	// init grpc Server
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{GRPCApp: grpcApp}
}
