
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

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

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
