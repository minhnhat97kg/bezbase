package main

import (
	"log"
	"os"

	"bezbase/internal/config"
	"bezbase/internal/database"
	"bezbase/internal/docs"
	"bezbase/internal/handlers"
	"bezbase/internal/middleware"
	"bezbase/internal/services/rbac"

	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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

	// Swagger documentation
	docs.SwaggerInfo()

	// Initialize Echo
	e := echo.New()

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Initialize RBAC service
	rbacService, err := rbac.NewRBACService(db)
	if err != nil {
		log.Fatal("Failed to initialize RBAC service:", err)
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)

	// Initialize handlers
	h := handlers.NewHandler(authService, userService)
	rbacHandler := handlers.NewRBACHandler(rbacService)

	// Routes
	setupRoutes(e, h, rbacHandler, rbacService, cfg.JWTSecret)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}

func setupRoutes(e *echo.Echo, h *handlers.Handler, rbacHandler *handlers.RBACHandler, rbacService *rbac.RBACService, jwtSecret string) {
	// Public routes
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// Protected routes
	apiV1 := api.Group("/v1")
	apiV1.Use(middleware.JWTMiddleware(jwtSecret))

	// Profile routes (users can access their own profile)
	apiV1.GET("/profile", h.GetProfile, middleware.RequirePermission(rbacService, "profile", "read"))
	apiV1.PUT("/profile", h.UpdateProfile, middleware.RequirePermission(rbacService, "profile", "update"))

	// User management routes (admin only)
	userGroup := apiV1.Group("/users")
	userGroup.Use(middleware.RequireRole(rbacService, "admin"))
	userGroup.GET("", h.GetUsers)

	// RBAC management routes (admin only)
	rbacGroup := apiV1.Group("/rbac")
	rbacGroup.Use(middleware.RequireRole(rbacService, "admin"))

	// Role management
	rbacGroup.POST("/roles", rbacHandler.CreateRole)
	rbacGroup.GET("/roles", rbacHandler.GetRoles)
	rbacGroup.GET("/roles/:role_id", rbacHandler.GetRole)
	rbacGroup.PUT("/roles/:role_id", rbacHandler.UpdateRole)
	rbacGroup.DELETE("/roles/:role", rbacHandler.DeleteRole)
	rbacGroup.GET("/roles/:role/users", rbacHandler.GetUsersWithRole)
	rbacGroup.GET("/roles/:role/permissions", rbacHandler.GetRolePermissions)

	// User role management
	rbacGroup.POST("/users/assign-role", rbacHandler.AssignRole)
	rbacGroup.POST("/users/remove-role", rbacHandler.RemoveRole)
	rbacGroup.GET("/users/:user_id/roles", rbacHandler.GetUserRoles)

	// Permission management
	rbacGroup.GET("/permissions", rbacHandler.GetPermissions)
	rbacGroup.POST("/permissions", rbacHandler.AddPermission)
	rbacGroup.DELETE("/permissions", rbacHandler.RemovePermission)
	rbacGroup.GET("/users/:user_id/check-permission", rbacHandler.CheckPermission)

	// Health check
	api.GET("/health", h.HealthCheck)
}
