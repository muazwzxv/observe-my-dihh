package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/muazwzxv/otel_api_demo/internal/handlers/util"
	"github.com/muazwzxv/otel_api_demo/internal/models"
	"github.com/muazwzxv/otel_api_demo/internal/telemetry"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
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

	dummyErr := errors.New("FORBIDDEN_ACCESS_TO_WRITES")
	span.RecordError(dummyErr)
	span.SetStatus(codes.Error, dummyErr.Error())

	slog.ErrorContext(r.Context(), fmt.Sprintf("Error serving reqeust, error: %+v", dummyErr))

	util.WriteJSON(w, http.StatusInternalServerError, util.SuccessResponse{Message: "SOMETHING WRONG HAPPENED"})
}
