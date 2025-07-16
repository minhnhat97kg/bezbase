# BezBase - Fullstack Application

A modern fullstack web application built with Go (Echo framework) backend, React frontend, and PostgreSQL database with JWT authentication.

## 🚀 Features

- **Backend**: Go with Echo framework
- **Frontend**: React with JavaScript
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT-based authentication
- **Architecture**: Clean separation of handlers and services
- **Containerization**: Docker and Docker Compose
- **Development**: Hot reloading for both frontend and backend

## 📋 Prerequisites

- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Go](https://golang.org/doc/install) 1.21+ (for local development)
- [Node.js](https://nodejs.org/) 18+ (for local development)
- [PostgreSQL](https://www.postgresql.org/download/) (for local development)

## 🛠️ Project Structure

```
bezbase/
├── backend/                 # Go backend application
│   ├── cmd/                 # Application entry point
│   │   └── main.go         # Main application file
│   ├── internal/            # Internal packages
│   │   ├── auth/           # Authentication utilities
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection & migrations
│   │   ├── handlers/       # HTTP request/response handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Data models & DTOs
│   │   └── services/       # Business logic layer
│   │       └── user_service.go  # User business logic
│   ├── Dockerfile          # Production Docker image
│   ├── Dockerfile.dev      # Development Docker image
│   └── go.mod              # Go dependencies
├── frontend/               # React frontend application
│   ├── src/                # Source code
│   │   ├── components/     # React components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── context/        # React context
│   │   └── hooks/          # Custom hooks
│   ├── public/             # Public assets
│   ├── Dockerfile          # Production Docker image
│   └── package.json        # Node.js dependencies
├── database/               # Database files
│   ├── migrations/         # SQL migration files
│   ├── init.sql           # Database initialization
│   └── README.md          # Database documentation
├── docker/                 # Docker configuration
│   └── nginx.conf         # Nginx configuration
├── scripts/                # Deployment scripts
│   ├── deploy-prod.sh     # Production deployment
│   └── setup-dev.sh       # Development setup
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
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Full application (via Nginx): http://localhost:80

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

2. **Run migrations**
   ```bash
   psql -U bezbase_user -d bezbase -f database/migrations/001_create_users_table.sql
   ```

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

### Authentication Endpoints

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Protected Endpoints

- `GET /api/protected/profile` - Get user profile
- `PUT /api/protected/profile` - Update user profile

### Public Endpoints

- `GET /api/health` - Health check

### Example API Usage

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
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
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Get profile (with JWT token):**
```bash
curl -X GET http://localhost:8080/api/protected/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
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

### Benefits
- **Separation of Concerns**: HTTP handling vs business logic
- **Testability**: Services can be unit tested independently
- **Reusability**: Business logic can be used by different handlers
- **Maintainability**: Clear boundaries between layers

## 🔐 Authentication

The application uses JWT (JSON Web Tokens) for authentication:

- Tokens are generated on login/register
- Tokens expire after 24 hours
- Protected routes require `Authorization: Bearer <token>` header
- Frontend automatically handles token storage and API requests

## 🌐 Environment Variables

### Backend (.env)
```
DATABASE_URL=postgres://user:password@localhost/bezbase?sslmode=disable
JWT_SECRET=your-secret-key-change-this-in-production
PORT=8080
```

### Frontend (.env)
```
REACT_APP_API_URL=http://localhost:8080/api
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

## 📁 Key Files

- `docker-compose.yml` - Production Docker configuration
- `docker-compose.dev.yml` - Development Docker configuration
- `backend/cmd/main.go` - Backend application entry point
- `frontend/src/App.js` - Frontend application entry point
- `database/init.sql` - Database initialization script

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

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

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

### Getting Help

If you encounter issues:
1. Check the logs: `docker-compose logs [service-name]`
2. Verify environment variables
3. Check database connectivity
4. Review API endpoints and authentication

---

**Happy coding! 🎉**# bezbase
