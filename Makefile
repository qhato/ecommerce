
.PHONY: help build run-admin run-storefront test clean docker-build docker-up docker-down migrate

# Variables
ADMIN_BINARY=bin/admin
STOREFRONT_BINARY=bin/storefront
DOCKER_COMPOSE=docker-compose

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: ## Build all binaries
	@echo "Building Admin API..."
	go build -o $(ADMIN_BINARY) cmd/admin/main.go
	@echo "Building Storefront API..."
	go build -o $(STOREFRONT_BINARY) cmd/storefront/main.go
	@echo "Build complete!"

build-admin: ## Build Admin API binary
	@echo "Building Admin API..."
	go build -o $(ADMIN_BINARY) cmd/admin/main.go

build-storefront: ## Build Storefront API binary
	@echo "Building Storefront API..."
	go build -o $(STOREFRONT_BINARY) cmd/storefront/main.go

run-admin: ## Run Admin API (development)
	@echo "Starting Admin API..."
	go run cmd/admin/main.go

run-storefront: ## Run Storefront API (development)
	@echo "Starting Storefront API..."
	go run cmd/storefront/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v -race -short ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	go test -v -race -run Integration ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

test-coverage-ci: ## Run tests and check coverage threshold (70%)
	@echo "Running tests with coverage check..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "Checking coverage threshold..."
	@go tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 70.0) {print "Coverage " $$3 " is below 70%"; exit 1} else {print "Coverage " $$3 " is above 70%"}}'

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

test-watch: ## Watch for changes and run tests (requires gotestsum)
	@which gotestsum > /dev/null || (echo "Installing gotestsum..." && go install gotest.tools/gotestsum@latest)
	gotestsum --watch

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete!"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

lint-fix: ## Run linter and fix issues
	@echo "Running linter with auto-fix..."
	golangci-lint run --fix ./...

# OpenAPI/Swagger Documentation

swagger-install: ## Install swag CLI for OpenAPI generation
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)

swagger-gen: swagger-install ## Generate OpenAPI/Swagger documentation
	@echo "Generating OpenAPI documentation..."
	@swag init -g cmd/admin/main.go -o docs/admin --parseDependency --parseInternal
	@swag init -g cmd/storefront/main.go -o docs/storefront --parseDependency --parseInternal
	@echo "OpenAPI documentation generated in docs/"

swagger-serve: ## Serve Swagger UI locally (requires swagger-ui-dist)
	@echo "Serving Swagger UI at http://localhost:8090"
	@echo "Visit http://localhost:8090/admin for Admin API docs"
	@echo "Visit http://localhost:8090/storefront for Storefront API docs"

# Observability commands

obs-stack-up: ## Start observability stack (Prometheus, Grafana, Jaeger)
	@echo "Starting observability stack..."
	cd examples/observability-server && docker-compose up -d prometheus grafana jaeger redis postgres
	@echo "Observability stack started!"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"
	@echo "  - Jaeger UI: http://localhost:16686"

obs-stack-down: ## Stop observability stack
	@echo "Stopping observability stack..."
	cd examples/observability-server && docker-compose down
	@echo "Observability stack stopped!"

obs-example: ## Run observability example server
	@echo "Running observability example..."
	go run examples/observability-server/main.go

metrics: ## View Prometheus metrics
	@echo "Opening metrics endpoint..."
	@curl -s http://localhost:8080/metrics || echo "Server not running on port 8080"

health: ## Check application health
	@echo "Checking application health..."
	@curl -s http://localhost:8080/health | jq . || echo "Server not running or jq not installed"

health-live: ## Check liveness probe
	@curl -s http://localhost:8080/health/live || echo "Server not running on port 8080"

health-ready: ## Check readiness probe
	@curl -s http://localhost:8080/health/ready || echo "Server not running on port 8080"

# Docker commands

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up: ## Start all services with Docker Compose
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d
	@echo "Services started!"

docker-down: ## Stop all services
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down
	@echo "Services stopped!"

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show running containers
	$(DOCKER_COMPOSE) ps

# Database commands

db-create: ## Create database
	@echo "Creating database..."
	createdb ecommerce || echo "Database already exists"

db-drop: ## Drop database
	@echo "Dropping database..."
	dropdb ecommerce

db-migrate: ## Run database migrations
	@echo "Running migrations..."
	psql ecommerce < database.sql
	@echo "Migrations complete!"

db-reset: db-drop db-create db-migrate ## Reset database (drop, create, migrate)

db-shell: ## Open PostgreSQL shell
	psql ecommerce

# Development commands

dev-admin: ## Run Admin API in development mode with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air -c .air-admin.toml

dev-storefront: ## Run Storefront API in development mode with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air -c .air-storefront.toml

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Production commands

build-prod: ## Build production binaries with optimizations
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o $(ADMIN_BINARY) cmd/admin/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o $(STOREFRONT_BINARY) cmd/storefront/main.go
	@echo "Production build complete!"

# Multi-platform build commands

build-all-platforms: ## Build for all platforms (Linux, macOS, Windows)
	@./scripts/build.sh all

build-linux: ## Build for Linux (amd64, arm64)
	@./scripts/build.sh linux

build-macos: ## Build for macOS (amd64, arm64)
	@./scripts/build.sh macos

build-windows: ## Build for Windows (amd64)
	@./scripts/build.sh windows

build-release: ## Build release archives for all platforms
	@./scripts/build.sh release

build-clean: ## Clean build directory
	@./scripts/build.sh clean

build-admin-only: ## Build only admin binary for all platforms
	@./scripts/build.sh all --admin-only

build-storefront-only: ## Build only storefront binary for all platforms
	@./scripts/build.sh all --storefront-only

# Utilities

.DEFAULT_GOAL := help
