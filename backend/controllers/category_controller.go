package controllers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory creates a new product category
func CreateCategory(c *fiber.Ctx) error {
	// Parse request body
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}

	// Check if category with this name already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", req.Name).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Category with this name already exists",
		})
	}

	// Create the category
	result, err := database.DB.Exec(
		"INSERT INTO categories (name, description, image_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		req.Name,
		req.Description,
		req.ImageURL,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	// Get the category ID
	categoryID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
		"id":      categoryID,
	})
}

// GetAllCategories returns all product categories
func GetAllCategories(c *fiber.Ctx) error {
	// Query to get categories
	rows, err := database.DB.Query(`
		SELECT id, name, description, image_url, created_at, updated_at 
		FROM categories 
		ORDER BY name ASC`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description,
			&category.ImageURL, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			continue
		}
		categories = append(categories, category)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": categories,
	})
}

// GetCategoryByID returns a specific category by ID
func GetCategoryByID(c *fiber.Ctx) error {
	// Get the category ID from the URL parameter
	id := c.Params("id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Get the category from the database
	var category models.Category
	err = database.DB.QueryRow(`
		SELECT id, name, description, image_url, created_at, updated_at 
		FROM categories 
		WHERE id = ?`,
		categoryID).Scan(
		&category.ID, &category.Name, &category.Description,
		&category.ImageURL, &category.CreatedAt, &category.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Count products in this category
	var productCount int
	database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE category_id = ?", categoryID).Scan(&productCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"category": category,
		"product_count": productCount,
	})
}

// UpdateCategory updates a category
func UpdateCategory(c *fiber.Ctx) error {
	// Get the category ID from the URL parameter
	id := c.Params("id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Check if category exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Parse request body
	var req models.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}

	// Check if another category with this name already exists
	var nameExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ? AND id != ?)", req.Name, categoryID).Scan(&nameExists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if nameExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Another category with this name already exists",
		})
	}

	// Update the category
	_, err = database.DB.Exec(
		"UPDATE categories SET name = ?, description = ?, image_url = ?, updated_at = ? WHERE id = ?",
		req.Name,
		req.Description,
		req.ImageURL,
		time.Now(),
		categoryID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}

// DeleteCategory deletes a category
func DeleteCategory(c *fiber.Ctx) error {
	// Get the category ID from the URL parameter
	id := c.Params("id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	// Check if category exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", categoryID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Check if there are products in this category
	var productCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM products WHERE category_id = ?", categoryID).Scan(&productCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// If there are products, update their category to NULL
	if productCount > 0 {
		_, err = database.DB.Exec("UPDATE products SET category_id = NULL WHERE category_id = ?", categoryID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update products before category deletion",
			})
		}
	}

	// Delete the category
	_, err = database.DB.Exec("DELETE FROM categories WHERE id = ?", categoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
		"products_affected": productCount,
	})
} 