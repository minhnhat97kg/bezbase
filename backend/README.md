# BezBase Backend

A robust Go backend built with Echo framework, featuring JWT authentication, comprehensive RBAC system, and interactive API documentation.

## 🚀 Features

- **Framework**: Echo - High performance, minimalist Go web framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT-based secure authentication
- **Authorization**: Advanced RBAC with Casbin integration
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Migrations**: Automated database migrations with gormigrate
- **Architecture**: Clean layered architecture (Handlers → Services → Models)
- **Containerization**: Docker support for development and production
- **Hot Reloading**: Air for live reloading during development

## 📁 Project Structure

```
backend/
├── cmd/                     # Application entry points
│   ├── main.go             # Main application
│   └── migrate/            # Database migration CLI tool
│       └── main.go
├── internal/               # Internal packages (not exported)
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── database/          # Database connection & migrations
│   │   ├── database.go    # Connection setup
│   │   └── migrations/    # Migration definitions
│   ├── docs/              # Swagger documentation setup
│   │   └── swagger.go
│   ├── dto/               # Data Transfer Objects
│   │   ├── auth.go        # Authentication DTOs
│   │   ├── user.go        # User management DTOs
│   │   ├── role.go        # Role management DTOs
│   │   ├── permission.go  # Permission management DTOs
│   │   ├── rbac.go        # RBAC resource/action DTOs
│   │   └── pagination.go  # Pagination utilities
│   ├── handlers/          # HTTP request handlers
│   │   ├── auth.go        # Authentication endpoints
│   │   ├── user.go        # User management endpoints
│   │   ├── rbac.go        # RBAC management endpoints
│   │   └── common.go      # Health check and common endpoints
│   ├── middleware/        # HTTP middleware
│   │   ├── jwt.go         # JWT authentication middleware
│   │   └── rbac.go        # RBAC authorization middleware
│   ├── models/            # Database models
│   │   ├── user.go        # User entity
│   │   ├── user_info.go   # User profile information
│   │   ├── role.go        # Role entity
│   │   ├── rule.go        # Casbin rules (permissions)
│   │   └── auth_provider.go # Authentication provider info
│   ├── pkg/               # Shared packages
│   │   └── auth/          # Authentication utilities
│   └── services/          # Business logic layer
│       ├── auth.go        # Authentication business logic
│       ├── user.go        # User management business logic
│       └── rbac.go        # RBAC business logic
├── docs/                  # Generated API documentation
│   ├── docs.go           # Generated Swagger code
│   ├── swagger.json      # Swagger specification (JSON)
│   └── swagger.yaml      # Swagger specification (YAML)
├── tmp/                   # Temporary files (Air hot reloading)
├── vendor/                # Go module dependencies (optional)
├── .env.example           # Environment variables template
├── .air.toml             # Air configuration for hot reloading
├── Dockerfile            # Production container image
├── Dockerfile.dev        # Development container image
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── MIGRATIONS.md         # Database migration documentation
└── RBAC_USAGE.md        # RBAC system usage guide
```

## 🛠️ Quick Start

### Prerequisites

- Go 1.23+ (for local development)
- PostgreSQL 13+ (for local development)
- Docker & Docker Compose (for containerized development)

### 1. Local Development Setup

```bash
# Clone the repository
git clone <repository-url>
cd bezbase/backend

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
nano .env

# Install dependencies
go mod download

# Run the application
go run cmd/main.go
```

### 2. Docker Development

```bash
# Start with Docker Compose
docker-compose up -d

# Or start only the backend service
docker-compose up backend
```

### 3. Access the Application

- **API Server**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/
- **Health Check**: http://localhost:8080/api/health

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the backend directory:

```bash
# Database Configuration
DATABASE_URL=postgres://bezbase_user:bezbase_password@localhost/bezbase?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-change-this-in-production

# Server Configuration
PORT=8080
ENVIRONMENT=development

# Optional: Casbin Configuration
CASBIN_MODEL_PATH=./configs/rbac_model.conf
```

### Database Setup

