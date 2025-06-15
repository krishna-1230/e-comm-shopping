package controllers

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RegisterUser handles user registration
func RegisterUser(c *fiber.Ctx) error {
	// Parse request body
	var userRegister models.UserRegister
	if err := c.BodyParser(&userRegister); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if userRegister.Name == "" || userRegister.Email == "" || userRegister.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email, and password are required",
		})
	}

	// Check if user with this email already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", userRegister.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "User with this email already exists",
		})
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(userRegister.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	// Create the user
	result, err := database.DB.Exec(
		"INSERT INTO users (name, email, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		userRegister.Name,
		userRegister.Email,
		hashedPassword,
		"customer", // Default role
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Get the user ID
	userID, _ := result.LastInsertId()

	// Generate JWT token
	token, err := utils.GenerateToken(userID, userRegister.Email, "customer")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return the user and token
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":    userID,
			"name":  userRegister.Name,
			"email": userRegister.Email,
			"role":  "customer",
		},
		"token": token,
	})
}

// LoginUser handles user authentication
func LoginUser(c *fiber.Ctx) error {
	// Parse request body
	var userLogin models.UserLogin
	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if userLogin.Email == "" || userLogin.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Find the user
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email = ?",
		userLogin.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Check password
	if !utils.CheckPasswordHash(userLogin.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Return the user and token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"token":   token,
	})
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser(c *fiber.Ctx) error {
	// Get the user ID from the context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get the user from the database
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Return the user
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user.ToResponse(),
	})
} 