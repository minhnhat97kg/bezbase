package handlers

import (
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/i18n"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type PasswordResetHandler struct {
	passwordResetService *services.PasswordResetService
}

func NewPasswordResetHandler(passwordResetService *services.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
	}
}

// RequestPasswordReset initiates a password reset request
// @Summary Request password reset
// @Description Send a password reset email to the user's registered email address
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body dto.PasswordResetRequest true "Password reset request"
// @Success 200 {object} dto.PasswordResetResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/request-password-reset [post]
func (h *PasswordResetHandler) RequestPasswordReset(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.PasswordResetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.passwordResetService.RequestPasswordReset(ctx, req.Email); err != nil {
		if err.Error() == "password reset was requested recently, please wait before requesting again" {
			return echo.NewHTTPError(http.StatusTooManyRequests, t.Error("password.reset_rate_limited"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("password.reset_request_failed"))
	}

	return c.JSON(http.StatusOK, dto.PasswordResetResponse{
		Message: t.Success("password.reset_email_sent"),
		Success: true,
	})
}

// ResetPassword resets the user's password using the reset token
// @Summary Reset password
// @Description Reset the user's password using the provided reset token
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} dto.PasswordResetResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/reset-password [post]
func (h *PasswordResetHandler) ResetPassword(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.passwordResetService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		if err.Error() == "invalid or expired reset token" {
			return echo.NewHTTPError(http.StatusNotFound, t.Error("password.invalid_token"))
		}
		if err.Error() == "password reset token is expired or already used" {
			return echo.NewHTTPError(http.StatusBadRequest, t.Error("password.token_expired"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("password.reset_failed"))
	}

	return c.JSON(http.StatusOK, dto.PasswordResetResponse{
		Message: t.Success("password.reset_success"),
		Success: true,
	})
}

// ValidateResetToken validates a password reset token
// @Summary Validate reset token
// @Description Validate if a password reset token is valid and not expired
// @Tags Password Reset
// @Accept json
// @Produce json
// @Param request body dto.ValidateResetTokenRequest true "Validate reset token request"
// @Success 200 {object} dto.PasswordResetResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/validate-reset-token [post]
func (h *PasswordResetHandler) ValidateResetToken(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.ValidateResetTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.passwordResetService.ValidateResetToken(ctx, req.Token); err != nil {
		if err.Error() == "invalid or expired reset token" {
			return echo.NewHTTPError(http.StatusNotFound, t.Error("password.invalid_token"))
		}
		if err.Error() == "password reset token is expired or already used" {
			return echo.NewHTTPError(http.StatusBadRequest, t.Error("password.token_expired"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("password.validation_failed"))
	}

	return c.JSON(http.StatusOK, dto.PasswordResetResponse{
		Message: t.Success("password.token_valid"),
		Success: true,
	})
}

// ValidateResetTokenByParam validates a password reset token from URL parameter
// @Summary Validate reset token by parameter
// @Description Validate if a password reset token is valid and not expired using URL parameter
// @Tags Password Reset
// @Produce json
// @Param token query string true "Reset token"
// @Success 200 {object} dto.PasswordResetResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/validate-reset-token/{token} [get]
func (h *PasswordResetHandler) ValidateResetTokenByParam(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	token := c.QueryParam("token")
	if token == "" {
		token = c.Param("token")
	}

	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, t.Error("password.token_required"))
	}

	if err := h.passwordResetService.ValidateResetToken(ctx, token); err != nil {
		if err.Error() == "invalid or expired reset token" {
			return echo.NewHTTPError(http.StatusNotFound, t.Error("password.invalid_token"))
		}
		if err.Error() == "password reset token is expired or already used" {
			return echo.NewHTTPError(http.StatusBadRequest, t.Error("password.token_expired"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("password.validation_failed"))
	}

	return c.JSON(http.StatusOK, dto.PasswordResetResponse{
		Message: t.Success("password.token_valid"),
		Success: true,
	})
}