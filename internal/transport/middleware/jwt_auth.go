package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go-beresin/internal/transport/response"
	"go-beresin/pkg/utils"
)

// JWTAuth parses, validates, and extracts claims from JWT Bearer token.
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization token")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header must be in the format 'Bearer <token>'")
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired authorization token")
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
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: missing user role")
		}

		userRole, ok := roleVal.(string)
		if !ok {
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: invalid role format")
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
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: you do not have permission to access this resource")
		}

		return c.Next()
	}
}
