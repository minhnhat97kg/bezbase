package handlers

import (
	"net/http"
	"strconv"

	"bezbase/internal/dto"
	"bezbase/internal/i18n"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type EmailVerificationHandler struct {
	verificationService *services.EmailVerificationService
}

func NewEmailVerificationHandler(verificationService *services.EmailVerificationService) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		verificationService: verificationService,
	}
}

// SendVerificationEmail sends a verification email to the user
// @Summary Send verification email
// @Description Send a verification email to the user's registered email address
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param request body dto.SendVerificationEmailRequest true "Send verification email request"
// @Success 200 {object} dto.EmailVerificationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/send-verification-email [post]
func (h *EmailVerificationHandler) SendVerificationEmail(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.SendVerificationEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.verificationService.SendVerificationEmail(ctx, req.UserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("email.send_failed"))
	}

	return c.JSON(http.StatusOK, dto.EmailVerificationResponse{
		Message: t.Success("email.verification_sent"),
		Success: true,
	})
}

// VerifyEmail verifies the user's email address using the verification token
// @Summary Verify email address
// @Description Verify the user's email address using the verification token
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Verify email request"
// @Success 200 {object} dto.EmailVerificationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/verify-email [post]
func (h *EmailVerificationHandler) VerifyEmail(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.verificationService.VerifyEmail(ctx, req.Token); err != nil {
		if err.Error() == "invalid or expired verification token" {
			return echo.NewHTTPError(http.StatusNotFound, t.Error("email.invalid_token"))
		}
		return echo.NewHTTPError(http.StatusBadRequest, t.Error("email.verification_failed"))
	}

	return c.JSON(http.StatusOK, dto.EmailVerificationResponse{
		Message: t.Success("email.verification_success"),
		Success: true,
	})
}

// VerifyEmailByToken verifies the user's email address using the verification token from URL parameter
// @Summary Verify email address by token
// @Description Verify the user's email address using the verification token from URL parameter
// @Tags Email Verification
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} dto.EmailVerificationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/verify-email/{token} [get]
func (h *EmailVerificationHandler) VerifyEmailByToken(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	token := c.QueryParam("token")
	if token == "" {
		token = c.Param("token")
	}

	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, t.Error("email.token_required"))
	}

	if err := h.verificationService.VerifyEmail(ctx, token); err != nil {
		if err.Error() == "invalid or expired verification token" {
			return echo.NewHTTPError(http.StatusNotFound, t.Error("email.invalid_token"))
		}
		return echo.NewHTTPError(http.StatusBadRequest, t.Error("email.verification_failed"))
	}

	return c.JSON(http.StatusOK, dto.EmailVerificationResponse{
		Message: t.Success("email.verification_success"),
		Success: true,
	})
}

// ResendVerificationEmail resends the verification email to the user
// @Summary Resend verification email
// @Description Resend the verification email to the user's registered email address
// @Tags Email Verification
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationEmailRequest true "Resend verification email request"
// @Success 200 {object} dto.EmailVerificationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 429 {object} echo.HTTPError
// @Router /auth/resend-verification-email [post]
func (h *EmailVerificationHandler) ResendVerificationEmail(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	ctx := contextx.NewWithRequestContext(c)

	var req dto.ResendVerificationEmailRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	if err := h.verificationService.ResendVerificationEmail(ctx, req.UserID); err != nil {
		if err.Error() == "verification email was sent recently, please wait before requesting again" {
			return echo.NewHTTPError(http.StatusTooManyRequests, t.Error("email.rate_limited"))
		}
		if err.Error() == "email already verified" {
			return echo.NewHTTPError(http.StatusBadRequest, t.Error("email.already_verified"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("email.send_failed"))
	}

	return c.JSON(http.StatusOK, dto.EmailVerificationResponse{
		Message: t.Success("email.verification_sent"),
		Success: true,
	})
}

// GetVerificationStatus gets the current email verification status for a user
// @Summary Get verification status
// @Description Get the current email verification status for a user
// @Tags Email Verification
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.EmailVerificationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Router /auth/verification-status/{user_id} [get]
func (h *EmailVerificationHandler) GetVerificationStatus(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())

	userIDStr := c.Param("user_id")
	_, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.Error("validation.invalid_user_id"))
	}

	// This would need to be implemented in the verification service
	// For now, we'll just return a placeholder response
	return c.JSON(http.StatusOK, dto.EmailVerificationResponse{
		Message: "Feature not yet implemented",
		Success: false,
	})
}