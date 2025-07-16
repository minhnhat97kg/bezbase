package main

import (
	"bezbase/internal/config"
	"bezbase/internal/database"
	"bezbase/internal/handlers"
	"bezbase/internal/middleware"
	"bezbase/internal/rbac"
	"bezbase/internal/services"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Initialize Casbin RBAC enforcer
	enforcer, err := rbac.NewEnforcer(db)
	if err != nil {
		log.Fatal("Failed to initialize Casbin enforcer:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)

	// Initialize handlers
	h := handlers.NewHandlerWithEnforcer(authService, userService, enforcer)

	// Routes
	setupRoutes(e, h, cfg.JWTSecret)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}

func setupRoutes(e *echo.Echo, h *handlers.Handler, jwtSecret string) {
	// Public routes
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// Protected routes
	apiV1 := api.Group("/v1")
	apiV1.Use(middleware.JWTMiddleware(jwtSecret))
	apiV1.GET("/profile", h.GetProfile)
	apiV1.PUT("/profile", h.UpdateProfile)

	// Health check
	api.GET("/health", h.HealthCheck)
}
