.PHONY: dev dev-logs dev-down dev-restart dev-rebuild prod prod-logs prod-down prod-restart build build-cli-docker rebuild-cli-docker clean shell test help cli gas-prices-update

# ============================================
# Local Build Commands
# ============================================

build:
	@echo "🔨 Building server binary..."
	@go build -o bin/server ./cmd/server
	@echo "🔨 Building CLI binary..."
	@go build -o bin/cli ./cmd/cli
	@echo "✅ Build complete: bin/server and bin/cli"

build-server:
	@echo "🔨 Building server binary..."
	@go build -o bin/server ./cmd/server
	@echo "✅ Server binary ready: bin/server"

build-cli:
	@echo "🔨 Building CLI binary..."
	@go build -o bin/cli ./cmd/cli
	@echo "✅ CLI binary ready: bin/cli"

# ============================================
# Docker Build Commands (for CLI usage)
# ============================================

build-cli-docker:
	@echo "🔨 Building CLI binary in Docker container..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "✅ CLI binary ready in container: bin/cli"

rebuild-cli-docker:
	@echo "🔨 Rebuilding CLI binary in Docker container..."
	@docker-compose exec api rm -f bin/cli
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "✅ CLI binary rebuilt in container"

cli:
	@docker-compose exec api ./bin/cli

cli-help:
	@docker-compose exec api ./bin/cli --help

# ============================================
# Development commands (docker-compose.yml)
# ============================================

dev:
	docker-compose up -d

dev-logs:
	docker-compose logs -f

dev-down:
	docker-compose down

dev-restart:
	docker-compose restart

dev-rebuild:
	docker-compose down
	docker-compose up --build

# ============================================
# Production commands (docker-compose.prod.yml)
# ============================================

prod:
	docker-compose -f docker-compose.prod.yml up --build

prod-d:
	docker-compose -f docker-compose.prod.yml up -d --build

prod-logs:
	docker-compose -f docker-compose.prod.yml logs -f

prod-down:
	docker-compose -f docker-compose.prod.yml down

prod-restart:
	docker-compose -f docker-compose.prod.yml restart

prod-rebuild:
	docker-compose -f docker-compose.prod.yml down
	docker-compose -f docker-compose.prod.yml up --build

# ============================================
# Build specific images
# ============================================

build-dev:
	docker build --target development -t api-api:dev .

build-prod:
	docker build --target production -t api-api:prod .

# ============================================
# CLI Commands (via Docker)
# ============================================

gas-prices-update:
	@echo "🔨 Building CLI..."
	@docker-compose exec api go build -o bin/cli ./cmd/cli
	@echo "🔄 Running gas:prices:update command..."
	@docker-compose exec api ./bin/cli gas:prices:update

# ============================================
# Utility commands
# ============================================

shell:
	docker-compose exec api sh

shell-prod:
	docker-compose -f docker-compose.prod.yml exec api sh

test:
	docker-compose exec api go test ./... -v

clean:
	docker-compose down -v
	docker-compose -f docker-compose.prod.yml down -v
	rm -rf tmp/ bin/
	docker system prune -f

clean-all:
	docker-compose down -v --rmi all
	docker-compose -f docker-compose.prod.yml down -v --rmi all
	rm -rf tmp/ bin/
	docker system prune -af --volumes

lint:
	docker-compose exec api golangci-lint run --fix

# ============================================
# Help
# ============================================

help:
	@echo "╔════════════════════════════════════════════════════════╗"
	@echo "║        StationHub API - Available Commands            ║"
	@echo "╚════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📦 Build Commands:"
	@echo "  make build              Build both server and CLI binaries (local)"
	@echo "  make build-server       Build only server binary (local)"
	@echo "  make build-cli          Build only CLI binary (local)"
	@echo "  make build-cli-docker   Build CLI binary in Docker container"
	@echo "  make rebuild-cli-docker Rebuild CLI binary in Docker container"
	@echo ""
	@echo "🚀 Development:"
	@echo "  make dev                Start development environment"
	@echo "  make dev-logs           View development logs"
	@echo "  make dev-down           Stop development environment"
	@echo "  make dev-restart        Restart development containers"
	@echo "  make dev-rebuild        Rebuild development environment"
	@echo ""
	@echo "🏭 Production:"
	@echo "  make prod               Build and start production"
	@echo "  make prod-d             Build and start production (detached)"
	@echo "  make prod-logs          View production logs"
	@echo "  make prod-down          Stop production"
	@echo "  make prod-restart       Restart production"
	@echo ""
	@echo "🔨 Docker Build:"
	@echo "  make build-dev          Build development image"
	@echo "  make build-prod         Build production image"
	@echo ""
	@echo "🤖 CLI Commands (via Docker):"
	@echo "  make cli                Run CLI (interactive)"
	@echo "  make cli-help           Show CLI help"
	@echo "  make gas-prices-update  Update gas prices"
	@echo ""
	@echo "🛠️  Utility:"
	@echo "  make shell              Open shell in dev container"
	@echo "  make shell-prod         Open shell in prod container"
	@echo "  make test               Run tests"
	@echo "  make lint               Run linter"
	@echo "  make clean              Clean up containers and tmp files"
	@echo "  make clean-all          Deep clean (removes images too)"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  See CLI.md for detailed CLI usage"
	@echo "  Use ./cli.sh for helper script (local/docker)"
	@echo ""