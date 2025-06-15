package controllers

import (
	"backend/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AddProductColor adds a new color to a product
func AddProductColor(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Check if product exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Parse request body
	var color struct {
		ColorName string `json:"color_name"`
		ColorHex  string `json:"color_hex"`
	}
	if err := c.BodyParser(&color); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if color.ColorName == "" || color.ColorHex == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Color name and hex code are required",
		})
	}

	// Insert the color
	result, err := database.DB.Exec(
		"INSERT INTO product_colors (product_id, color_name, color_hex, created_at) VALUES (?, ?, ?, ?)",
		productID, color.ColorName, color.ColorHex, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add color",
		})
	}

	// Get the color ID
	colorID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Color added successfully",
		"id":      colorID,
	})
}

// AddProductSize adds a new size to a product
func AddProductSize(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Check if product exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Parse request body
	var size struct {
		SizeName string `json:"size_name"`
	}
	if err := c.BodyParser(&size); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if size.SizeName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Size name is required",
		})
	}

	// Insert the size
	result, err := database.DB.Exec(
		"INSERT INTO product_sizes (product_id, size_name, created_at) VALUES (?, ?, ?)",
		productID, size.SizeName, time.Now())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add size",
		})
	}

	// Get the size ID
	sizeID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Size added successfully",
		"id":      sizeID,
	})
}

// UpdateInventory updates the inventory for a product variant
func UpdateInventory(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Check if product exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Parse request body
	var inventory struct {
		ColorID  int64 `json:"color_id"`
		SizeID   int64 `json:"size_id"`
		Quantity int   `json:"quantity"`
	}
	if err := c.BodyParser(&inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if inventory.ColorID <= 0 || inventory.SizeID <= 0 || inventory.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Valid color ID, size ID, and quantity are required",
		})
	}

	// Check if color and size exist for this product
	var colorExists, sizeExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_colors WHERE id = ? AND product_id = ?)", inventory.ColorID, productID).Scan(&colorExists)
	if err != nil || !colorExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid color for this product",
		})
	}

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_sizes WHERE id = ? AND product_id = ?)", inventory.SizeID, productID).Scan(&sizeExists)
	if err != nil || !sizeExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid size for this product",
		})
	}

	// Check if inventory entry exists
	var inventoryExists bool
	var inventoryID int64
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM product_inventory WHERE product_id = ? AND color_id = ? AND size_id = ?)",
		productID, inventory.ColorID, inventory.SizeID).Scan(&inventoryExists)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update or insert inventory
	if inventoryExists {
		// Update existing inventory
		_, err = database.DB.Exec(
			"UPDATE product_inventory SET quantity = ?, updated_at = ? WHERE product_id = ? AND color_id = ? AND size_id = ?",
			inventory.Quantity, time.Now(), productID, inventory.ColorID, inventory.SizeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update inventory",
			})
		}

		// Get the inventory ID
		database.DB.QueryRow(
			"SELECT id FROM product_inventory WHERE product_id = ? AND color_id = ? AND size_id = ?",
			productID, inventory.ColorID, inventory.SizeID).Scan(&inventoryID)

	} else {
		// Insert new inventory
		result, err := database.DB.Exec(
			"INSERT INTO product_inventory (product_id, color_id, size_id, quantity, updated_at) VALUES (?, ?, ?, ?, ?)",
			productID, inventory.ColorID, inventory.SizeID, inventory.Quantity, time.Now())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add inventory",
			})
		}

		// Get the inventory ID
		inventoryID, _ = result.LastInsertId()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Inventory updated successfully",
		"id":      inventoryID,
	})
}

// AddProductImage adds a new image to a product
func AddProductImage(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Check if product exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ?)", productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Parse request body
	var image struct {
		ImageURL  string `json:"image_url"`
		IsPrimary bool   `json:"is_primary"`
	}
	if err := c.BodyParser(&image); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if image.ImageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image URL is required",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// If this is the primary image, update all other images to not be primary
	if image.IsPrimary {
		_, err = tx.Exec("UPDATE product_images SET is_primary = 0 WHERE product_id = ?", productID)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update existing images",
			})
		}
	}

	// Insert the image
	result, err := tx.Exec(
		"INSERT INTO product_images (product_id, image_url, is_primary, created_at) VALUES (?, ?, ?, ?)",
		productID, image.ImageURL, image.IsPrimary, time.Now())
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add image",
		})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Get the image ID
	imageID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Image added successfully",
		"id":      imageID,
	})
}

// DeleteProductColor deletes a color from a product
func DeleteProductColor(c *fiber.Ctx) error {
	// Get the product ID and color ID from the URL parameters
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	colorID, err := strconv.ParseInt(c.Params("colorId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid color ID",
		})
	}

	// Check if color exists for this product
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_colors WHERE id = ? AND product_id = ?)", colorID, productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Color not found for this product",
		})
	}

	// Delete the color
	_, err = database.DB.Exec("DELETE FROM product_colors WHERE id = ? AND product_id = ?", colorID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete color",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Color deleted successfully",
	})
}

// DeleteProductSize deletes a size from a product
func DeleteProductSize(c *fiber.Ctx) error {
	// Get the product ID and size ID from the URL parameters
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	sizeID, err := strconv.ParseInt(c.Params("sizeId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid size ID",
		})
	}

	// Check if size exists for this product
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_sizes WHERE id = ? AND product_id = ?)", sizeID, productID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Size not found for this product",
		})
	}

	// Delete the size
	_, err = database.DB.Exec("DELETE FROM product_sizes WHERE id = ? AND product_id = ?", sizeID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete size",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Size deleted successfully",
	})
}

// DeleteProductImage deletes an image from a product
func DeleteProductImage(c *fiber.Ctx) error {
	// Get the product ID and image ID from the URL parameters
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	imageID, err := strconv.ParseInt(c.Params("imageId"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid image ID",
		})
	}

	// Check if image exists for this product
	var exists bool
	var isPrimary bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM product_images WHERE id = ? AND product_id = ?), is_primary FROM product_images WHERE id = ?", imageID, productID, imageID).Scan(&exists, &isPrimary)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found for this product",
		})
	}

	// Delete the image
	_, err = database.DB.Exec("DELETE FROM product_images WHERE id = ? AND product_id = ?", imageID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete image",
		})
	}

	// If this was the primary image, set another image as primary
	if isPrimary {
		database.DB.Exec("UPDATE product_images SET is_primary = 1 WHERE product_id = ? ORDER BY id LIMIT 1", productID)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Image deleted successfully",
	})
}
