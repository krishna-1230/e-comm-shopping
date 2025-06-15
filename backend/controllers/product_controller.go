package controllers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateProduct handles product creation
func CreateProduct(c *fiber.Ctx) error {
	// Parse request body
	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.BasePrice <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and price are required",
		})
	}

	// Create the product
	result, err := database.DB.Exec(
		"INSERT INTO products (name, description, category_id, base_price, discount_percentage, featured, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		req.Name,
		req.Description,
		req.CategoryID,
		req.BasePrice,
		req.DiscountPercentage,
		req.Featured,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	// Get the product ID
	productID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"id":      productID,
	})
}

// GetAllProducts returns all products
func GetAllProducts(c *fiber.Ctx) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query to get products
	rows, err := database.DB.Query(`
		SELECT p.id, p.name, p.description, p.category_id, p.base_price, 
			   p.discount_percentage, p.featured, p.created_at, p.updated_at,
			   IFNULL(c.name, 'Uncategorized') as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id DESC
		LIMIT ? OFFSET ?`,
		limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	products := []map[string]interface{}{}
	for rows.Next() {
		var product models.Product
		var categoryName string
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.CategoryID,
			&product.BasePrice, &product.DiscountPercentage, &product.Featured,
			&product.CreatedAt, &product.UpdatedAt, &categoryName)
		if err != nil {
			continue
		}

		// Calculate final price
		finalPrice := product.BasePrice * (1 - product.DiscountPercentage/100)

		productMap := map[string]interface{}{
			"id":                  product.ID,
			"name":                product.Name,
			"description":         product.Description,
			"category_id":         product.CategoryID,
			"category_name":       categoryName,
			"base_price":          product.BasePrice,
			"discount_percentage": product.DiscountPercentage,
			"final_price":         finalPrice,
			"featured":            product.Featured,
			"created_at":          product.CreatedAt,
			"updated_at":          product.UpdatedAt,
		}

		// Get primary image
		var imageURL sql.NullString
		database.DB.QueryRow(`
			SELECT image_url FROM product_images 
			WHERE product_id = ? AND is_primary = 1
			LIMIT 1`, product.ID).Scan(&imageURL)

		if imageURL.Valid {
			productMap["primary_image"] = imageURL.String
		} else {
			productMap["primary_image"] = nil
		}

		products = append(products, productMap)
	}

	// Count total products for pagination
	var total int
	database.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&total)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
		"meta": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetProductByID returns a specific product by ID
func GetProductByID(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	id := c.Params("id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Get the product from the database
	var product models.Product
	var categoryName string
	err = database.DB.QueryRow(`
		SELECT p.id, p.name, p.description, p.category_id, p.base_price, 
			   p.discount_percentage, p.featured, p.created_at, p.updated_at,
			   IFNULL(c.name, 'Uncategorized') as category_name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?`,
		productID).Scan(
		&product.ID, &product.Name, &product.Description, &product.CategoryID,
		&product.BasePrice, &product.DiscountPercentage, &product.Featured,
		&product.CreatedAt, &product.UpdatedAt, &categoryName)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Calculate final price
	finalPrice := product.BasePrice * (1 - product.DiscountPercentage/100)

	// Get all images
	rows, err := database.DB.Query("SELECT id, product_id, image_url, is_primary, created_at FROM product_images WHERE product_id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product images",
		})
	}
	defer rows.Close()

	var images []models.ProductImage
	for rows.Next() {
		var image models.ProductImage
		if err := rows.Scan(&image.ID, &image.ProductID, &image.ImageURL, &image.IsPrimary, &image.CreatedAt); err == nil {
			images = append(images, image)
		}
	}

	// Get all colors
	rows, err = database.DB.Query("SELECT id, product_id, color_name, color_hex, created_at FROM product_colors WHERE product_id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product colors",
		})
	}
	defer rows.Close()

	var colors []models.ProductColor
	for rows.Next() {
		var color models.ProductColor
		if err := rows.Scan(&color.ID, &color.ProductID, &color.ColorName, &color.ColorHex, &color.CreatedAt); err == nil {
			colors = append(colors, color)
		}
	}

	// Get all sizes
	rows, err = database.DB.Query("SELECT id, product_id, size_name, created_at FROM product_sizes WHERE product_id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product sizes",
		})
	}
	defer rows.Close()

	var sizes []models.ProductSize
	for rows.Next() {
		var size models.ProductSize
		if err := rows.Scan(&size.ID, &size.ProductID, &size.SizeName, &size.CreatedAt); err == nil {
			sizes = append(sizes, size)
		}
	}

	// Get inventory
	rows, err = database.DB.Query("SELECT color_id, size_id, quantity FROM product_inventory WHERE product_id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product inventory",
		})
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ColorID, &item.SizeID, &item.Quantity); err == nil {
			inventory = append(inventory, item)
		}
	}

	// Combine all data into product response
	response := models.ProductResponse{
		ID:                 product.ID,
		Name:               product.Name,
		Description:        product.Description,
		CategoryID:         product.CategoryID,
		CategoryName:       categoryName,
		BasePrice:          product.BasePrice,
		DiscountPercentage: product.DiscountPercentage,
		FinalPrice:         finalPrice,
		Featured:           product.Featured,
		Images:             images,
		Colors:             colors,
		Sizes:              sizes,
		Inventory:          inventory,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product": response,
	})
}

// UpdateProduct updates a product
func UpdateProduct(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	id := c.Params("id")
	productID, err := strconv.ParseInt(id, 10, 64)
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
	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.BasePrice <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and price are required",
		})
	}

	// Update the product
	_, err = database.DB.Exec(
		`UPDATE products SET 
			name = ?, 
			description = ?, 
			category_id = ?, 
			base_price = ?, 
			discount_percentage = ?, 
			featured = ?, 
			updated_at = ? 
		WHERE id = ?`,
		req.Name,
		req.Description,
		req.CategoryID,
		req.BasePrice,
		req.DiscountPercentage,
		req.Featured,
		time.Now(),
		productID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

// DeleteProduct deletes a product
func DeleteProduct(c *fiber.Ctx) error {
	// Get the product ID from the URL parameter
	id := c.Params("id")
	productID, err := strconv.ParseInt(id, 10, 64)
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

	// Delete the product
	_, err = database.DB.Exec("DELETE FROM products WHERE id = ?", productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
} 