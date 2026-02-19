package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

// AuthRequired is a middleware that requires a valid JWT session.
func AuthRequired(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// 2. Extract bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}
		tokenString := parts[1]

		// 3. Validate session using AuthService
		sessionInfo, err := authService.ValidateSession(c.Context(), tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired session",
			})
		}

		// 4. Store person ID and email in locals for downstream handlers
		c.Locals("person_id", sessionInfo.PersonID)
		c.Locals("email", sessionInfo.Email)

		return c.Next()
	}
}

// OptionalAuth is a middleware that extracts session if present but doesn't require it.
func OptionalAuth(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]
			sessionInfo, err := authService.ValidateSession(c.Context(), tokenString)
			if err == nil {
				c.Locals("person_id", sessionInfo.PersonID)
				c.Locals("email", sessionInfo.Email)
			}
		}

		return c.Next()
	}
}
