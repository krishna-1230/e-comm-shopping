package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupCartRoutes sets up all cart routes
func SetupCartRoutes(app *fiber.App) {
	// All cart routes require authentication
	cartRoutes := app.Group("/api/cart", middlewares.Protected())

	// Cart endpoints
	cartRoutes.Get("/", controllers.GetCart)
	cartRoutes.Post("/", controllers.AddToCart)
	cartRoutes.Put("/:id", controllers.UpdateCartItem)
	cartRoutes.Delete("/:id", controllers.RemoveFromCart)
	cartRoutes.Delete("/", controllers.ClearCart)
}

// SetupWishlistRoutes sets up all wishlist routes
func SetupWishlistRoutes(app *fiber.App) {
	// All wishlist routes require authentication
	wishlistRoutes := app.Group("/api/wishlist", middlewares.Protected())

	// Wishlist endpoints
	wishlistRoutes.Get("/", controllers.GetWishlist)
	wishlistRoutes.Post("/", controllers.AddToWishlist)
	wishlistRoutes.Delete("/:id", controllers.RemoveFromWishlist)
	wishlistRoutes.Delete("/", controllers.ClearWishlist)
}

// SetupAddressRoutes sets up all address routes
func SetupAddressRoutes(app *fiber.App) {
	// All address routes require authentication
	addressRoutes := app.Group("/api/addresses", middlewares.Protected())

	// Address endpoints
	addressRoutes.Get("/", controllers.GetAllAddresses)
	addressRoutes.Get("/:id", controllers.GetAddress)
	addressRoutes.Post("/", controllers.CreateAddress)
	addressRoutes.Put("/:id", controllers.UpdateAddress)
	addressRoutes.Delete("/:id", controllers.DeleteAddress)
}

// SetupOrderRoutes sets up all order routes
func SetupOrderRoutes(app *fiber.App) {
	// All order routes require authentication
	orderRoutes := app.Group("/api/orders", middlewares.Protected())

	// Order endpoints for all users
	orderRoutes.Post("/", controllers.PlaceOrder)
	orderRoutes.Get("/", controllers.GetAllOrders)
	orderRoutes.Get("/:id", controllers.GetOrderByID)

	// Admin only endpoints
	adminRoutes := orderRoutes.Use(middlewares.AdminOnly())
	adminRoutes.Put("/:id/status", controllers.UpdateOrderStatus)
}
