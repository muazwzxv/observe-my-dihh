package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/muazwzxv/otel_api_demo/cmd"
	"github.com/muazwzxv/otel_api_demo/internal/config"
	"github.com/muazwzxv/otel_api_demo/internal/handlers"
	"github.com/muazwzxv/otel_api_demo/internal/models"
	"github.com/muazwzxv/otel_api_demo/internal/repository"
	"github.com/muazwzxv/otel_api_demo/internal/telemetry"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bootupCtx := context.Background()

	shutdown, err := telemetry.Setup(bootupCtx, cfg)
	if err != nil {
		slog.Error("failed to setup telemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			slog.Error("failed to shutdown telemetry", "error", err)
		}
	}()

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient, err := repository.NewRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	apiService := &cmd.APIService{
		DB:      db,
		Queries: models.New(),
		Redis:   redisClient,
	}

	router := chi.NewRouter()

	handlers.SetupHandler(bootupCtx, router, apiService)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	slog.InfoContext(bootupCtx, fmt.Sprintf("Server starting on port %s", cfg.ServerPort))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.ErrorContext(bootupCtx, fmt.Sprintf("Failed to start server: %v", err))
		os.Exit(1)
	}
}
