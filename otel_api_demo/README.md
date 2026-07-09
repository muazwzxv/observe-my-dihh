# otel_api_demo

Scaffolded with chi, MySQL, Redis, Kafka, sqlc, Goose, and OpenTelemetry.

## Setup

```bash
# Start infrastructure (MySQL, Redis, Kafka, OTel Collector)
make docker-up

# Run migrations
make migrate-up

# Generate sqlc models
make sqlc-generate

# Run app
make run-api
```

## API

The API will be available at `http://localhost:8080`

### Endpoints

- `GET /v1/users/:id` - Get user by ID

## Observability

OpenTelemetry is configured out of the box. All HTTP requests, database queries, and Redis commands are automatically instrumented with traces and metrics.

### OTel Collector

Started by `make docker-up`. Endpoints:

| Port | Protocol | Purpose |
|------|----------|---------|
| 4317 | OTLP gRPC | Traces, metrics, logs ingestion |
| 4318 | OTLP HTTP | Alternative protocol |
| 8889 | Prometheus | Metrics scrape endpoint |

### Configuration

```bash
# .env or environment
OTEL_SERVICE_NAME=otel_api_demo          # Service name in traces
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317  # Collector address
OTEL_EXPORTER_OTLP_PROTOCOL=grpc            # grpc or http
```

### Traces

Traces are exported to the collector and printed to Docker logs. Add a visual backend by editing `otel-collector-config.yaml` — add Jaeger, Grafana Tempo, or any OTLP-compatible backend.

### Metrics

Application and infrastructure metrics are exported to the collector. Prometheus metrics are available at `http://localhost:8889/metrics`.

### Logs

Application logs use Go's `log/slog` with structured JSON output to stdout. Trace context (`trace_id`, `span_id`) is automatically included in log records when called with a context.

## Development Commands

```bash
make docker-up        # Start Docker services
make docker-down      # Stop Docker services
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make migrate-create   # Create new migration
make sqlc-generate    # Generate SQLC models
make run-api          # Run the API
make build-api        # Build binary
```
