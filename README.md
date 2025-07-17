# BezBase - Fullstack Application

A modern fullstack web application built with Go (Echo framework) backend, React frontend, and PostgreSQL database featuring JWT authentication and comprehensive RBAC (Role-Based Access Control) system.

## 🚀 Current Features (v1.0)

### Core Features
- **Backend**: Go with Echo framework and comprehensive API documentation
- **Frontend**: React with TypeScript, TailwindCSS, and responsive design
- **Database**: PostgreSQL with GORM ORM and automated migrations
- **Authentication**: JWT-based authentication with secure token management
- **Authorization**: Advanced RBAC system with roles, permissions, and resources
- **Internationalization**: Multi-language support (English, Vietnamese) with i18n
- **API Documentation**: Interactive Swagger/OpenAPI documentation
- **Architecture**: Clean separation of handlers, services, and middleware
- **Containerization**: Docker and Docker Compose for development and production
- **Development**: Hot reloading for both frontend and backend

### User Management
- Complete user CRUD operations with status management
- User profiles with extended information (bio, location, timezone)
- Role-based access control with fine-grained permissions
- User authentication with JWT tokens
- Password management and security

### Developer Experience
- Comprehensive API documentation with Swagger/OpenAPI
- Docker containerization for easy deployment
- Hot reloading for rapid development
- Makefile for development task automation
- Clean architecture with proper separation of concerns

## 📋 Prerequisites

- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Go](https://golang.org/doc/install) 1.23+ (for local development)
- [Node.js](https://nodejs.org/) 18+ (for local development)
- [PostgreSQL](https://www.postgresql.org/download/) (for local development)

## 🛠️ Project Structure

```
bezbase/
├── backend/                 # Go backend application
│   ├── cmd/                 # Application entry points
│   │   ├── main.go         # Main application
│   │   └── migrate/        # Database migration tool
│   ├── internal/            # Internal packages
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection & migrations
│   │   ├── docs/           # Swagger documentation setup
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── handlers/       # HTTP request/response handlers
│   │   ├── middleware/     # HTTP middleware (JWT, RBAC)
│   │   ├── models/         # Database models
│   │   ├── pkg/            # Shared packages
│   │   └── services/       # Business logic layer
│   ├── docs/               # Generated Swagger documentation
│   ├── Dockerfile          # Production Docker image
│   ├── Dockerfile.dev      # Development Docker image
│   ├── MIGRATIONS.md       # Database migration documentation
│   ├── RBAC_USAGE.md      # RBAC system usage guide
│   └── go.mod              # Go dependencies
├── frontend/               # React frontend application
│   ├── src/                # Source code
│   │   ├── components/     # React components
│   │   │   ├── common/     # Reusable UI components
│   │   │   └── rbac/       # RBAC-specific components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── context/        # React context (Auth, Theme)
│   │   └── hooks/          # Custom hooks
│   ├── public/             # Public assets
│   ├── Dockerfile          # Production Docker image
│   └── package.json        # Node.js dependencies
├── database/               # Database initialization
│   ├── init.sql           # Database and user creation
│   └── README.md          # Database setup documentation
├── docker/                 # Docker configuration
│   └── nginx.conf         # Nginx reverse proxy config
├── scripts/                # Deployment and setup scripts
└── docker-compose.yml     # Docker Compose configuration
```

## 🛠️ Development with Make

This project includes a Makefile to simplify development tasks. Here are the available commands:

### Initial Setup

```bash
# Install all dependencies (backend and frontend)
make install
```

### Development

```bash
# Start both backend and frontend development servers
make dev

# Start only backend with hot-reload
make dev-backend

# Start only frontend with hot-reload
make dev-frontend
```

### Testing

```bash
# Run all tests
make test

# Run backend tests only
make test-backend

# Run frontend tests only
make test-frontend
```

### Documentation

```bash
# Generate Swagger documentation
make swagger
```

### Building

```bash
# Build both backend and frontend for production
make build

# Clean build artifacts
make clean
```

### Help

```bash
# Show all available commands
make help
```

## 🚀 Quick Start with Docker

### 1. Clone the repository
```bash
git clone <repository-url>
cd bezbase
```

### 2. Start the application
```bash
docker-compose up -d
```

### 3. Access the application
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/
- **Full application (via Nginx)**: http://localhost:80

### 4. Stop the application
```bash
docker-compose down
```

## 🔧 Local Development Setup

### Backend Setup

1. **Navigate to backend directory**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your database configuration.

4. **Run the backend**
   ```bash
   go run cmd/main.go
   ```

### Frontend Setup

1. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```

4. **Run the frontend**
   ```bash
   npm start
   ```

### Database Setup

1. **Create database and user**
   ```bash
   psql -U postgres -f database/init.sql
   ```

2. **Migrations are automatic** - The application will run migrations on startup

For detailed database setup instructions, see [database/README.md](database/README.md).

## 🐳 Docker Development

For development with Docker and live reloading:

```bash
# Start only database and backend with live reloading
docker-compose -f docker-compose.dev.yml up -d

# Run frontend locally
cd frontend
npm install
npm start
```

## 📚 API Documentation

The API is fully documented using Swagger/OpenAPI 3.0. When the backend is running, visit:
- **Interactive Swagger UI**: http://localhost:8080/swagger/

### Main API Groups

#### Authentication Endpoints
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user

#### User Management (Protected)
- `GET /v1/profile` - Get current user profile
- `PUT /v1/profile` - Update current user profile
- `GET /v1/me/permissions` - Get current user permissions
- `GET /v1/users` - List all users (admin only)
- `POST /v1/users` - Create new user (admin only)
- `GET /v1/users/{id}` - Get user by ID (admin only)
- `PUT /v1/users/{id}` - Update user (admin only)
- `DELETE /v1/users/{id}` - Delete user (admin only)

#### RBAC Management (Protected)
- `GET /v1/rbac/roles` - List roles with pagination
- `POST /v1/rbac/roles` - Create new role
- `GET /v1/rbac/roles/{id}` - Get role by ID
- `PUT /v1/rbac/roles/{id}` - Update role
- `DELETE /v1/rbac/roles/{role}` - Delete role
- `GET /v1/rbac/permissions` - List permissions with pagination
- `POST /v1/rbac/permissions` - Add permission to role
- `DELETE /v1/rbac/permissions` - Remove permission from role
- `GET /v1/rbac/resources` - List available resources
- `GET /v1/rbac/actions` - List available actions
- `POST /v1/rbac/users/assign-role` - Assign role to user
- `POST /v1/rbac/users/remove-role` - Remove role from user
- `GET /v1/rbac/users/{id}/roles` - Get user roles
- `GET /v1/rbac/users/{id}/check-permission` - Check user permission

#### System Endpoints
- `GET /health` - Health check

### Authentication

All protected endpoints require a JWT token:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/v1/profile
```

### Example API Usage

**Register a new user:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@example.com",
    "password": "password123"
  }'
```

## 🏗️ Architecture

The backend follows a clean layered architecture pattern:

### Handler Layer (`/handlers`)
- Handles HTTP requests and responses
- Validates request data and formats responses
- Maps service errors to appropriate HTTP status codes
- Keeps business logic separate from HTTP concerns

### Service Layer (`/services`)
- Contains core business logic
- Handles database operations
- Manages authentication and authorization logic
- Independent of HTTP layer for better testability

### Middleware Layer (`/middleware`)
- JWT authentication middleware
- RBAC authorization middleware
- Request logging and error handling

### DTO Layer (`/dto`)
- Data Transfer Objects for API requests/responses
- Input validation and serialization
- Clean separation between API and database models

### Benefits
- **Separation of Concerns**: HTTP handling vs business logic
- **Testability**: Services can be unit tested independently
- **Reusability**: Business logic can be used by different handlers
- **Maintainability**: Clear boundaries between layers

## 🔐 Authentication & Authorization

### JWT Authentication
- Tokens are generated on login/register
- Tokens expire after 24 hours
- Protected routes require `Authorization: Bearer <token>` header
- Frontend automatically handles token storage and API requests

### RBAC (Role-Based Access Control)
- **Roles**: Admin, Moderator, User (with custom roles support)
- **Resources**: Users, Posts, Profile, Admin, Permissions, All
- **Actions**: Create, Read, Update, Delete, All
- **Permissions**: Fine-grained access control per role
- **Middleware**: Automatic permission checking on protected routes

For detailed RBAC usage, see [backend/RBAC_USAGE.md](backend/RBAC_USAGE.md).

## 🌐 Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://bezbase_user:bezbase_password@localhost/bezbase?sslmode=disable
JWT_SECRET=your-secret-key-change-this-in-production
PORT=8080
ENVIRONMENT=development
```

### Frontend (.env)
```
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENV=development
```

## 🧪 Testing

### Backend Testing
```bash
cd backend
go test ./...
```

### Frontend Testing
```bash
cd frontend
npm test
```

## 📦 Building for Production

### Docker Production Build
```bash
docker-compose up -d --build
```

### Manual Production Build

**Backend:**
```bash
cd backend
go build -o main cmd/main.go
```

**Frontend:**
```bash
cd frontend
npm run build
```

## 🔧 Development Tools

- **Air**: Live reloading for Go backend
- **React Scripts**: Development server with hot reloading
- **Docker Compose**: Container orchestration
- **PostgreSQL**: Database with GORM ORM
- **GORM**: Auto-migrations and database operations
- **Swagger**: API documentation generation
- **TailwindCSS**: Utility-first CSS framework
- **Echo**: High performance Go web framework

## 📁 Key Files

- `docker-compose.yml` - Production Docker configuration
- `docker-compose.dev.yml` - Development Docker configuration
- `backend/cmd/main.go` - Backend application entry point
- `frontend/src/App.js` - Frontend application entry point
- `database/init.sql` - Database initialization script
- `backend/docs/` - Generated API documentation
- `Makefile` - Development task automation

## 🚀 Deployment

### Docker Deployment
1. Build and push images to your registry
2. Update environment variables for production
3. Run with `docker-compose up -d`

### Manual Deployment
1. Build backend binary
2. Build frontend static files
3. Setup PostgreSQL database
4. Configure reverse proxy (Nginx)
5. Setup SSL certificates

## 🗺️ Feature Roadmap

BezBase is designed to be a comprehensive foundation for modern web applications. The following features are planned for future releases:

### 🔥 Priority 1: Essential Features (v1.1-1.2)
- **Enhanced Security**: 2FA, OAuth integration, password reset, email verification
- **Testing Infrastructure**: Comprehensive unit, integration, and e2e testing
- **API Enhancements**: Rate limiting, API versioning, enhanced error handling

### 🔮 Priority 2: User Experience (v1.3-1.4)
- **File Management**: Upload system, media library, cloud storage integration
- **Notifications**: Real-time notifications, email system, push notifications
- **Enhanced Profiles**: Extended user profiles, activity logs, profile pictures

### 🌟 Priority 3: Advanced Features (v1.5-1.6)
- **Content Management**: Blog system, rich text editor, content approval
- **Communication**: In-app messaging, email templates, notification center
- **Analytics**: Dashboard analytics, performance monitoring, user insights

### 🚀 Priority 4: Scalability (v1.7-1.8)
- **Performance**: Redis caching, CDN integration, background job processing
- **Architecture**: Event-driven design, microservices preparation
- **Security**: Advanced security headers, audit logging, compliance features

### 🏗️ Priority 5: Enterprise (v2.0+)
- **Multi-tenancy**: Organization-based isolation, tenant management
- **Advanced RBAC**: Hierarchical roles, context-based permissions
- **Compliance**: GDPR compliance, audit systems, data retention
- **Integrations**: Webhook system, third-party APIs, payment gateways

For detailed feature specifications and implementation guidelines, see [FEATURE_ROADMAP.md](FEATURE_ROADMAP.md).

## 💡 Use Cases

BezBase is perfect for:

### Startups & MVPs
- **User Management Systems**: Customer portals, admin dashboards
- **SaaS Applications**: Multi-tenant software with user authentication
- **Internal Tools**: Employee management, project tracking systems

### Educational Projects
- **Learning Fullstack Development**: Modern Go + React architecture
- **Authentication Systems**: JWT and RBAC implementation examples
- **API Development**: RESTful API design with comprehensive documentation

### Enterprise Applications
- **Admin Panels**: Role-based administrative interfaces
- **Customer Portals**: User self-service applications
- **Internal Applications**: HR systems, project management tools

## 🎯 Why Choose BezBase?

- **Production Ready**: Comprehensive authentication, authorization, and security
- **Developer Friendly**: Clean architecture, extensive documentation, hot reloading
- **Scalable**: Designed for growth with microservices-ready architecture
- **Modern Stack**: Latest versions of Go, React, PostgreSQL, and Docker
- **International**: Built-in i18n support for global applications
- **Extensible**: Clean interfaces and patterns for easy feature addition

## 🤝 Contributing

We welcome contributions! Here's how to get involved:

### Development Process
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with proper tests
4. Update documentation as needed
5. Submit a pull request with detailed description

### Contribution Guidelines
- Follow the existing code style and patterns
- Add tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting
- Reference relevant issues in your PR description

### Areas for Contribution
- **Feature Development**: Implement features from the roadmap
- **Testing**: Add unit, integration, and e2e tests
- **Documentation**: Improve guides and API documentation
- **Bug Fixes**: Fix reported issues
- **Performance**: Optimize queries and API responses

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Check if PostgreSQL is running
   - Verify database credentials in `.env`
   - Ensure database exists

2. **Frontend can't connect to backend**
   - Check if backend is running on port 8080
   - Verify `REACT_APP_API_URL` in frontend `.env`
   - Check CORS configuration

3. **Docker build fails**
   - Ensure Docker is running
   - Check Dockerfile syntax
   - Verify all required files exist

4. **RBAC permissions not working**
   - Check user roles assignment
   - Verify permission configuration
   - Review middleware logs

### Getting Help

If you encounter issues:
1. Check the logs: `docker-compose logs [service-name]`
2. Verify environment variables
3. Check database connectivity
4. Review API endpoints and authentication
5. Consult the API documentation at `/swagger/`

---

**Happy coding! 🎉**