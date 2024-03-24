package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// TODO: инициализировать проект
	cfg := config.EnvLoad()

	fmt.Println(cfg)

	// TODO: инициализировать логгер

	log := setupLogger(cfg.Env)

	log.Info("Starting Application ...", slog.String("env", cfg.Env))

	// TODO: инициализировать приложение

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// start
	application.GRPCApp.MustRun()
	// TODO: запустить gRPC-сервер приложения
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
