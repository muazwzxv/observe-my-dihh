.PHONY: docker-up docker-down migrate-up migrate-down migrate-create sqlc-generate run-api build-api

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

MIGRATION_DIR=database/migrations
DATABASE_URL="otel_api_demo_user:otel_api_demo_pass@tcp(localhost:3306)/otel_api_demo_db?parseTime=true"

migrate-up:
	goose -dir $(MIGRATION_DIR) mysql $(DATABASE_URL) up

migrate-down:
	goose -dir $(MIGRATION_DIR) mysql $(DATABASE_URL) down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATION_DIR) create $$name sql

sqlc-generate:
	sqlc generate

run-api:
	go run cmd/api/main.go

build-api:
	go build -o bin/api cmd/api/main.go
