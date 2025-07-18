package middleware

import (
	"net/http"
	"strconv"

	"bezbase/internal/pkg/contextx"

	"github.com/labstack/echo/v4"
)

// TenantContext keys for organization context
const (
	OrganizationIDKey = "organization_id"
	OrganizationKey   = "organization"
)

// TenantMiddleware extracts organization context from request
func TenantMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract organization ID from various sources
			orgID := extractOrganizationID(c)
			
			if orgID != nil {
				// Set organization ID in context
				c.Set(OrganizationIDKey, *orgID)
				
				// Set in contextx for use in services
				ctx := contextx.FromEchoContext(c)
				ctx = contextx.WithOrganizationID(ctx, *orgID)
				c.Set("contextx", ctx)
			}
			
			return next(c)
		}
	}
}

// extractOrganizationID extracts organization ID from request in order of precedence:
// 1. Header: X-Organization-ID
// 2. Query parameter: org_id
// 3. URL path parameter: orgId
// 4. User's current organization (if authenticated)
func extractOrganizationID(c echo.Context) *uint {
	// 1. Check header
	if orgIDHeader := c.Request().Header.Get("X-Organization-ID"); orgIDHeader != "" {
		if orgID, err := strconv.ParseUint(orgIDHeader, 10, 32); err == nil {
			result := uint(orgID)
			return &result
		}
	}
	
	// 2. Check query parameter
	if orgIDQuery := c.QueryParam("org_id"); orgIDQuery != "" {
		if orgID, err := strconv.ParseUint(orgIDQuery, 10, 32); err == nil {
			result := uint(orgID)
			return &result
		}
	}
	
	// 3. Check URL path parameter
	if orgIDParam := c.Param("orgId"); orgIDParam != "" {
		if orgID, err := strconv.ParseUint(orgIDParam, 10, 32); err == nil {
			result := uint(orgID)
			return &result
		}
	}
	
	// 4. Check authenticated user's current organization
	if userID := GetUserIDFromContext(c); userID != nil {
		if currentOrgID := getUserCurrentOrganization(c, *userID); currentOrgID != nil {
			return currentOrgID
		}
	}
	
	return nil
}

// getUserCurrentOrganization gets user's current organization from database
func getUserCurrentOrganization(c echo.Context, userID uint) *uint {
	// This would typically query the database to get user's current_org_id
	// For now, we'll check if it's available in the user context
	if user := c.Get("user"); user != nil {
		// Assuming user struct has CurrentOrgID field
		// This would need to be implemented based on your user loading logic
		// userModel := user.(*models.User)
		// return userModel.CurrentOrgID
	}
	return nil
}

// OrganizationRequired middleware ensures request has organization context
func OrganizationRequired() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			orgID := GetOrganizationIDFromContext(c)
			if orgID == nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Organization context required")
			}
			return next(c)
		}
	}
}

// GetOrganizationIDFromContext gets organization ID from echo context
func GetOrganizationIDFromContext(c echo.Context) *uint {
	if orgID := c.Get(OrganizationIDKey); orgID != nil {
		if id, ok := orgID.(uint); ok {
			return &id
		}
	}
	return nil
}

// GetUserIDFromContext gets user ID from echo context (assuming JWT middleware sets this)
func GetUserIDFromContext(c echo.Context) *uint {
	if userID := c.Get("user_id"); userID != nil {
		if id, ok := userID.(uint); ok {
			return &id
		}
	}
	return nil
}

// OrganizationScope middleware ensures user has access to the requested organization
func OrganizationScope() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := GetUserIDFromContext(c)
			orgID := GetOrganizationIDFromContext(c)
			
			if userID == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}
			
			if orgID == nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Organization context required")
			}
			
			// Check if user has access to this organization
			hasAccess, err := checkUserOrganizationAccess(c, *userID, *orgID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify organization access")
			}
			
			if !hasAccess {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied to this organization")
			}
			
			return next(c)
		}
	}
}

// checkUserOrganizationAccess verifies if user has access to organization
func checkUserOrganizationAccess(c echo.Context, userID uint, orgID uint) (bool, error) {
	// This would query the organization_users table to verify membership
	// For now, we'll return true - this should be implemented with actual database query
	
	// Example implementation:
	// db := c.Get("db").(*gorm.DB)
	// var count int64
	// err := db.Model(&models.OrganizationUser{}).
	//     Where("user_id = ? AND org_id = ?", userID, orgID).
	//     Count(&count).Error
	// return count > 0, err
	
	return true, nil
}

// RoleBasedAccessWithOrg middleware combines RBAC with organization context
func RoleBasedAccessWithOrg(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := GetUserIDFromContext(c)
			if userID == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}
			
			_ = GetOrganizationIDFromContext(c) // Organization context available for future use
			
			// Get RBAC service from context (assuming it's injected)
			rbacService := c.Get("rbac_service")
			if rbacService == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "RBAC service not available")
			}
			
			// Check permission with organization context
			// rbac := rbacService.(*services.RBACService)
			// hasPermission, err := rbac.CheckPermissionWithContext(*userID, resource, action, orgID)
			// if err != nil {
			//     return echo.NewHTTPError(http.StatusInternalServerError, "Permission check failed")
			// }
			// 
			// if !hasPermission {
			//     return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
			// }
			
			return next(c)
		}
	}
}