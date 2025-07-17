package middleware

import (
	"context"
	"strings"

	"bezbase/internal/i18n"

	"github.com/labstack/echo/v4"
)

// I18nMiddleware sets up the i18n localizer in the request context
func I18nMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get language from various sources (in order of preference)
			languages := getLanguagesFromRequest(c)
			
			// Create localizer
			localizer := i18n.GetLocalizer(languages...)
			
			// Set localizer in context
			ctx := context.WithValue(c.Request().Context(), i18n.LocalizerContextKey, localizer)
			c.SetRequest(c.Request().WithContext(ctx))
			
			return next(c)
		}
	}
}

// getLanguagesFromRequest extracts languages from the request in order of preference
func getLanguagesFromRequest(c echo.Context) []string {
	var languages []string
	
	// 1. Check for explicit language parameter in query or header
	if lang := c.QueryParam("lang"); lang != "" {
		if i18n.IsLanguageSupported(lang) {
			languages = append(languages, lang)
		}
	}
	
	// 2. Check for custom language header
	if lang := c.Request().Header.Get("X-Language"); lang != "" {
		if i18n.IsLanguageSupported(lang) {
			languages = append(languages, lang)
		}
	}
	
	// 3. Parse Accept-Language header
	acceptLanguage := c.Request().Header.Get("Accept-Language")
	if acceptLanguage != "" {
		acceptedLanguages := i18n.ParseAcceptLanguage(acceptLanguage)
		for _, lang := range acceptedLanguages {
			// Handle language variants (e.g., "en-US" -> "en")
			if strings.Contains(lang, "-") {
				lang = strings.Split(lang, "-")[0]
			}
			if i18n.IsLanguageSupported(lang) {
				languages = append(languages, lang)
			}
		}
	}
	
	// 4. Add default language if no supported language found
	if len(languages) == 0 {
		languages = append(languages, i18n.DefaultLanguage)
	}
	
	return languages
}