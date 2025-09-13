# Makefile for Todo App Database Setup and Management

# Database configuration (from .env file)
DB_NAME = todoapp
DB_USER = todouser
DB_PASSWORD = todopassword123
DB_HOST = localhost
DB_PORT = 5432

# Colors for output
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help setup-db start-postgres stop-postgres restart-postgres check-postgres create-db create-user grant-privileges run-migrations dev build clean

# Default target
help:
	@echo "$(GREEN)Todo App Database Management$(NC)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(NC)"
	@echo "  make setup-db       - Complete database setup (start postgres + create db + user)"
	@echo "  make start-postgres - Start PostgreSQL service"
	@echo "  make stop-postgres  - Stop PostgreSQL service"
	@echo "  make restart-postgres - Restart PostgreSQL service"
	@echo "  make check-postgres - Check PostgreSQL status"
	@echo "  make create-db      - Create database and user"
	@echo "  make dev           - Run application in development mode"
	@echo "  make build         - Build application for production"
	@echo "  make clean         - Clean build artifacts"

# Start PostgreSQL service
start-postgres:
	@echo "$(YELLOW)Starting PostgreSQL service...$(NC)"
	sudo systemctl start postgresql
	sudo systemctl enable postgresql
	@echo "$(GREEN)PostgreSQL started and enabled$(NC)"

# Stop PostgreSQL service
stop-postgres:
	@echo "$(YELLOW)Stopping PostgreSQL service...$(NC)"
	sudo systemctl stop postgresql
	@echo "$(GREEN)PostgreSQL stopped$(NC)"

# Restart PostgreSQL service
restart-postgres:
	@echo "$(YELLOW)Restarting PostgreSQL service...$(NC)"
	sudo systemctl restart postgresql
	@echo "$(GREEN)PostgreSQL restarted$(NC)"

# Check PostgreSQL status
check-postgres:
	@echo "$(YELLOW)Checking PostgreSQL status...$(NC)"
	sudo systemctl status postgresql --no-pager

# Create database and user
create-db:
	@echo "$(YELLOW)Creating database and user...$(NC)"
	sudo -u postgres psql -c "SELECT 1 FROM pg_database WHERE datname = '$(DB_NAME)';" | grep -q 1 || \
	sudo -u postgres psql -c "CREATE DATABASE $(DB_NAME);"
	sudo -u postgres psql -c "SELECT 1 FROM pg_user WHERE usename = '$(DB_USER)';" | grep -q 1 || \
	sudo -u postgres psql -c "CREATE USER $(DB_USER) WITH ENCRYPTED PASSWORD '$(DB_PASSWORD)';"
	sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $(DB_NAME) TO $(DB_USER);"
	@echo "$(GREEN)Database $(DB_NAME) and user $(DB_USER) created successfully$(NC)"

# Complete database setup
setup-db: start-postgres create-db
	@echo "$(GREEN)Database setup completed!$(NC)"
	@echo "$(YELLOW)Database credentials:$(NC)"
	@echo "  Host: $(DB_HOST)"
	@echo "  Port: $(DB_PORT)"
	@echo "  Database: $(DB_NAME)"
	@echo "  User: $(DB_USER)"
	@echo "  Password: $(DB_PASSWORD)"

# Install dependencies
install-deps:
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	go mod download
	@echo "$(YELLOW)Installing frontend dependencies...$(NC)"
	cd frontend && npm install
	@echo "$(GREEN)Dependencies installed$(NC)"

# Run in development mode
dev: check-postgres
	@echo "$(YELLOW)Starting application in development mode...$(NC)"
	wails dev

# Build for production
build: check-postgres
	@echo "$(YELLOW)Building application for production...$(NC)"
	wails build
	@echo "$(GREEN)Application built successfully$(NC)"

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf build/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/.cache/
	@echo "$(GREEN)Clean completed$(NC)"

# Test database connection
test-db:
	@echo "$(YELLOW)Testing database connection...$(NC)"
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "SELECT version();"
	@echo "$(GREEN)Database connection successful$(NC)"

# Full setup (for first time users)
full-setup: setup-db install-deps
	@echo "$(GREEN)Full setup completed!$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Create .env file with database credentials"
	@echo "  2. Run 'make dev' to start development server"
	@echo "  3. Run 'make build' to build for production"
