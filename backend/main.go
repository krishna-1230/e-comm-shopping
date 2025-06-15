package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment")
	}

	// Initialize database
	database.InitDatabase()
	defer database.CloseDatabase()

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		AppName: "E-Commerce API",
	})

	// Setup middlewares
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Routes setup
	routes.SetupUserRoutes(app)
	routes.SetupProductRoutes(app)
	routes.SetupCartRoutes(app)
	routes.SetupWishlistRoutes(app)
	routes.SetupAddressRoutes(app)
	routes.SetupOrderRoutes(app)

	// Health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// Handle not found routes
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Fatal(app.Listen(":" + port))
}
