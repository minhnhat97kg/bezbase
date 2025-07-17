package middleware

import (
	"fmt"
	"net/http"

	"bezbase/internal/models"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

func RBACMiddleware(rbacService *services.RBACService, permission models.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user claims from JWT middleware
			userClaims, ok := c.Get("user").(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
			}

			// Check permission
			allowed, err := rbacService.CheckPermission(userClaims.UserID, permission.Resource.String(), permission.Action.String())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Permission check failed: %v", err))
			}

			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Insufficient permissions: %s", permission.Permission))
			}

			return next(c)
		}
	}
}

func RequirePermission(rbacService *services.RBACService, permission models.Permission) echo.MiddlewareFunc {
	return RBACMiddleware(rbacService, permission)
}

// RequirePermissionModel accepts a Permission model directly
func RequirePermissionModel(rbacService *services.RBACService, permission models.Permission) echo.MiddlewareFunc {
	return RBACMiddleware(rbacService, permission)
}

func RequireRole(rbacService *services.RBACService, role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
			}

			// Check if role exists and is active
			roleModel, err := rbacService.GetRoleByName(role)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Role validation failed: %v", err))
			}

			if !roleModel.IsActive {
				return echo.NewHTTPError(http.StatusForbidden, "Role is not active")
			}

			userRoles, err := rbacService.GetUserRoles(userClaims.UserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Role check failed: %v", err))
			}

			hasRole := false
			for _, userRole := range userRoles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient role")
			}

			return next(c)
		}
	}
}

func RequireAnyRole(rbacService *services.RBACService, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
			}

			userRoles, err := rbacService.GetUserRoles(userClaims.UserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Role check failed: %v", err))
			}

			hasAnyRole := false
			for _, userRole := range userRoles {
				for _, requiredRole := range roles {
					if userRole == requiredRole {
						hasAnyRole = true
						break
					}
				}
				if hasAnyRole {
					break
				}
			}

			if !hasAnyRole {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient role")
			}

			return next(c)
		}
	}
}

func RequireAllRoles(rbacService *services.RBACService, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userClaims, ok := c.Get("user").(*auth.Claims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
			}

			userRoles, err := rbacService.GetUserRoles(userClaims.UserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Role check failed: %v", err))
			}

			// Check if user has all required roles
			for _, requiredRole := range roles {
				hasRole := false
				for _, userRole := range userRoles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if !hasRole {
					return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Missing required role: %s", requiredRole))
				}
			}

			return next(c)
		}
	}
}

func CheckPermissionForUser(rbacService *services.RBACService, userID uint, resource, action string) (bool, error) {
	return rbacService.CheckPermission(userID, resource, action)
}
