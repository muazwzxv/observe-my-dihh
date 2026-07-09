package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"github.com/muazwzxv/otel_api_demo/cmd"
	"github.com/muazwzxv/otel_api_demo/internal/handlers/user"
)

func SetupHandler(ctx context.Context, router chi.Router, svc *cmd.APIService) {
	router.Use(otelhttp.NewMiddleware("otel_api_demo"))
	router.Use(LoggingMiddleware())

	setupUserHandlers(ctx, router, svc)
}

func setupUserHandlers(ctx context.Context, router chi.Router, svc *cmd.APIService) {
	handler := &user.GetByIDHandler{
		DB:      svc.DB,
		Queries: svc.Queries,
		Redis:   svc.Redis,
	}
	router.Get("/v1/users/{id}", handler.Handle)
	slog.InfoContext(ctx, "Registered user handlers")
}

func LoggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)
			slog.InfoContext(r.Context(), "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.statusCode,
				"duration", time.Since(start),
				"ip", r.RemoteAddr,
			)
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
