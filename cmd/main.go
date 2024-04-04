package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// инициализировать проект
	cfg := config.EnvLoad()

	fmt.Println(cfg)

	// инициализировать логгер

	log := setupLogger(cfg.Env)

	log.Info("Starting Application ...", slog.String("env", cfg.Env))

	// инициализировать приложение

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// start
	go application.GRPCApp.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	signalStop := <-stop
	application.GRPCApp.Stop()
	log.Info("application stopped", slog.String("signal", signalStop.String()))
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
