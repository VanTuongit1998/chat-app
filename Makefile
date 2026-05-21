include .env

.PHONY: help run build test fmt tidy clean docker-up docker-down docker-logs docker-ps migrate migrate-users migrate-rooms migrate-messages

APP_NAME := chat-app
SERVER_ENTRY := ./cmd/server
BUILD_DIR := bin
BUILD_OUTPUT := $(BUILD_DIR)/server

help:
	@echo "Available commands:"
	@echo "  make run              Run the Go server"
	@echo "  make build            Build the server binary"
	@echo "  make test             Run Go tests"
	@echo "  make fmt              Format Go files"
	@echo "  make tidy             Run go mod tidy"
	@echo "  make clean            Remove build output"
	@echo "  make docker-up        Start Redis and PostgreSQL"
	@echo "  make docker-down      Stop Docker services"
	@echo "  make docker-logs      Show Docker logs"
	@echo "  make docker-ps        Show Docker services"
	@echo "  make migrate          Run all SQL migrations"

run:
	go run $(SERVER_ENTRY)

build:
	go build -o $(BUILD_OUTPUT) $(SERVER_ENTRY)

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name "*.go" -not -path "./vendor/*")

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-ps:
	docker compose ps

migrate: migrate-users migrate-rooms migrate-messages

migrate-users:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < migrations/001_create_users.sql

migrate-rooms:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < migrations/002_create_rooms.sql

migrate-messages:
	docker compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < migrations/003_create_messages.sql
