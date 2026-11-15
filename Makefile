
MIGRATIONS_DIR ?= ./migrations

export
-include .env

.PHONY: run generate migrate migrate-down

run:
	go run ./cmd/trip-service

migrate:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

generate:
	go tool oapi-codegen \
  	-generate types,chi-server \
  	-package api \
  	-o internal/generated/api.gen.go \
 	contracts/openapi/trip-service.openapi.yaml