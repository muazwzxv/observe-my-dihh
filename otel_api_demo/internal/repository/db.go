package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"github.com/muazwzxv/otel_api_demo/internal/config"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := otelsql.Open("mysql", cfg.GetDSN(),
		otelsql.WithAttributes(semconv.DBSystemMySQL),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.GetRedisAddr(),
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		slog.Warn("failed to instrument redis tracing", "error", err)
	}
	if err := redisotel.InstrumentMetrics(client); err != nil {
		slog.Warn("failed to instrument redis metrics", "error", err)
	}

	return client, nil
}
