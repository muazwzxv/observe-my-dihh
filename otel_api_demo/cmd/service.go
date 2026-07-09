package cmd

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/muazwzxv/otel_api_demo/internal/models"
)

type APIService struct {
	DB      *sql.DB
	Queries *models.Queries
	Redis   *redis.Client
}
