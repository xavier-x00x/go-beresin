package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"go-beresin/internal/domain"
	"go-beresin/internal/service"
	"go-beresin/internal/transport/response"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func mapServiceError(err error) int {
	if err == nil {
		return 0
	}

	var valErr *domain.ValidationError
	if errors.As(err, &valErr) {
		return fiber.StatusBadRequest
	}

	switch {
	case errors.Is(err, domain.ErrEmailConflict):
		return fiber.StatusConflict
	case errors.Is(err, domain.ErrPhoneConflict):
		return fiber.StatusConflict
	case errors.Is(err, domain.ErrInvalidCredentials):
		return fiber.StatusUnauthorized
	case errors.Is(err, domain.ErrIPBlocked):
		return fiber.StatusTooManyRequests
	default:
		return fiber.StatusInternalServerError
	}
}

func errorMessage(err error) string {
	var valErr *domain.ValidationError
	if errors.As(err, &valErr) {
		return valErr.Message
	}

	if errors.Is(err, domain.ErrIPBlocked) {
		var ttl int
		_, scanErr := fmt.Sscanf(err.Error(), "%*s: %d", &ttl)
		if scanErr == nil {
			return fmt.Sprintf("IP temporarily blocked due to multiple failed login attempts. Please retry after %d seconds.", ttl)
		}
		return "IP temporarily blocked due to multiple failed login attempts."
	}

	switch {
	case errors.Is(err, domain.ErrEmailConflict):
		return "Email is already registered"
	case errors.Is(err, domain.ErrPhoneConflict):
		return "Phone number is already registered"
	case errors.Is(err, domain.ErrInvalidCredentials):
		return "Invalid email or password credentials"
	default:
		return "An unexpected error occurred"
	}
}

// Register handles new user registration.
// @Summary User Registration
// @Description Register a new client or talent account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterReq true "Register Payload"
// @Success 201 {object} RegisterResp
// @Failure 400 {object} ErrorResp
// @Failure 409 {object} ErrorResp
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(domain.RegisterReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := h.service.Register(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusCreated, "Registration completed successfully. Please check your email to verify your account.", fiber.Map{
		"user_id":   resp.UserID,
		"email":     resp.Email,
		"phone":     resp.Phone,
		"full_name": resp.FullName,
		"role":      resp.Role,
	})
}

// Login handles email and password authentication.
// @Summary User Login
// @Description Login with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginReq true "Login Payload"
// @Success 200 {object} LoginResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Failure 429 {object} ErrorResp
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(domain.LoginReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := h.service.Login(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "Successfully authenticated", fiber.Map{
		"user_id":       resp.UserID,
		"full_name":     resp.FullName,
		"role":          resp.Role,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// GoogleLogin handles Google OAuth authentication.
// @Summary Google Login
// @Description Authenticate using Google OAuth ID token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.GoogleLoginReq true "Google Login Payload"
// @Success 200 {object} LoginResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /api/v1/auth/login/google [post]
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	req := new(domain.GoogleLoginReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := h.service.GoogleLogin(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "Successfully authenticated via Google", fiber.Map{
		"user_id":       resp.UserID,
		"full_name":     resp.FullName,
		"role":          resp.Role,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// SendWhatsAppOTP dispatches OTP via WhatsApp.
// @Summary Send WhatsApp OTP
// @Description Send OTP code to phone number via WhatsApp
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.SendOtpReq true "Send OTP Payload"
// @Success 200 {object} MessageResp
// @Failure 400 {object} ErrorResp
// @Router /api/v1/auth/login/whatsapp/send-otp [post]
func (h *AuthHandler) SendWhatsAppOTP(c *fiber.Ctx) error {
	req := new(domain.SendOtpReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	err := h.service.SendWhatsAppOTP(ctx, *req)
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "OTP has been successfully dispatched to WhatsApp", nil)
}

// VerifyWhatsAppOTP validates OTP and authenticates user.
// @Summary Verify WhatsApp OTP
// @Description Verify OTP code and authenticate via WhatsApp
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.VerifyOtpReq true "Verify OTP Payload"
// @Success 200 {object} LoginResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /api/v1/auth/login/whatsapp/verify [post]
func (h *AuthHandler) VerifyWhatsAppOTP(c *fiber.Ctx) error {
	req := new(domain.VerifyOtpReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := h.service.VerifyWhatsAppOTP(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "Successfully authenticated via WhatsApp OTP", fiber.Map{
		"user_id":       resp.UserID,
		"full_name":     resp.FullName,
		"role":          resp.Role,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// RefreshToken handles token rotation.
// @Summary Refresh Token
// @Description Rotate access and refresh token pair
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenReq true "Refresh Token Payload"
// @Success 200 {object} RefreshTokenResp
// @Failure 400 {object} ErrorResp
// @Failure 401 {object} ErrorResp
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(domain.RefreshTokenReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	resp, err := h.service.RefreshToken(ctx, *req)
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "Tokens successfully rotated", fiber.Map{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

// Logout terminates user session.
// @Summary Logout
// @Description Invalidate refresh token and terminate session
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResp
// @Failure 401 {object} ErrorResp
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "Missing authenticated user details")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Missing authenticated user details")
	}

	ctx := context.Background()
	_ = h.service.Logout(ctx, userIDStr)

	return response.Success(c, fiber.StatusOK, "Successfully logged out, session terminated", nil)
}

// ForgotPassword initiates password reset flow.
// @Summary Forgot Password
// @Description Send password reset link to email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.ForgotPasswordReq true "Forgot Password Payload"
// @Success 200 {object} MessageResp
// @Failure 400 {object} ErrorResp
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	req := new(domain.ForgotPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	err := h.service.ForgotPassword(ctx, *req)
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "If the email exists, password reset instructions have been dispatched.", nil)
}

// Profile handles retrieving authenticated user details from local context.
// @Summary Get Profile
// @Description Fetch authenticated user identity from Bearer token claims
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MessageResp
// @Failure 401 {object} ErrorResp
// @Router /api/v1/profile [get]
func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	role := c.Locals("role")

	return response.Success(c, fiber.StatusOK, "Profile retrieved successfully", fiber.Map{
		"user_id": userID,
		"role":    role,
	})
}

// ResetPassword completes password reset.
// @Summary Reset Password
// @Description Reset password using email reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.ResetPasswordReq true "Reset Password Payload"
// @Success 200 {object} MessageResp
// @Failure 400 {object} ErrorResp
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req := new(domain.ResetPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	if err := response.Check(c, req); err != nil {
		return err
	}

	ctx := context.Background()
	err := h.service.ResetPassword(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return response.Error(c, mapServiceError(err), errorMessage(err))
	}

	return response.Success(c, fiber.StatusOK, "Password updated successfully. All other sessions have been logged out.", nil)
}
