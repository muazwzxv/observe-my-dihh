package telemetry

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type traceLogHandler struct {
	base slog.Handler
}

func (h *traceLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func (h *traceLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		r.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return h.base.Handle(ctx, r)
}

func (h *traceLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceLogHandler{base: h.base.WithAttrs(attrs)}
}

func (h *traceLogHandler) WithGroup(name string) slog.Handler {
	return &traceLogHandler{base: h.base.WithGroup(name)}
}
