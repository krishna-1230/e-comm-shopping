package controllers

import (
	"backend/database"
	"backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AddToWishlist adds a product to the user's wishlist
func AddToWishlist(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Parse request body
	var req struct {
		ProductID int64 `json:"product_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	// Check if product exists
	var productExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", req.ProductID).Scan(&productExists)
	if err != nil || !productExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Check if item already exists in wishlist
	var wishlistExists bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM wishlist WHERE user_id = ? AND product_id = ?)",
		userID, req.ProductID).Scan(&wishlistExists)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if wishlistExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Product already in wishlist",
		})
	}

	// Add to wishlist
	result, err := database.DB.Exec(
		"INSERT INTO wishlist (user_id, product_id, created_at) VALUES (?, ?, ?)",
		userID, req.ProductID, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to wishlist",
		})
	}

	// Get the wishlist item ID
	wishlistItemID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Added to wishlist successfully",
		"id":      wishlistItemID,
	})
}

// GetWishlist retrieves the user's wishlist
func GetWishlist(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Query to get wishlist items with product details
	rows, err := database.DB.Query(`
		SELECT 
			w.id, w.product_id, w.created_at,
			p.name, p.description, p.base_price, p.discount_percentage,
			(SELECT image_url FROM product_images WHERE product_id = p.id AND is_primary = 1 LIMIT 1) as image_url,
			(SELECT COUNT(*) > 0 FROM product_inventory pi 
				JOIN product_colors pc ON pi.color_id = pc.id 
				JOIN product_sizes ps ON pi.size_id = ps.id 
				WHERE pi.product_id = p.id AND pi.quantity > 0) as in_stock
		FROM wishlist w
		JOIN products p ON w.product_id = p.id
		WHERE w.user_id = ?
		ORDER BY w.id DESC`,
		userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	var wishlistItems []models.WishlistItemResponse

	for rows.Next() {
		var item models.WishlistItemResponse
		var basePrice, discountPercentage float64
		var inStock bool

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.CreatedAt,
			&item.ProductName, &item.ProductDescription, &basePrice, &discountPercentage,
			&item.ImageURL, &inStock)
		if err != nil {
			continue
		}

		// Calculate final price
		item.BasePrice = basePrice
		item.DiscountPercentage = discountPercentage
		item.FinalPrice = basePrice * (1 - discountPercentage/100)
		item.InStock = inStock

		wishlistItems = append(wishlistItems, item)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items": wishlistItems,
	})
}

// RemoveFromWishlist removes a product from the user's wishlist
func RemoveFromWishlist(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get wishlist item ID from URL parameter
	id := c.Params("id")
	wishlistItemID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid wishlist item ID",
		})
	}

	// Check if wishlist item exists and belongs to user
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM wishlist WHERE id = ? AND user_id = ?)", wishlistItemID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Wishlist item not found",
		})
	}

	// Delete wishlist item
	_, err = database.DB.Exec("DELETE FROM wishlist WHERE id = ? AND user_id = ?", wishlistItemID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove from wishlist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Removed from wishlist successfully",
	})
}

// ClearWishlist removes all items from the user's wishlist
func ClearWishlist(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Delete all wishlist items for this user
	_, err := database.DB.Exec("DELETE FROM wishlist WHERE user_id = ?", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear wishlist",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Wishlist cleared successfully",
	})
}
