.PHONY: help build run test proto docker-up docker-down clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the worker binary
	go build -o bin/worker ./cmd/server

run: ## Run the worker locally
	go run ./cmd/server/main.go

test: ## Run all tests
	go test -v ./...

test-infra: ## Run infrastructure tests (requires Redis & MinIO)
	go test -v ./internal/infrastructure/...

proto: ## Generate protobuf files
	./scripts/generate-proto.sh

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show worker logs
	docker-compose logs -f worker

docker-build: ## Build worker Docker image
	docker-compose build worker

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

.DEFAULT_GOAL := help
