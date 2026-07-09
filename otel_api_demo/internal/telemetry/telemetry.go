package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/muazwzxv/otel_api_demo/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	Tracer = otel.Tracer("otel_api_demo")
	Meter  = otel.Meter("otel_api_demo")
)

func Setup(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	res, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(semconv.ServiceName(cfg.OTELServiceName)),
		sdkresource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	traceExp, err := newTraceExporter(ctx, cfg.OTELExporterEndpoint, cfg.OTELExporterProtocol)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	metricExp, err := newMetricExporter(ctx, cfg.OTELExporterEndpoint, cfg.OTELExporterProtocol)
	if err != nil {
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	baseHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(&traceLogHandler{base: baseHandler}))

	slog.InfoContext(ctx, "opentelemetry initialized",
		"service", cfg.OTELServiceName,
		"protocol", cfg.OTELExporterProtocol,
		"endpoint", cfg.OTELExporterEndpoint,
	)

	return func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("trace provider: %w", err)
		}
		return mp.Shutdown(ctx)
	}, nil
}

func newTraceExporter(ctx context.Context, endpoint, protocol string) (sdktrace.SpanExporter, error) {
	if protocol == "http" {
		return otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
	}
	return otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
}

func newMetricExporter(ctx context.Context, endpoint, protocol string) (sdkmetric.Exporter, error) {
	if protocol == "http" {
		return otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(endpoint), otlpmetrichttp.WithInsecure())
	}
	return otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(endpoint), otlpmetricgrpc.WithInsecure())
}
