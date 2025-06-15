package controllers

import (
	"backend/database"
	"backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AddToCart adds a product to the user's cart
func AddToCart(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Parse request body
	var req models.CartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.ProductID <= 0 || req.ColorID <= 0 || req.SizeID <= 0 || req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID, color ID, size ID, and quantity are required and must be positive",
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

	// Check if color exists for this product
	var colorExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_colors WHERE id = ? AND product_id = ?)", req.ColorID, req.ProductID).Scan(&colorExists)
	if err != nil || !colorExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Color not found for this product",
		})
	}

	// Check if size exists for this product
	var sizeExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_sizes WHERE id = ? AND product_id = ?)", req.SizeID, req.ProductID).Scan(&sizeExists)
	if err != nil || !sizeExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Size not found for this product",
		})
	}

	// Check if there's enough inventory
	var availableQuantity int
	err = database.DB.QueryRow(
		"SELECT quantity FROM product_inventory WHERE product_id = ? AND color_id = ? AND size_id = ?",
		req.ProductID, req.ColorID, req.SizeID).Scan(&availableQuantity)

	if err != nil || availableQuantity < req.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Not enough inventory",
			"available": availableQuantity,
		})
	}

	// Check if item already exists in cart
	var cartItemExists bool
	var existingCartID int64
	var existingQuantity int
	err = database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM cart 
			WHERE user_id = ? AND product_id = ? AND color_id = ? AND size_id = ?
		), id, quantity 
		FROM cart 
		WHERE user_id = ? AND product_id = ? AND color_id = ? AND size_id = ?`,
		userID, req.ProductID, req.ColorID, req.SizeID,
		userID, req.ProductID, req.ColorID, req.SizeID).Scan(
		&cartItemExists, &existingCartID, &existingQuantity)

	if err != nil && cartItemExists {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// If item exists, update quantity
	if cartItemExists {
		_, err = database.DB.Exec(
			"UPDATE cart SET quantity = ?, updated_at = ? WHERE id = ?",
			req.Quantity, time.Now(), existingCartID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cart",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart updated successfully",
			"id":      existingCartID,
		})
	}

	// Otherwise, add new item to cart
	result, err := database.DB.Exec(
		"INSERT INTO cart (user_id, product_id, color_id, size_id, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, req.ProductID, req.ColorID, req.SizeID, req.Quantity, time.Now(), time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add to cart",
		})
	}

	// Get the cart item ID
	cartItemID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Added to cart successfully",
		"id":      cartItemID,
	})
}

// GetCart retrieves the user's cart
func GetCart(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Query to get cart items with product details
	rows, err := database.DB.Query(`
		SELECT 
			c.id, c.product_id, c.color_id, c.size_id, c.quantity,
			p.name, p.description, p.base_price, p.discount_percentage,
			pc.color_name, pc.color_hex,
			ps.size_name,
			pi.quantity as in_stock,
			(SELECT image_url FROM product_images WHERE product_id = p.id AND is_primary = 1 LIMIT 1) as image_url
		FROM cart c
		JOIN products p ON c.product_id = p.id
		JOIN product_colors pc ON c.color_id = pc.id
		JOIN product_sizes ps ON c.size_id = ps.id
		LEFT JOIN product_inventory pi ON c.product_id = pi.product_id AND c.color_id = pi.color_id AND c.size_id = pi.size_id
		WHERE c.user_id = ?
		ORDER BY c.id DESC`,
		userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	var cartItems []models.CartItemResponse
	var totalItems int
	var subTotal float64

	for rows.Next() {
		var item models.CartItemResponse
		var basePrice, discountPercentage float64
		var inStock int

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.ColorID, &item.SizeID, &item.Quantity,
			&item.ProductName, &item.ProductDescription, &basePrice, &discountPercentage,
			&item.ColorName, &item.ColorHex,
			&item.SizeName,
			&inStock,
			&item.ImageURL)
		if err != nil {
			continue
		}

		// Calculate final price and subtotal
		item.BasePrice = basePrice
		item.DiscountPercentage = discountPercentage
		item.FinalPrice = basePrice * (1 - discountPercentage/100)
		item.SubTotal = item.FinalPrice * float64(item.Quantity)
		item.InStock = inStock

		cartItems = append(cartItems, item)
		totalItems += item.Quantity
		subTotal += item.SubTotal
	}

	// Calculate cart summary
	shippingCost := 0.0
	taxRate := 0.1 // 10% tax
	tax := subTotal * taxRate

	summary := models.CartSummary{
		TotalItems:   totalItems,
		SubTotal:     subTotal,
		ShippingCost: shippingCost,
		Tax:          tax,
		Total:        subTotal + shippingCost + tax,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"items":   cartItems,
		"summary": summary,
	})
}

// UpdateCartItem updates the quantity of a cart item
func UpdateCartItem(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get cart item ID from URL parameter
	id := c.Params("id")
	cartItemID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart item ID",
		})
	}

	// Parse request body
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity must be positive",
		})
	}

	// Check if cart item exists and belongs to user
	var exists bool
	var productID, colorID, sizeID int64
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM cart WHERE id = ? AND user_id = ?), product_id, color_id, size_id FROM cart WHERE id = ?",
		cartItemID, userID, cartItemID).Scan(&exists, &productID, &colorID, &sizeID)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	// Check if there's enough inventory
	var availableQuantity int
	err = database.DB.QueryRow(
		"SELECT quantity FROM product_inventory WHERE product_id = ? AND color_id = ? AND size_id = ?",
		productID, colorID, sizeID).Scan(&availableQuantity)

	if err != nil || availableQuantity < req.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Not enough inventory",
			"available": availableQuantity,
		})
	}

	// Update cart item quantity
	_, err = database.DB.Exec(
		"UPDATE cart SET quantity = ?, updated_at = ? WHERE id = ? AND user_id = ?",
		req.Quantity, time.Now(), cartItemID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart updated successfully",
	})
}

// RemoveFromCart removes an item from the user's cart
func RemoveFromCart(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get cart item ID from URL parameter
	id := c.Params("id")
	cartItemID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart item ID",
		})
	}

	// Check if cart item exists and belongs to user
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM cart WHERE id = ? AND user_id = ?)", cartItemID, userID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	// Delete cart item
	_, err = database.DB.Exec("DELETE FROM cart WHERE id = ? AND user_id = ?", cartItemID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove from cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Removed from cart successfully",
	})
}

// ClearCart removes all items from the user's cart
func ClearCart(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Delete all cart items for this user
	_, err := database.DB.Exec("DELETE FROM cart WHERE user_id = ?", userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart cleared successfully",
	})
}
