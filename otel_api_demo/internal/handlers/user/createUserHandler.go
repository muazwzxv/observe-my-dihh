package user

import (
	"database/sql"
	"net/http"

	"github.com/muazwzxv/otel_api_demo/internal/models"
	"github.com/muazwzxv/otel_api_demo/internal/telemetry"
	"github.com/redis/go-redis/v9"
)

type CreateUserHandler struct {
	DB      *sql.DB
	Queries *models.Queries
	Redis   *redis.Client
}

func (h *CreateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := telemetry.Tracer.Start(ctx, "CreateUserHandler")
	defer span.End()
}
