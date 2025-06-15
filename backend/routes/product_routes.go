package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

// SetupProductRoutes sets up all product routes
func SetupProductRoutes(app *fiber.App) {
	// Product endpoints
	productRoutes := app.Group("/api/products")
	
	// Public routes
	productRoutes.Get("/", controllers.GetAllProducts)
	productRoutes.Get("/:id", controllers.GetProductByID)
	
	// Protected routes (admin only)
	admin := productRoutes.Use(middlewares.AdminOnly())
	admin.Post("/", controllers.CreateProduct)
	admin.Put("/:id", controllers.UpdateProduct)
	admin.Delete("/:id", controllers.DeleteProduct)
	
	// Product inventory management (admin only)
	admin.Post("/:id/colors", controllers.AddProductColor)
	admin.Post("/:id/sizes", controllers.AddProductSize)
	admin.Post("/:id/inventory", controllers.UpdateInventory)
	admin.Post("/:id/images", controllers.AddProductImage)
	
	// Delete product attributes (admin only)
	admin.Delete("/:id/colors/:colorId", controllers.DeleteProductColor)
	admin.Delete("/:id/sizes/:sizeId", controllers.DeleteProductSize)
	admin.Delete("/:id/images/:imageId", controllers.DeleteProductImage)
	
	// Category endpoints
	categoryRoutes := app.Group("/api/categories")
	
	// Public routes
	categoryRoutes.Get("/", controllers.GetAllCategories)
	categoryRoutes.Get("/:id", controllers.GetCategoryByID)
	
	// Protected routes (admin only)
	adminCategory := categoryRoutes.Use(middlewares.AdminOnly())
	adminCategory.Post("/", controllers.CreateCategory)
	adminCategory.Put("/:id", controllers.UpdateCategory)
	adminCategory.Delete("/:id", controllers.DeleteCategory)
} 