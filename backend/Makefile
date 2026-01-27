# GoFund Backend Makefile

.PHONY: help dev-setup atlas-install atlas-generate atlas-apply atlas-inspect clean test build docker-up docker-down

# Default target
help:
	@echo "GoFund Backend Development Commands"
	@echo "=================================="
	@echo ""
	@echo "Setup Commands:"
	@echo "  dev-setup      - Setup development environment (databases, migrations)"
	@echo "  atlas-install  - Install Atlas CLI"
	@echo ""
	@echo "Database Commands:"
	@echo "  atlas-generate - Generate new migration from schema changes"
	@echo "  atlas-apply    - Apply pending migrations"
	@echo "  atlas-inspect  - Inspect current database schema"
	@echo ""
	@echo "Development Commands:"
	@echo "  docker-up      - Start all infrastructure services"
	@echo "  docker-down    - Stop all services"
	@echo "  test           - Run tests"
	@echo "  build          - Build all services"
	@echo "  clean          - Clean build artifacts"

# Development setup
dev-setup:
	@echo "🚀 Setting up GoFund development environment..."
	@chmod +x scripts/*.sh
	@./scripts/setup-dev-db.sh

# Atlas CLI installation
atlas-install:
	@echo "📦 Installing Atlas CLI..."
	@curl -sSf https://atlasgo.sh | sh

# Database migration commands
atlas-generate:
	@echo "🔄 Generating Atlas migration..."
	@chmod +x scripts/atlas-generate.sh
	@./scripts/atlas-generate.sh dev

atlas-apply:
	@echo "🔄 Applying Atlas migrations..."
	@chmod +x scripts/atlas-apply.sh
	@./scripts/atlas-apply.sh dev

atlas-inspect:
	@echo "🔍 Inspecting database schema..."
	@chmod +x scripts/atlas-inspect.sh
	@./scripts/atlas-inspect.sh dev

# Docker commands
docker-up:
	@echo "🐳 Starting Docker services..."
	@docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker services..."
	@docker-compose down

# Build commands
build:
	@echo "🔨 Building all services..."
	@go work sync
	@cd services/users-service && go build -o ../../bin/users-service ./cmd
	@cd services/goals-service && go build -o ../../bin/goals-service ./cmd
	@cd services/ledger-service && go build -o ../../bin/ledger-service ./cmd
	@cd services/payments-service && go build -o ../../bin/payments-service ./cmd
	@cd services/notifications-service && go build -o ../../bin/notifications-service ./cmd

# Test commands
test:
	@echo "🧪 Running tests..."
	@go work sync
	@go test ./...

# Clean commands
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@go clean -cache
	@go clean -modcache

# Development workflow
dev: dev-setup atlas-apply
	@echo "✅ Development environment ready!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Start services: make build && ./bin/users-service"
	@echo "  2. Test API: curl http://localhost:8080/health"

# Production migration
migrate-prod:
	@echo "🚀 Applying production migrations..."
	@chmod +x scripts/atlas-apply.sh
	@./scripts/atlas-apply.sh prod