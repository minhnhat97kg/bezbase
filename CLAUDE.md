# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

BezBase is a fullstack application with Go (Echo framework) backend, React TypeScript frontend, and PostgreSQL database. It features JWT authentication and comprehensive RBAC (Role-Based Access Control) system with internationalization support.

## Development Commands

### Backend Development
```bash
# Backend development with hot-reload
make dev-backend
cd backend && air

# Run backend tests
make test-backend
cd backend && go test ./... -v

# Build backend
make build-backend
cd backend && go build -o main cmd/main.go

# Generate Swagger docs
make swagger
cd backend && swag init -g cmd/main.go --parseDependency --parseInternal
```

### Frontend Development
```bash
# Frontend development with hot-reload
make dev-frontend
cd frontend && npm start

# Run frontend tests
make test-frontend
cd frontend && npm test

# Build frontend
make build-frontend
cd frontend && npm run build
```

### Database
- Migrations run automatically on startup
- Database initialization: `psql -U postgres -f database/init.sql`
- Connection: PostgreSQL with GORM ORM

## Architecture

### Backend Architecture (Go + Echo)
- **Clean layered architecture** with proper separation of concerns
- **Entry point**: `backend/cmd/main.go`
- **Handlers** (`/handlers`): HTTP request/response handling and validation
- **Services** (`/services`): Business logic layer, database operations
- **Repository** (`/repository`): Data access layer with interfaces
- **Models** (`/models`): Database models (User, Role, Permission, Organization, etc.)
- **DTOs** (`/dto`): Data Transfer Objects for API requests/responses
- **Middleware** (`/middleware`): JWT auth, RBAC, i18n, rate limiting, versioning
- **Config** (`/config`): Configuration management

### Frontend Architecture (React + TypeScript)
- **Entry point**: `frontend/src/App.tsx`
- **Context providers**: AuthContext, OrganizationContext, ThemeContext
- **Services** (`/services`): API communication layer with axios
- **Components**: Reusable UI components with RBAC-specific components
- **Pages**: Route-based page components
- **Hooks**: Custom hooks for auth, resource actions
- **i18n**: Multi-language support (English, Vietnamese)

### Key Models and Relationships
- **User** with **UserInfo** (extended profile data)
- **Role** with **Permission** (resource + action based)
- **Organization** with multi-tenancy support
- **AuthProvider** for authentication methods
- Advanced RBAC with contextual permissions and role inheritance

### API Structure
- **Authentication**: `/auth/register`, `/auth/login`
- **Protected routes**: `/v1/*` with JWT middleware
- **API versioning**: Hardcoded to v1, managed via middleware
- **RBAC endpoints**: `/v1/rbac/*` for role/permission management
- **User management**: `/v1/users/*` and `/v1/profile`

## RBAC System

The application uses a sophisticated RBAC system:
- **Resources**: Users, Posts, Profile, Admin, Permissions, All
- **Actions**: Create, Read, Update, Delete, All
- **Roles**: Admin, Moderator, User (with custom roles support)
- **Context-aware permissions** for organizational boundaries
- **Role inheritance** and **contextual permissions**

## Technology Stack

### Backend
- Go 1.23+ with Echo framework
- GORM ORM with PostgreSQL
- JWT authentication with golang-jwt/jwt/v5
- Casbin for RBAC enforcement
- Swagger/OpenAPI documentation
- Air for hot-reloading in development

### Frontend
- React 18+ with TypeScript
- TailwindCSS for styling
- axios for API communication
- react-i18next for internationalization
- React Router for navigation

## Environment Setup

### Required Tools
- Docker and Docker Compose
- Go 1.23+ (for local development)
- Node.js 18+ (for local development)
- PostgreSQL (for local development)

### Development Workflow
1. Use `make install` for initial setup
2. Use `make dev` to start both backend and frontend
3. Use `make test` to run all tests
4. Backend runs on port 8080, Frontend on port 3000
5. API documentation available at `/swagger/`

## Coding Conventions

### Backend Method Signatures
- **Services and Repository methods**: Always receive `contextx.Contextx` as the first parameter
- **Pattern**: `func (s *Service) MethodName(ctx contextx.Contextx, ...other params) (result, error)`
- **Example**: `func (s *AuthService) Register(ctx contextx.Contextx, req dto.RegisterRequest) (*dto.AuthResponse, error)`
- **Purpose**: Enables context-aware operations including tenant isolation, user context, and tracing

## Important Notes

- **Multi-language support**: English and Vietnamese with i18n middleware
- **Database migrations**: Automatic on startup via GORM
- **Hot reloading**: Available for both backend (Air) and frontend
- **API versioning**: Currently hardcoded to v1
- **Authentication**: JWT tokens with 24-hour expiration
- **CORS**: Configured for cross-origin requests
- **Error handling**: Internationalized error messages