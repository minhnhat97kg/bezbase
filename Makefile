.PHONY: all install dev test build clean swagger help

# Default target
all: help

# Colors for terminal output
COLOR_RESET = \033[0m
COLOR_BLUE = \033[34m
COLOR_GREEN = \033[32m

# Development server ports
BACKEND_PORT = 8080
FRONTEND_PORT = 3000

# Install all dependencies
install: install-backend install-frontend

install-backend:
	@echo "$(COLOR_BLUE)Installing backend dependencies...$(COLOR_RESET)"
	cd backend && go mod download && go mod tidy
	cd backend && go install github.com/swaggo/swag/cmd/swag@latest
	cd backend && go install github.com/cosmtrek/air@latest

install-frontend:
	@echo "$(COLOR_BLUE)Installing frontend dependencies...$(COLOR_RESET)"
	cd frontend && npm install

# Run development servers
dev: dev-backend dev-frontend

dev-backend:
	@echo "$(COLOR_GREEN)Starting backend development server on port $(BACKEND_PORT)...$(COLOR_RESET)"
	cd backend && air

dev-frontend:
	@echo "$(COLOR_GREEN)Starting frontend development server on port $(FRONTEND_PORT)...$(COLOR_RESET)"
	cd frontend && npm start

# Run tests
test: test-backend test-frontend

test-backend:
	@echo "$(COLOR_BLUE)Running backend tests...$(COLOR_RESET)"
	cd backend && go test ./... -v

test-frontend:
	@echo "$(COLOR_BLUE)Running frontend tests...$(COLOR_RESET)"
	cd frontend && npm test

# Build for production
build: build-backend build-frontend

build-backend:
	@echo "$(COLOR_BLUE)Building backend...$(COLOR_RESET)"
	cd backend && go build -o main cmd/main.go

build-frontend:
	@echo "$(COLOR_BLUE)Building frontend...$(COLOR_RESET)"
	cd frontend && npm run build

# Generate Swagger documentation
swagger:
	@echo "$(COLOR_BLUE)Generating Swagger documentation...$(COLOR_RESET)"
	cd backend && swag init -g cmd/main.go --parseDependency --parseInternal

# Clean build artifacts
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf backend/main backend/tmp
	rm -rf frontend/build
	rm -rf backend/docs

# Help information
help:
	@echo "Available commands:"
	@echo "  make install         - Install all dependencies"
	@echo "  make dev            - Start development servers (backend and frontend)"
	@echo "  make dev-backend    - Start backend development server only"
	@echo "  make dev-frontend   - Start frontend development server only"
	@echo "  make test          - Run all tests"
	@echo "  make test-backend   - Run backend tests only"
	@echo "  make test-frontend  - Run frontend tests only"
	@echo "  make build         - Build for production"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make clean         - Clean build artifacts"
