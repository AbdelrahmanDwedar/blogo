.PHONY: help run build clean test db-create db-drop redis-start install

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

run: ## Run the application
	@echo "🚀 Starting Blogo..."
	@go run cmd/api/main.go

build: ## Build the application
	@echo "🔨 Building Blogo..."
	@go build -o blogo cmd/api/main.go
	@echo "✅ Build complete: ./blogo"

clean: ## Clean build files
	@echo "🧹 Cleaning..."
	@rm -f blogo
	@echo "✅ Clean complete"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

db-create: ## Create PostgreSQL database
	@echo "📊 Creating database..."
	@createdb blogo || echo "Database might already exist"
	@echo "✅ Database ready"

db-drop: ## Drop PostgreSQL database
	@echo "⚠️  Dropping database..."
	@dropdb blogo || echo "Database might not exist"
	@echo "✅ Database dropped"

redis-start: ## Start Redis server
	@echo "🔴 Starting Redis..."
	@redis-server &

setup: install ## Setup project (install deps and create env file)
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env file..."; \
		cp env.example .env; \
		echo "✅ .env file created"; \
		echo "⚠️  Please edit .env and update the configuration values!"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

dev: ## Run in development mode with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run with: make run"; \
	fi

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t blogo:latest .
	@echo "✅ Docker image built"

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 --env-file .env blogo:latest

fmt: ## Format Go code
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not found. Install it from: https://golangci-lint.run/usage/install/"; \
	fi

deps-update: ## Update dependencies
	@echo "📦 Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@go run cmd/seed/main.go

reset-db: db-drop db-create seed ## Reset database and seed with sample data
	@echo "✅ Database reset complete"

