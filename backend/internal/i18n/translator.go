package i18n

import (
	"context"
	"fmt"
)

// Translator provides convenient methods for translation
type Translator struct {
	ctx context.Context
}

// NewTranslator creates a new translator with the given context
func NewTranslator(ctx context.Context) *Translator {
	return &Translator{ctx: ctx}
}

// Error translates error messages
func (t *Translator) Error(key string, args ...interface{}) string {
	fullKey := fmt.Sprintf("errors.%s", key)
	if len(args) > 0 {
		return T(t.ctx, fullKey, args[0])
	}
	return T(t.ctx, fullKey)
}

// Success translates success messages
func (t *Translator) Success(key string, args ...interface{}) string {
	fullKey := fmt.Sprintf("success.%s", key)
	if len(args) > 0 {
		return T(t.ctx, fullKey, args[0])
	}
	return T(t.ctx, fullKey)
}

// Status translates status messages
func (t *Translator) Status(key string, args ...interface{}) string {
	fullKey := fmt.Sprintf("status.%s", key)
	if len(args) > 0 {
		return T(t.ctx, fullKey, args[0])
	}
	return T(t.ctx, fullKey)
}

// Warning translates warning messages
func (t *Translator) Warning(key string, args ...interface{}) string {
	fullKey := fmt.Sprintf("warnings.%s", key)
	if len(args) > 0 {
		return T(t.ctx, fullKey, args[0])
	}
	return T(t.ctx, fullKey)
}

// Default translates default values
func (t *Translator) Default(key string, args ...interface{}) string {
	fullKey := fmt.Sprintf("defaults.%s", key)
	if len(args) > 0 {
		return T(t.ctx, fullKey, args[0])
	}
	return T(t.ctx, fullKey)
}

// Trans translates any message by full key
func (t *Translator) Trans(key string, args ...interface{}) string {
	if len(args) > 0 {
		return T(t.ctx, key, args[0])
	}
	return T(t.ctx, key)
}

// TransWithDefault translates with a default fallback
func (t *Translator) TransWithDefault(key string, defaultMessage string, args ...interface{}) string {
	if len(args) > 0 {
		return TWithDefault(t.ctx, key, defaultMessage, args[0])
	}
	return TWithDefault(t.ctx, key, defaultMessage)
}

// Helper functions for common error patterns

// InvalidRequestBody returns the translated invalid request body message
func (t *Translator) InvalidRequestBody() string {
	return t.Error("invalid_request_body")
}

// InvalidCredentials returns the translated invalid credentials message
func (t *Translator) InvalidCredentials() string {
	return t.Error("invalid_credentials")
}

// UserNotFound returns the translated user not found message
func (t *Translator) UserNotFound() string {
	return t.Error("user_not_found")
}

// UsernameAlreadyTaken returns the translated username already taken message
func (t *Translator) UsernameAlreadyTaken() string {
	return t.Error("username_already_taken")
}

// EmailAlreadyTaken returns the translated email already taken message
func (t *Translator) EmailAlreadyTaken() string {
	return t.Error("email_already_taken")
}

// InsufficientPermissions returns the translated insufficient permissions message
func (t *Translator) InsufficientPermissions(permission string) string {
	return t.Error("insufficient_permissions", map[string]interface{}{
		"Permission": permission,
	})
}

// RoleNotFound returns the translated role not found message
func (t *Translator) RoleNotFound() string {
	return t.Error("role_not_found")
}

// PasswordChangedSuccessfully returns the translated password changed success message
func (t *Translator) PasswordChangedSuccessfully() string {
	return t.Success("password_changed_successfully")
}

// UserDeletedSuccessfully returns the translated user deleted success message
func (t *Translator) UserDeletedSuccessfully() string {
	return t.Success("user_deleted_successfully")
}