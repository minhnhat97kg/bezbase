package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// VersionConfig holds the versioning configuration
type VersionConfig struct {
	DefaultVersion string
	HeaderKey      string
	QueryParam     string
	URLPrefix      string
	MinVersion     int
	MaxVersion     int
}

// DefaultVersionConfig returns a default versioning configuration
func DefaultVersionConfig() VersionConfig {
	return VersionConfig{
		DefaultVersion: "v1",
		HeaderKey:      "API-Version",
		QueryParam:     "version",
		URLPrefix:      "/api",
		MinVersion:     1,
		MaxVersion:     2,
	}
}

// VersionInfo holds version information
type VersionInfo struct {
	Version     string
	VersionNum  int
	IsSupported bool
}

// parseVersion extracts version number from version string
func parseVersion(version string) int {
	re := regexp.MustCompile(`v?(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		if num, err := strconv.Atoi(matches[1]); err == nil {
			return num
		}
	}
	return 1
}

// getVersionFromRequest extracts version from request
func getVersionFromRequest(c echo.Context, config VersionConfig) string {
	// 1. Check URL path
	path := c.Request().URL.Path
	if strings.HasPrefix(path, config.URLPrefix) {
		pathParts := strings.Split(strings.TrimPrefix(path, config.URLPrefix), "/")
		if len(pathParts) > 1 && strings.HasPrefix(pathParts[1], "v") {
			return pathParts[1]
		}
	}
	
	// 2. Check header
	if version := c.Request().Header.Get(config.HeaderKey); version != "" {
		return version
	}
	
	// 3. Check query parameter
	if version := c.QueryParam(config.QueryParam); version != "" {
		return version
	}
	
	// 4. Return default
	return config.DefaultVersion
}

// Versioning returns a middleware that handles API versioning
func Versioning() echo.MiddlewareFunc {
	return VersioningWithConfig(DefaultVersionConfig())
}

// VersioningWithConfig returns a versioning middleware with custom configuration
func VersioningWithConfig(config VersionConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			version := getVersionFromRequest(c, config)
			versionNum := parseVersion(version)
			
			// Create version info
			versionInfo := VersionInfo{
				Version:     version,
				VersionNum:  versionNum,
				IsSupported: versionNum >= config.MinVersion && versionNum <= config.MaxVersion,
			}
			
			// Store version info in context
			c.Set("version", versionInfo)
			
			// Set response headers
			c.Response().Header().Set("API-Version", version)
			c.Response().Header().Set("API-Version-Supported", strconv.FormatBool(versionInfo.IsSupported))
			
			// Check if version is supported
			if !versionInfo.IsSupported {
				return echo.NewHTTPError(http.StatusNotAcceptable, map[string]interface{}{
					"error": "Unsupported API version",
					"version": version,
					"supported_versions": map[string]int{
						"min": config.MinVersion,
						"max": config.MaxVersion,
					},
				})
			}
			
			return next(c)
		}
	}
}

// GetVersionInfo retrieves version information from context
func GetVersionInfo(c echo.Context) *VersionInfo {
	if info, ok := c.Get("version").(VersionInfo); ok {
		return &info
	}
	return nil
}

// RequireVersion returns a middleware that requires a specific version
func RequireVersion(requiredVersion string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			versionInfo := GetVersionInfo(c)
			if versionInfo == nil || versionInfo.Version != requiredVersion {
				return echo.NewHTTPError(http.StatusNotAcceptable, map[string]interface{}{
					"error": "This endpoint requires a specific API version",
					"required_version": requiredVersion,
					"current_version": versionInfo.Version,
				})
			}
			return next(c)
		}
	}
}

// RequireMinVersion returns a middleware that requires a minimum version
func RequireMinVersion(minVersion int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			versionInfo := GetVersionInfo(c)
			if versionInfo == nil || versionInfo.VersionNum < minVersion {
				return echo.NewHTTPError(http.StatusNotAcceptable, map[string]interface{}{
					"error": "This endpoint requires a minimum API version",
					"required_min_version": minVersion,
					"current_version": versionInfo.Version,
				})
			}
			return next(c)
		}
	}
}