1. **Create database and user** (if not using Docker):
   ```bash
   psql -U postgres -f ../database/init.sql
   ```

2. **Migrations run automatically** when the application starts

For detailed migration information, see [MIGRATIONS.md](MIGRATIONS.md).

## 📚 API Documentation

The API is fully documented using Swagger/OpenAPI 3.0. When the backend is running:

- **Interactive Documentation**: http://localhost:8080/swagger/
- **JSON Specification**: http://localhost:8080/docs/swagger.json
- **YAML Specification**: http://localhost:8080/docs/swagger.yaml

### API Groups

#### Authentication (`/auth`)
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

#### User Management (`/v1/users`) - Protected
- `GET /v1/profile` - Get current user profile
- `PUT /v1/profile` - Update current user profile
- `GET /v1/me/permissions` - Get current user permissions
- `GET /v1/users` - List all users (admin)
- `POST /v1/users` - Create user (admin)
- `GET /v1/users/{id}` - Get user by ID (admin)
- `PUT /v1/users/{id}` - Update user (admin)
- `DELETE /v1/users/{id}` - Delete user (admin)

#### RBAC Management (`/v1/rbac`) - Protected
- `GET /v1/rbac/roles` - List roles with pagination
- `POST /v1/rbac/roles` - Create new role
- `GET /v1/rbac/roles/{id}` - Get role details
- `PUT /v1/rbac/roles/{id}` - Update role
- `DELETE /v1/rbac/roles/{role}` - Delete role
- `GET /v1/rbac/permissions` - List permissions
- `POST /v1/rbac/permissions` - Add permission
- `DELETE /v1/rbac/permissions` - Remove permission
- `GET /v1/rbac/resources` - List available resources
- `GET /v1/rbac/actions` - List available actions
- User role assignments and permission checks

#### System (`/api`)
- `GET /api/health` - Health check endpoint

## 🔐 Authentication & Authorization

### JWT Authentication

The backend uses JWT tokens for authentication:

```go
// Token structure
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

**Usage:**
```bash
# Include in request headers
Authorization: Bearer <jwt-token>
```

### RBAC (Role-Based Access Control)

The system implements comprehensive RBAC using Casbin:

#### Default Roles
- **admin**: Full system access
- **moderator**: User and content management
- **user**: Basic user permissions

#### Resources
- **users**: User management
- **posts**: Content management
- **profile**: User profile access
- **admin**: Administrative functions
- **permissions**: Role and permission management
- **all**: Global access

#### Actions
- **create**: Create new entities
- **read**: View existing entities
- **update**: Modify existing entities
- **delete**: Remove entities
- **all**: All actions

For detailed RBAC usage, see [RBAC_USAGE.md](RBAC_USAGE.md).

## 🏗️ Architecture

### Clean Architecture Layers

#### 1. Handler Layer (`/handlers`)
```go
// Example handler structure
type UserHandler struct {
    userService *services.UserService
    rbacService *services.RBACService
}

func (h *UserHandler) GetProfile(c echo.Context) error {
    // Handle HTTP request/response
    // Delegate business logic to service layer
}
```

#### 2. Service Layer (`/services`)
```go
// Example service structure
type UserService struct {
    db          *gorm.DB
    rbacService *RBACService
}

