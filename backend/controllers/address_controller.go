package controllers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateAddress creates a new address for the user
func CreateAddress(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Parse request body
	var req models.AddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.Street == "" || req.City == "" || req.State == "" || req.PostalCode == "" || req.Country == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All address fields are required",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// If this is the default address, update all other addresses to not be default
	if req.IsDefault {
		_, err = tx.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", userID)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update existing addresses",
			})
		}
	}

	// Create the address
	result, err := tx.Exec(
		`INSERT INTO addresses (user_id, name, street, city, state, postal_code, country, phone, is_default, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Name, req.Street, req.City, req.State, req.PostalCode, req.Country, req.Phone, req.IsDefault, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create address",
		})
	}

	// If no default address yet, make this the default
	if !req.IsDefault {
		var count int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM addresses WHERE user_id = ?", userID).Scan(&count)
		if err == nil && count == 1 {
			addressID, _ := result.LastInsertId()
			_, err = tx.Exec("UPDATE addresses SET is_default = 1 WHERE id = ?", addressID)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to set default address",
				})
			}
			req.IsDefault = true
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Get the address ID
	addressID, _ := result.LastInsertId()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Address created successfully",
		"id":         addressID,
		"is_default": req.IsDefault,
	})
}

// GetAllAddresses returns all addresses for the user
func GetAllAddresses(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Query to get addresses
	rows, err := database.DB.Query(`
		SELECT id, user_id, name, street, city, state, postal_code, country, phone, is_default, created_at, updated_at
		FROM addresses
		WHERE user_id = ?
		ORDER BY is_default DESC, id DESC`,
		userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	addresses := []models.Address{}
	for rows.Next() {
		var address models.Address
		err := rows.Scan(
			&address.ID, &address.UserID, &address.Name, &address.Street, &address.City,
			&address.State, &address.PostalCode, &address.Country, &address.Phone,
			&address.IsDefault, &address.CreatedAt, &address.UpdatedAt)
		if err != nil {
			continue
		}
		addresses = append(addresses, address)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"addresses": addresses,
	})
}

// GetAddress returns a specific address for the user
func GetAddress(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get address ID from URL parameter
	id := c.Params("id")
	addressID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid address ID",
		})
	}

	// Get the address
	var address models.Address
	err = database.DB.QueryRow(`
		SELECT id, user_id, name, street, city, state, postal_code, country, phone, is_default, created_at, updated_at
		FROM addresses
		WHERE id = ? AND user_id = ?`,
		addressID, userID).Scan(
		&address.ID, &address.UserID, &address.Name, &address.Street, &address.City,
		&address.State, &address.PostalCode, &address.Country, &address.Phone,
		&address.IsDefault, &address.CreatedAt, &address.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Address not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"address": address,
	})
}

// UpdateAddress updates a specific address for the user
func UpdateAddress(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get address ID from URL parameter
	id := c.Params("id")
	addressID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid address ID",
		})
	}

	// Check if address exists and belongs to user
	var exists bool
	var currentIsDefault bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM addresses WHERE id = ? AND user_id = ?), is_default FROM addresses WHERE id = ?",
		addressID, userID, addressID).Scan(&exists, &currentIsDefault)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Address not found",
		})
	}

	// Parse request body
	var req models.AddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.Street == "" || req.City == "" || req.State == "" || req.PostalCode == "" || req.Country == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All address fields are required",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// If this address will be the default, update all other addresses to not be default
	if req.IsDefault && !currentIsDefault {
		_, err = tx.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", userID)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update existing addresses",
			})
		}
	} else if !req.IsDefault && currentIsDefault {
		// Check if this is the only address
		var count int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM addresses WHERE user_id = ?", userID).Scan(&count)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to count addresses",
			})
		}
		// If this is the only address, it must remain the default
		if count == 1 {
			req.IsDefault = true
		} else {
			// Make another address the default
			_, err = tx.Exec(
				"UPDATE addresses SET is_default = 1 WHERE user_id = ? AND id != ? ORDER BY id DESC LIMIT 1",
				userID, addressID)
			if err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to set new default address",
				})
			}
		}
	}

	// Update the address
	_, err = tx.Exec(
		`UPDATE addresses SET 
			name = ?, 
			street = ?, 
			city = ?, 
			state = ?, 
			postal_code = ?, 
			country = ?, 
			phone = ?, 
			is_default = ?, 
			updated_at = ? 
		WHERE id = ? AND user_id = ?`,
		req.Name, req.Street, req.City, req.State, req.PostalCode, req.Country, req.Phone, req.IsDefault, time.Now(),
		addressID, userID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update address",
		})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Address updated successfully",
		"is_default": req.IsDefault,
	})
}

// DeleteAddress deletes a specific address for the user
func DeleteAddress(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get address ID from URL parameter
	id := c.Params("id")
	addressID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid address ID",
		})
	}

	// Check if address exists and belongs to user
	var exists bool
	var isDefault bool
	err = database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM addresses WHERE id = ? AND user_id = ?), is_default FROM addresses WHERE id = ?",
		addressID, userID, addressID).Scan(&exists, &isDefault)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Address not found",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// Delete the address
	_, err = tx.Exec("DELETE FROM addresses WHERE id = ? AND user_id = ?", addressID, userID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete address",
		})
	}

	// If this was the default address, make another address the default
	if isDefault {
		_, err = tx.Exec(
			"UPDATE addresses SET is_default = 1 WHERE user_id = ? ORDER BY id DESC LIMIT 1",
			userID)
		// It's ok if this fails (e.g. if there are no more addresses)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Address deleted successfully",
	})
}
