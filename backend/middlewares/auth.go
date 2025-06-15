package middlewares

import (
	"backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Protected is a middleware that checks if the user is authenticated
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No token provided",
			})
		}

		// Check if the header format is correct
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token format",
			})
		}

		// Get the token
		tokenString := parts[1]

		// Validate the token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid token",
			})
		}

		// Set user data in context
		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		// Continue
		return c.Next()
	}
}

// AdminOnly is a middleware that checks if the user is an admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First run the Protected middleware
		err := Protected()(c)
		if err != nil {
			return err
		}

		// Check if the user is an admin
		role := c.Locals("role")
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Admin access required",
			})
		}

		// Continue
		return c.Next()
	}
} 