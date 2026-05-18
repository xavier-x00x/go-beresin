package handler

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/gofiber/fiber/v2"

	"go-beresin/internal/domain"
	"go-beresin/internal/service"
)

var phoneRegex = regexp.MustCompile(`^(?:\+62|62|0)8[1-9][0-9]{6,10}$`)

func validateRegisterInput(email, phone, password string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.NewValidationError("invalid email format")
	}
	if !phoneRegex.MatchString(phone) {
		return domain.NewValidationError("invalid Indonesian phone number format")
	}
	if len(password) < 8 {
		return domain.NewValidationError("password must be at least 8 characters long")
	}
	return nil
}

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

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(domain.RegisterReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if err := validateRegisterInput(req.Email, req.Phone, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	ctx := context.Background()
	resp, err := h.service.Register(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:  "success",
		Message: "Registration completed successfully. Please check your email to verify your account.",
		Data: fiber.Map{
			"user_id":   resp.UserID,
			"email":     resp.Email,
			"phone":     resp.Phone,
			"full_name": resp.FullName,
			"role":      resp.Role,
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(domain.LoginReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	resp, err := h.service.Login(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated",
		Data: fiber.Map{
			"user_id":       resp.UserID,
			"full_name":     resp.FullName,
			"role":          resp.Role,
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		},
	})
}

func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	req := new(domain.GoogleLoginReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	resp, err := h.service.GoogleLogin(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated via Google",
		Data: fiber.Map{
			"user_id":       resp.UserID,
			"full_name":     resp.FullName,
			"role":          resp.Role,
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		},
	})
}

func (h *AuthHandler) SendWhatsAppOTP(c *fiber.Ctx) error {
	req := new(domain.SendOtpReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	if !phoneRegex.MatchString(req.Phone) {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "invalid Indonesian phone number format",
		})
	}

	ctx := context.Background()
	err := h.service.SendWhatsAppOTP(ctx, *req)
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "OTP has been successfully dispatched to WhatsApp",
	})
}

func (h *AuthHandler) VerifyWhatsAppOTP(c *fiber.Ctx) error {
	req := new(domain.VerifyOtpReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	resp, err := h.service.VerifyWhatsAppOTP(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully authenticated via WhatsApp OTP",
		Data: fiber.Map{
			"user_id":       resp.UserID,
			"full_name":     resp.FullName,
			"role":          resp.Role,
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	req := new(domain.RefreshTokenReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	resp, err := h.service.RefreshToken(ctx, *req)
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Tokens successfully rotated",
		Data: fiber.Map{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Missing authenticated user details",
		})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(Response{
			Status:  "error",
			Message: "Missing authenticated user details",
		})
	}

	ctx := context.Background()
	_ = h.service.Logout(ctx, userIDStr)

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Successfully logged out, session terminated",
	})
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	req := new(domain.ForgotPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	err := h.service.ForgotPassword(ctx, *req)
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "If the email exists, password reset instructions have been dispatched.",
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req := new(domain.ResetPasswordReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	ctx := context.Background()
	err := h.service.ResetPassword(ctx, *req, c.IP(), c.Get("User-Agent"))
	if err != nil {
		return c.Status(mapServiceError(err)).JSON(Response{
			Status:  "error",
			Message: errorMessage(err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Status:  "success",
		Message: "Password updated successfully. All other sessions have been logged out.",
	})
}
