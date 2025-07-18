package main

import (
	"log"

	"bezbase/internal/config"
	"bezbase/internal/database"
	"bezbase/internal/docs"
	"bezbase/internal/handlers"
	"bezbase/internal/i18n"
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
	db, err := database.Connect(cfg.Database.URL)
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

	// Initialize i18n
	if err := i18n.Initialize(); err != nil {
		log.Fatal("Failed to initialize i18n:", err)
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
	e.Use(middleware.I18nMiddleware()) // Add i18n middleware
	e.Use(middleware.Versioning())     // Add versioning middleware

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	userInfoRepo := repository.NewUserInfoRepository(db)
	authProviderRepo := repository.NewAuthProviderRepository(db)
	emailVerificationRepo := repository.NewEmailVerificationRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	ruleRepo := repository.NewRuleRepository(db)

	// Advanced RBAC repositories
	roleTemplateRepo := repository.NewRoleTemplateRepository(db)
	contextualPermissionRepo := repository.NewContextualPermissionRepository(db)
	roleInheritanceRepo := repository.NewRoleInheritanceRepository(db)

	// Organization repositories
	orgRepo := repository.NewOrganizationRepository(db)
	orgUserRepo := repository.NewOrganizationUserRepository(db)
	orgInvitationRepo := repository.NewOrganizationInvitationRepository(db)

	// Initialize services
	rbacService, err := services.NewRBACService(roleRepo, ruleRepo, db)
	if err != nil {
		log.Fatal("Failed to initialize RBAC service:", err)
	}
	emailService := services.NewEmailService(emailVerificationRepo, &cfg.Email, cfg.Server.BaseURL)
	emailVerificationService := services.NewEmailVerificationService(userRepo, emailVerificationRepo, emailService)
	passwordResetService := services.NewPasswordResetService(userRepo, userInfoRepo, authProviderRepo, passwordResetRepo, emailService)
	authService := services.NewAuthService(userRepo, userInfoRepo, authProviderRepo, &cfg.Auth, db)
	userService := services.NewUserService(userRepo, userInfoRepo, authProviderRepo, rbacService, db)
	orgService := services.NewOrganizationService(orgRepo, orgUserRepo, orgInvitationRepo, userRepo, rbacService, emailService, db)

	// Initialize handlers
	commonHandler := handlers.NewCommonHandler()
	rbacHandler := handlers.NewRBACHandler(rbacService)
	advancedRbacHandler := handlers.NewAdvancedRBACHandler(rbacService, roleTemplateRepo, contextualPermissionRepo, roleInheritanceRepo, db)
	orgHandler := handlers.NewOrganizationHandler(orgService)
	userHandler := handlers.NewUserHandler(userService, rbacService)
	authHandler := handlers.NewAuthHandler(authService)
	emailVerificationHandler := handlers.NewEmailVerificationHandler(emailVerificationService)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)

	// Routes

	// Public routes
	api := e.Group("/api")

	// API v1 routes
	apiV1 := api.Group("/v1")
	apiV1.Use(middleware.APIRateLimit()) // Add rate limiting for API endpoints

	// Auth routes (public, but versioned)
	auth := apiV1.Group("/auth")
	auth.Use(middleware.AuthRateLimit()) // Add rate limiting for auth endpoints
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// Email verification routes (public)
	auth.POST("/send-verification-email", emailVerificationHandler.SendVerificationEmail)
	auth.POST("/verify-email", emailVerificationHandler.VerifyEmail)
	auth.GET("/verify-email", emailVerificationHandler.VerifyEmailByToken)
	auth.POST("/resend-verification-email", emailVerificationHandler.ResendVerificationEmail)

	// Password reset routes (public)
	auth.POST("/request-password-reset", passwordResetHandler.RequestPasswordReset)
	auth.POST("/reset-password", passwordResetHandler.ResetPassword)
	auth.POST("/validate-reset-token", passwordResetHandler.ValidateResetToken)
	auth.GET("/validate-reset-token", passwordResetHandler.ValidateResetTokenByParam)

	// Protected routes (add JWT middleware after auth routes)
	apiV1.Use(middleware.JWTMiddleware(cfg.Auth.JWTSecret))

	// Profile routes (users can access their own profile)
	apiV1.GET("/profile", userHandler.GetProfile, middleware.RequirePermission(rbacService, models.PermissionViewProfile))
	apiV1.PUT("/profile", userHandler.UpdateProfile, middleware.RequirePermission(rbacService, models.PermissionEditProfile))
	apiV1.PUT("/profile/password", userHandler.ChangePassword, middleware.RequirePermission(rbacService, models.PermissionEditProfile))
	apiV1.GET("/me/permissions", userHandler.GetCurrentUserPermissions)

	// User management routes (admin only)
	userGroup := apiV1.Group("/users")
	// userGroup.Use(middleware.RequirePermission(rbacService,models.ResourceTypeUser, models.ActionTypeAll))
	userGroup.GET("", userHandler.GetUsers, middleware.RequirePermission(rbacService, models.PermissionViewUsers))
	userGroup.GET("/:id", userHandler.GetUser, middleware.RequirePermission(rbacService, models.PermissionViewUsers))
	userGroup.POST("", userHandler.CreateUser, middleware.RequirePermission(rbacService, models.PermissionCreateUsers))
	userGroup.PUT("/:id", userHandler.UpdateUser, middleware.RequirePermission(rbacService, models.PermissionEditUsers))
	userGroup.DELETE("/:id", userHandler.DeleteUser, middleware.RequirePermission(rbacService, models.PermissionDeleteUsers))

	// RBAC management routes (admin only)
	rbacGroup := apiV1.Group("/rbac")
	// rbacGroup.Use(middleware.RequireRole(rbacService, "admin"))

	// Role management
	rbacGroup.POST("/roles", rbacHandler.CreateRole,
		middleware.RequirePermission(rbacService, models.PermissionCreateRoles))
	rbacGroup.GET("/roles", rbacHandler.GetRoles, middleware.RequirePermission(rbacService, models.PermissionViewRoles))
	rbacGroup.GET("/roles/:role_id", rbacHandler.GetRole, middleware.RequirePermission(rbacService, models.PermissionViewRoles))
	rbacGroup.PUT("/roles/:role_id", rbacHandler.UpdateRole, middleware.RequirePermission(rbacService, models.PermissionEditRoles))
	rbacGroup.DELETE("/roles/:role", rbacHandler.DeleteRole, middleware.RequirePermission(rbacService, models.PermissionDeleteRoles))
	rbacGroup.GET("/roles/:role/users", rbacHandler.GetUsersWithRole, middleware.RequirePermission(rbacService, models.PermissionViewRoles))
	rbacGroup.GET("/roles/:role/permissions", rbacHandler.GetRolePermissions, middleware.RequirePermission(rbacService, models.PermissionViewRoles))

	// User role management
	rbacGroup.POST("/users/assign-role", rbacHandler.AssignRole, middleware.RequirePermission(rbacService, models.PermissionEditPermissions))
	rbacGroup.POST("/users/remove-role", rbacHandler.RemoveRole, middleware.RequirePermission(rbacService, models.PermissionEditPermissions))
	rbacGroup.GET("/users/:user_id/roles", rbacHandler.GetUserRoles, middleware.RequirePermission(rbacService, models.PermissionViewPermissions))

	// Permission management
	rbacGroup.GET("/permissions", rbacHandler.GetPermissions, middleware.RequirePermission(rbacService, models.PermissionViewPermissions))
	rbacGroup.GET("/permissions/available", rbacHandler.GetAvailablePermissions, middleware.RequirePermission(rbacService, models.PermissionViewPermissions))
	rbacGroup.POST("/permissions", rbacHandler.AddPermission, middleware.RequirePermission(rbacService, models.PermissionCreatePermissions))
	rbacGroup.DELETE("/permissions", rbacHandler.RemovePermission, middleware.RequirePermission(rbacService, models.PermissionDeletePermissions))
	rbacGroup.GET("/users/:user_id/check-permission", rbacHandler.CheckPermission)

	// Advanced RBAC endpoints
	// Role templates
	rbacGroup.GET("/role-templates", advancedRbacHandler.GetRoleTemplates, middleware.RequirePermission(rbacService, models.PermissionViewRoles))

	// Role hierarchy and template creation
	rbacGroup.POST("/roles/from-template", advancedRbacHandler.CreateRoleFromTemplate, middleware.RequirePermission(rbacService, models.PermissionCreateRoles))
	rbacGroup.PUT("/roles/:role_id/parent", advancedRbacHandler.SetRoleParent, middleware.RequirePermission(rbacService, models.PermissionEditRoles))
	rbacGroup.GET("/roles/:role_id/hierarchy", advancedRbacHandler.GetRoleHierarchy, middleware.RequirePermission(rbacService, models.PermissionViewRoles))

	// Contextual permissions
	rbacGroup.POST("/contextual-permissions", advancedRbacHandler.CreateContextualPermission, middleware.RequirePermission(rbacService, models.PermissionCreatePermissions))
	rbacGroup.GET("/users/:user_id/effective-permissions", advancedRbacHandler.GetEffectivePermissions, middleware.RequirePermission(rbacService, models.PermissionViewPermissions))

	// Organization management routes (TODO: Implement missing handler methods)
	orgGroup := apiV1.Group("/organizations")

	// Basic organization endpoints that exist
	orgGroup.POST("", orgHandler.CreateOrganization, middleware.RequirePermission(rbacService, models.PermissionCreateUsers)) // Reuse user permission for now
	orgGroup.GET("/:id", orgHandler.GetOrganization, middleware.RequirePermission(rbacService, models.PermissionViewUsers))

	// TODO: Add remaining organization endpoints once handler methods are implemented

	// API v2 (future version example)
	apiV2 := api.Group("/v2")
	apiV2.Use(middleware.APIRateLimit()) // Add rate limiting for API endpoints
	apiV2.Use(middleware.JWTMiddleware(cfg.Auth.JWTSecret))
	apiV2.Use(middleware.RequireMinVersion(2)) // Require minimum version 2

	// Health check
	api.GET("/health", commonHandler.HealthCheck)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(e.Start(":" + cfg.Server.Port))
}
