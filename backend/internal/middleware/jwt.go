package middleware

import (
	"net/http"
	"strings"

	"bezbase/internal/i18n"
	"bezbase/internal/pkg/auth"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := i18n.NewTranslator(c.Request().Context())
			
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, t.Error("missing_authorization_header"))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, t.Error("invalid_authorization_header"))
			}

			claims, err := auth.ValidateToken(tokenString, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, t.Error("invalid_token"))
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
