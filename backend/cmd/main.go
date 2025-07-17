package main

import (
	"log"
	"os"

	"bezbase/internal/config"
	"bezbase/internal/database"
	"bezbase/internal/docs"
	"bezbase/internal/handlers"
	"bezbase/internal/middleware"
	"bezbase/internal/models"
	"bezbase/internal/repository"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	userInfoRepo := repository.NewUserInfoRepository(db)
	authProviderRepo := repository.NewAuthProviderRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	ruleRepo := repository.NewRuleRepository(db)

	// Initialize services
	rbacService, err := services.NewRBACService(roleRepo, ruleRepo, db)
	if err != nil {
		log.Fatal("Failed to initialize RBAC service:", err)
	}
	authService := services.NewAuthService(userRepo, userInfoRepo, authProviderRepo, cfg.JWTSecret, db)
	userService := services.NewUserService(userRepo, userInfoRepo, authProviderRepo, rbacService, db)

	// Initialize handlers
	commonHandler := handlers.NewCommonHandler()
	rbacHandler := handlers.NewRBACHandler(rbacService)
	userHandler := handlers.NewUserHandler(userService, rbacService)
	authHandler := handlers.NewAuthHandler(authService)

	// Routes

	// Public routes
	api := e.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Protected routes
	apiV1 := api.Group("/v1")
	apiV1.Use(middleware.JWTMiddleware(cfg.JWTSecret))

	// Profile routes (users can access their own profile)
	apiV1.GET("/profile", userHandler.GetProfile, middleware.RequirePermission(rbacService, models.ResourceTypeProfile, models.ActionTypeRead))
	apiV1.PUT("/profile", userHandler.UpdateProfile, middleware.RequirePermission(rbacService, models.ResourceTypeProfile, models.ActionTypeUpdate))
	apiV1.PUT("/profile/password", userHandler.ChangePassword, middleware.RequirePermission(rbacService, models.ResourceTypeProfile, models.ActionTypeUpdate))
	apiV1.GET("/me/permissions", userHandler.GetCurrentUserPermissions)

	// User management routes (admin only)
	userGroup := apiV1.Group("/users")
	// userGroup.Use(middleware.RequirePermission(rbacService,models.ResourceTypeUser, models.ActionTypeAll))
	userGroup.GET("", userHandler.GetUsers, middleware.RequirePermission(rbacService, models.ResourceTypeUser, models.ActionTypeRead))
	userGroup.GET("/:id", userHandler.GetUser, middleware.RequirePermission(rbacService, models.ResourceTypeUser, models.ActionTypeRead))
	userGroup.POST("", userHandler.CreateUser, middleware.RequirePermission(rbacService, models.ResourceTypeUser, models.ActionTypeCreate))
	userGroup.PUT("/:id", userHandler.UpdateUser, middleware.RequirePermission(rbacService, models.ResourceTypeUser, models.ActionTypeUpdate))
	userGroup.DELETE("/:id", userHandler.DeleteUser, middleware.RequirePermission(rbacService, models.ResourceTypeUser, models.ActionTypeDelete))

	// RBAC management routes (admin only)
	rbacGroup := apiV1.Group("/rbac")
	// rbacGroup.Use(middleware.RequireRole(rbacService, "admin"))

	// Role management
	rbacGroup.POST("/roles", rbacHandler.CreateRole,
		middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeCreate))
	rbacGroup.GET("/roles", rbacHandler.GetRoles, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))
	rbacGroup.GET("/roles/:role_id", rbacHandler.GetRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))
	rbacGroup.PUT("/roles/:role_id", rbacHandler.UpdateRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeUpdate))
	rbacGroup.DELETE("/roles/:role", rbacHandler.DeleteRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeDelete))
	rbacGroup.GET("/roles/:role/users", rbacHandler.GetUsersWithRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))
	rbacGroup.GET("/roles/:role/permissions", rbacHandler.GetRolePermissions, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))

	// User role management
	rbacGroup.POST("/users/assign-role", rbacHandler.AssignRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeUpdate))
	rbacGroup.POST("/users/remove-role", rbacHandler.RemoveRole, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeUpdate))
	rbacGroup.GET("/users/:user_id/roles", rbacHandler.GetUserRoles, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))

	// Permission management
	rbacGroup.GET("/permissions", rbacHandler.GetPermissions, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))
	rbacGroup.POST("/permissions", rbacHandler.AddPermission, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeCreate))
	rbacGroup.DELETE("/permissions", rbacHandler.RemovePermission, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeDelete))
	rbacGroup.GET("/users/:user_id/check-permission", rbacHandler.CheckPermission)

	// Resource and Action management
	rbacGroup.GET("/resources", rbacHandler.GetResources, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))
	rbacGroup.GET("/actions", rbacHandler.GetActions, middleware.RequirePermission(rbacService, models.ResourceTypePermission, models.ActionTypeRead))

	// Health check
	api.GET("/health", commonHandler.HealthCheck)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}
