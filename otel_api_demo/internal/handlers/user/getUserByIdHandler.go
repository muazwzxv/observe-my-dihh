package user

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/codes"

	"github.com/muazwzxv/otel_api_demo/internal/handlers/util"
	"github.com/muazwzxv/otel_api_demo/internal/models"
	"github.com/muazwzxv/otel_api_demo/internal/telemetry"
	"github.com/redis/go-redis/v9"
)

var (
	UserLookupsCounter, _ = telemetry.Meter.Int64Counter("user_lookups_total")
)

type GetByIDHandler struct {
	DB      *sql.DB
	Queries *models.Queries
	Redis   *redis.Client
}

func (h *GetByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := telemetry.Tracer.Start(ctx, "GetByIDHandler")
	defer span.End()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.HandleError(w, util.BuildErrorWithCode(
			http.StatusBadRequest,
			"Invalid ID format",
			"INVALID_ID_FORMAT",
		))
		return
	}

	user, err := h.Queries.GetUser(ctx, h.DB, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.WarnContext(ctx, "user not found", "user_id", id)
			util.HandleError(w, util.BuildErrorWithCode(
				http.StatusNotFound,
				"User not found",
				"USER_NOT_FOUND",
			))
			return
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to get user", "user_id", id, "error", err)
		util.HandleError(w, util.BuildErrorWithCode(
			http.StatusInternalServerError,
			"Failed to retrieve user",
			"USER_LOOKUP_FAILED",
		))
		return
	}

	UserLookupsCounter.Add(ctx, 1)
	slog.InfoContext(ctx, "user retrieved", "user_id", user.ID, "user_name", user.Name)
	util.WriteJSON(w, http.StatusOK, util.SuccessResponse{Data: user})
}