func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
    // Business logic implementation
    // Database operations
    // Return domain objects
}
```

#### 3. Model Layer (`/models`)
```go
// Example model structure
type User struct {
    ID            uint `gorm:"primarykey"`
    Status        UserStatus
    EmailVerified bool
    LastLoginAt   *time.Time
    UserInfo      *UserInfo `gorm:"foreignKey:UserID"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     gorm.DeletedAt
}
```

#### 4. DTO Layer (`/dto`)
```go
// Example DTO structure
type UserResponse struct {
    ID        uint      `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Email     string    `json:"email"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    Roles     []string  `json:"roles,omitempty"`
}
```

### Benefits
- **Separation of Concerns**: Clear boundaries between layers
- **Testability**: Each layer can be unit tested independently
- **Maintainability**: Easy to modify and extend
- **Reusability**: Business logic can be reused across different handlers

## 🧪 Development

### Hot Reloading

The project uses [Air](https://github.com/cosmtrek/air) for hot reloading:

```bash
# Install Air (if not already installed)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air

# Or use the Docker development setup
docker-compose -f docker-compose.dev.yml up
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/services/...
```

### Database Migrations

```bash
# Create migration binary
go build -o migrate-tool ./cmd/migrate/main.go

# Run migrations
./migrate-tool up

# Rollback migrations
./migrate-tool down

# Check migration status
./migrate-tool status
```

### Generate Swagger Documentation

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/main.go -o docs/
```

## 🚀 Deployment

### Docker Production

```bash
# Build production image
docker build -t bezbase-backend .

# Run production container
docker run -p 8080:8080 \
  -e DATABASE_URL="your-database-url" \
  -e JWT_SECRET="your-jwt-secret" \
  bezbase-backend
```

### Manual Production Build

```bash
# Build binary
go build -o main cmd/main.go

# Run binary
./main
```

### Environment-Specific Configuration

**Development (.env.development):**
```bash
DATABASE_URL=postgres://bezbase_user:bezbase_password@localhost/bezbase_dev?sslmode=disable
JWT_SECRET=development-secret-key
PORT=8080
ENVIRONMENT=development
```

**Production (.env.production):**
```bash
DATABASE_URL=postgres://user:password@prod-db:5432/bezbase?sslmode=require
JWT_SECRET=super-secure-production-secret-key
PORT=8080
ENVIRONMENT=production
```

## 🔧 Tools & Dependencies

### Core Dependencies
- **Echo v4**: Web framework
- **GORM**: ORM for database operations
- **Casbin**: Authorization library
- **JWT-Go**: JWT token handling
- **PostgreSQL Driver**: Database connectivity
- **Swaggo**: API documentation generation

### Development Dependencies
- **Air**: Hot reloading
- **Testify**: Testing utilities
- **Gormigrate**: Database migrations

### Build Tools
- **Go 1.23+**: Programming language
- **Docker**: Containerization
- **Make**: Build automation (if using Makefile)

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check if PostgreSQL is running
   pg_isready -h localhost -p 5432
   
   # Verify credentials in .env file
   cat .env | grep DATABASE_URL
   ```

2. **Migration Errors**
   ```bash
   # Check migration status
   ./migrate-tool status
   
   # Reset database (caution: destroys data)
   ./migrate-tool reset
   ```

3. **JWT Token Issues**
   ```bash
   # Verify JWT secret is set
   echo $JWT_SECRET
   
   # Check token in request headers
   curl -H "Authorization: Bearer <token>" http://localhost:8080/v1/profile
   ```

4. **RBAC Permission Denied**
   ```bash
   # Check user roles
   curl http://localhost:8080/v1/rbac/users/{user_id}/roles
   
   # Verify permissions
   curl "http://localhost:8080/v1/rbac/users/{user_id}/check-permission?resource=users&action=read"
   ```

### Debugging

```bash
# Enable debug logging
export ENVIRONMENT=development

# View application logs
tail -f server.log

# Check database connections
# Add to .env: DB_LOG_LEVEL=info
```

### Performance Monitoring

```bash
# Enable pprof (development only)
import _ "net/http/pprof"

# Access profiling data
http://localhost:8080/debug/pprof/
```

## 📖 Additional Documentation

- [MIGRATIONS.md](MIGRATIONS.md) - Database migration system
- [RBAC_USAGE.md](RBAC_USAGE.md) - RBAC implementation guide
- [API Documentation](http://localhost:8080/swagger/) - Interactive API docs

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Add tests for new functionality
5. Update documentation
6. Run tests: `go test ./...`
7. Generate updated Swagger docs: `swag init -g cmd/main.go -o docs/`
8. Submit a pull request

### Code Style

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions small and focused
- Write tests for new functionality

---

For questions or support, please refer to the main [project README](../README.md) or create an issue in the repository.