package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up all the user routes
func SetupUserRoutes(app *fiber.App) {
	// Public routes
	app.Post("/api/auth/register", controllers.RegisterUser)
	app.Post("/api/auth/login", controllers.LoginUser)

	// Protected routes
	app.Get("/api/auth/me", middlewares.Protected(), controllers.GetCurrentUser)
} 