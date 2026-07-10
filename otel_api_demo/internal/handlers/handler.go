package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/muazwzxv/otel_api_demo/cmd"
	"github.com/muazwzxv/otel_api_demo/internal/handlers/user"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupHandler(ctx context.Context, router chi.Router, svc *cmd.APIService) {
	router.Use(otelhttp.NewMiddleware("otel_api_demo"))
	router.Use(LoggingMiddleware())

	setupUserHandlers(ctx, router, svc)
}

func setupUserHandlers(ctx context.Context, router chi.Router, svc *cmd.APIService) {
	getByIdHandlder := &user.GetByIDHandler{
		DB:      svc.DB,
		Queries: svc.Queries,
		Redis:   svc.Redis,
	}
	router.Get("/api/v1/users/{id}", getByIdHandlder.Handle)

	createUserHandler := &user.CreateUserHandler{
		DB:      svc.DB,
		Queries: svc.Queries,
		Redis:   svc.Redis,
	}

	router.Post("/api/v1/user", createUserHandler.Handle)

	slog.InfoContext(ctx, "Registered user handlers")
}

func LoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			slog.InfoContext(r.Context(), fmt.Sprintf("request, path: %s", r.URL.Path),
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.statusCode,
				"ip", r.RemoteAddr,
			)
			next.ServeHTTP(lrw, r)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
