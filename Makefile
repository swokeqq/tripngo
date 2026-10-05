
MIGRATIONS_DIR ?= ./migrations

export
-include .env

.PHONY: dev generate migrate-up migrate-down

dev:
	go run ./cmd/trip-service

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

generate:
	go tool oapi-codegen \
  	-generate types,chi-server \
  	-package api \
  	-o internal/generated/api.gen.go \
 	contracts/openapi/trip-service.openapi.yaml