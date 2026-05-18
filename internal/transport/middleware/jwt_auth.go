package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-beresin/pkg/utils"
)

// Response represents a standard JSON response structure matching our API.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JWTAuth parses, validates, and extracts claims from JWT Bearer token.
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(Response{
				Status:  "error",
				Message: "Missing authorization token",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(Response{
				Status:  "error",
				Message: "Authorization header must be in the format 'Bearer <token>'",
			})
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(Response{
				Status:  "error",
				Message: "Invalid or expired authorization token",
			})
		}

		// Store user details in fiber context locals
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RequireRole restricts access to specific roles.
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")
		if roleVal == nil {
			return c.Status(fiber.StatusForbidden).JSON(Response{
				Status:  "error",
				Message: "Access forbidden: missing user role",
			})
		}

		userRole, ok := roleVal.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(Response{
				Status:  "error",
				Message: "Access forbidden: invalid role format",
			})
		}

		// Check if the user's role matches any of the allowed roles
		isAllowed := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userRole, role) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(Response{
				Status:  "error",
				Message: "Access forbidden: you do not have permission to access this resource",
			})
		}

		return c.Next()
	}
}
