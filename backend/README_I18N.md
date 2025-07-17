# i18n (Internationalization) Implementation

This backend now supports multi-language functionality using the `go-i18n` library.

## Features

- **Multi-language support**: English (en) and Vietnamese (vi) translations
- **Auto-detection**: Language detection from Accept-Language header, query parameters, or custom headers
- **Template support**: String templating with variable substitution
- **Context-aware**: Translations are context-aware and thread-safe
- **Fallback**: Automatic fallback to default language (English) if translation not found

## Language Detection Priority

1. `?lang=` query parameter
2. `X-Language` custom header
3. `Accept-Language` header
4. Default language (English)

## Usage Examples

### Basic Usage in Handlers

```go
func (h *Handler) SomeEndpoint(c echo.Context) error {
    t := i18n.NewTranslator(c.Request().Context())
    
    // Simple translation
    message := t.Error("invalid_request_body")
    
    // Helper methods
    message := t.InvalidRequestBody()
    message := t.InvalidCredentials()
    message := t.UserNotFound()
    
    return echo.NewHTTPError(http.StatusBadRequest, message)
}
```

### Template Data

```go
func (h *Handler) SomeEndpoint(c echo.Context) error {
    t := i18n.NewTranslator(c.Request().Context())
    
    // With template data
    message := t.InsufficientPermissions("admin")
    // Returns: "Insufficient permissions: admin" (EN) or "Không đủ quyền: admin" (VI)
    
    return echo.NewHTTPError(http.StatusForbidden, message)
}
```

### Testing Different Languages

```bash
# Test English (default)
curl -H "Accept-Language: en" http://localhost:8080/api/health

# Test Vietnamese
curl -H "Accept-Language: vi" http://localhost:8080/api/health

# Test with query parameter
curl http://localhost:8080/api/health?lang=vi

# Test with custom header
curl -H "X-Language: vi" http://localhost:8080/api/health
```

## File Structure

```
internal/
├── i18n/
│   ├── locales/
│   │   ├── en.json          # English translations
│   │   └── vi.json          # Vietnamese translations
│   ├── i18n.go              # Core i18n functionality
│   ├── translator.go        # Translator helper with convenience methods
│   └── test_i18n.go         # Test functionality
├── middleware/
│   └── i18n.go              # i18n middleware for language detection
```

## Translation Categories

### Errors (`errors.*`)
- `invalid_request_body`: Invalid request body
- `invalid_credentials`: Invalid credentials
- `username_already_taken`: Username already taken
- `user_not_found`: User not found
- `insufficient_permissions`: Insufficient permissions
- And many more...

### Success Messages (`success.*`)
- `password_changed_successfully`: Password changed successfully
- `user_deleted_successfully`: User deleted successfully
- `role_assigned_successfully`: Role assigned successfully

### Status Messages (`status.*`)
- `healthy`: Server health status
- `server_running`: Server running message
- `server_starting`: Server starting message

### Warnings (`warnings.*`)
- `failed_to_initialize_default_roles`: Warning about role initialization

### Defaults (`defaults.*`)
- `language`: Default language code
- `timezone`: Default timezone

## Adding New Languages

1. Create a new JSON file in `internal/i18n/locales/` (e.g., `fr.json`)
2. Add translations for all keys from `en.json`
3. Update `GetSupportedLanguages()` in `i18n.go`
4. Update `loadLanguageFile()` calls in `Initialize()` function

## Adding New Translation Keys

1. Add the key to all language files (`en.json`, `vi.json`, etc.)
2. Use the translation in your code:
   ```go
   t := i18n.NewTranslator(c.Request().Context())
   message := t.Error("your_new_key")
   ```

## Translation Key Naming Convention

- Use lowercase with underscores: `invalid_request_body`
- Group by category: `errors.`, `success.`, `status.`, `warnings.`, `defaults.`
- Be descriptive but concise
- Use consistent naming patterns

## Template Data Format

For messages with variables, use Go template syntax:

```json
{
  "errors": {
    "insufficient_permissions": "Insufficient permissions: {{.Permission}}"
  }
}
```

Usage:
```go
t.Error("insufficient_permissions", map[string]interface{}{
    "Permission": "admin",
})
```

## Performance Considerations

- Translations are loaded once at startup and cached in memory
- Context-based localizers are lightweight and created per request
- No database queries needed for translations
- Embedded file system for translation files (no external file dependencies)

## Testing

Run the i18n test to verify functionality:

```bash
go run cmd/test_i18n/main.go
```

## Error Handling

- If a translation key is not found, the key itself is returned
- If translation fails, a fallback message is used
- All errors are logged for debugging
- The application continues to function even with translation errors