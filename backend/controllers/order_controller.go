package controllers

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// PlaceOrder creates a new order from the user's cart
func PlaceOrder(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Parse request body
	var req models.OrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.AddressID <= 0 || req.PaymentMethod == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Address ID and payment method are required",
		})
	}

	// Check if address exists and belongs to user
	var addressExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM addresses WHERE id = ? AND user_id = ?)", req.AddressID, userID).Scan(&addressExists)
	if err != nil || !addressExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid address",
		})
	}

	// Check if cart is empty
	var cartCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM cart WHERE user_id = ?", userID).Scan(&cartCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	if cartCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cart is empty",
		})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// Calculate total amount
	var totalAmount float64
	err = tx.QueryRow(`
		SELECT SUM(p.base_price * (1 - p.discount_percentage / 100) * c.quantity) 
		FROM cart c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?`,
		userID).Scan(&totalAmount)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to calculate total amount",
		})
	}

	// Create the order
	result, err := tx.Exec(
		`INSERT INTO orders (user_id, address_id, total_amount, payment_method, payment_status, order_status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, req.AddressID, totalAmount, req.PaymentMethod, "pending", "processing", time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Get the order ID
	orderID, _ := result.LastInsertId()

	// Get cart items
	cartRows, err := tx.Query(`
		SELECT c.product_id, c.color_id, c.size_id, c.quantity, 
			p.base_price * (1 - p.discount_percentage / 100) as price_per_unit
		FROM cart c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?`,
		userID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch cart items",
		})
	}
	defer cartRows.Close()

	// Insert order items and update inventory
	for cartRows.Next() {
		var productID, colorID, sizeID int64
		var quantity int
		var pricePerUnit float64

		err := cartRows.Scan(&productID, &colorID, &sizeID, &quantity, &pricePerUnit)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process cart items",
			})
		}

		// Check inventory again (in case it changed since adding to cart)
		var availableQuantity int
		err = tx.QueryRow(
			"SELECT quantity FROM product_inventory WHERE product_id = ? AND color_id = ? AND size_id = ?",
			productID, colorID, sizeID).Scan(&availableQuantity)
		if err != nil || availableQuantity < quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Not enough inventory for one or more items",
			})
		}

		// Insert order item
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, product_id, color_id, size_id, quantity, price_per_unit) VALUES (?, ?, ?, ?, ?, ?)",
			orderID, productID, colorID, sizeID, quantity, pricePerUnit)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create order item",
			})
		}

		// Update inventory
		_, err = tx.Exec(
			"UPDATE product_inventory SET quantity = quantity - ? WHERE product_id = ? AND color_id = ? AND size_id = ?",
			quantity, productID, colorID, sizeID)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update inventory",
			})
		}
	}

	// Clear the cart
	_, err = tx.Exec("DELETE FROM cart WHERE user_id = ?", userID)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Order placed successfully",
		"order_id":     orderID,
		"total_amount": totalAmount,
	})
}

// GetAllOrders returns all orders for the user
func GetAllOrders(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get the user's role
	role := c.Locals("role").(string)

	// Parse query parameters for pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Prepare base query
	baseQuery := `
		SELECT o.id, o.user_id, o.address_id, o.total_amount, o.payment_method, 
			o.payment_status, o.order_status, o.created_at, o.updated_at,
			u.name as user_name, u.email as user_email
		FROM orders o
		JOIN users u ON o.user_id = u.id
	`

	// Add filter for non-admin users
	var countQuery, selectQuery string
	var queryArgs []interface{}

	if role != "admin" {
		// Regular user can only see their own orders
		countQuery = "SELECT COUNT(*) FROM orders WHERE user_id = ?"
		selectQuery = baseQuery + " WHERE o.user_id = ? ORDER BY o.id DESC LIMIT ? OFFSET ?"
		queryArgs = append(queryArgs, userID, limit, offset)
	} else {
		// Admin can see all orders
		countQuery = "SELECT COUNT(*) FROM orders"
		selectQuery = baseQuery + " ORDER BY o.id DESC LIMIT ? OFFSET ?"
		queryArgs = append(queryArgs, limit, offset)
	}

	// Count total orders for pagination
	var total int
	var err error
	if role != "admin" {
		err = database.DB.QueryRow(countQuery, userID).Scan(&total)
	} else {
		err = database.DB.QueryRow(countQuery).Scan(&total)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Query to get orders
	var rows *sql.Rows
	if role != "admin" {
		rows, err = database.DB.Query(selectQuery, userID, limit, offset)
	} else {
		rows, err = database.DB.Query(selectQuery, limit, offset)
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	// Process results
	type OrderWithUser struct {
		models.Order
		UserName  string `json:"user_name"`
		UserEmail string `json:"user_email"`
	}

	var orders []OrderWithUser
	for rows.Next() {
		var order OrderWithUser
		err := rows.Scan(
			&order.ID, &order.UserID, &order.AddressID, &order.TotalAmount,
			&order.PaymentMethod, &order.PaymentStatus, &order.OrderStatus,
			&order.CreatedAt, &order.UpdatedAt, &order.UserName, &order.UserEmail)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orders,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetOrderByID returns a specific order
func GetOrderByID(c *fiber.Ctx) error {
	// Get user ID from context (set by the Protected middleware)
	userID := c.Locals("userID").(int64)

	// Get the user's role
	role := c.Locals("role").(string)

	// Get order ID from URL parameter
	id := c.Params("id")
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	// Check if order exists and belongs to user (if not admin)
	var order models.Order
	var exists bool
	var query string
	var args []interface{}

	if role != "admin" {
		// Regular users can only view their own orders
		query = `
			SELECT EXISTS(SELECT 1 FROM orders WHERE id = ? AND user_id = ?),
			id, user_id, address_id, total_amount, payment_method, payment_status, order_status, created_at, updated_at
			FROM orders WHERE id = ? AND user_id = ?`
		args = []interface{}{orderID, userID, orderID, userID}
	} else {
		// Admins can view any order
		query = `
			SELECT EXISTS(SELECT 1 FROM orders WHERE id = ?),
			id, user_id, address_id, total_amount, payment_method, payment_status, order_status, created_at, updated_at
			FROM orders WHERE id = ?`
		args = []interface{}{orderID, orderID}
	}

	err = database.DB.QueryRow(query, args...).Scan(
		&exists,
		&order.ID, &order.UserID, &order.AddressID, &order.TotalAmount,
		&order.PaymentMethod, &order.PaymentStatus, &order.OrderStatus,
		&order.CreatedAt, &order.UpdatedAt)

	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Get the address
	var address models.Address
	err = database.DB.QueryRow(`
		SELECT id, user_id, name, street, city, state, postal_code, country, phone, is_default, created_at, updated_at
		FROM addresses
		WHERE id = ?`,
		order.AddressID).Scan(
		&address.ID, &address.UserID, &address.Name, &address.Street, &address.City,
		&address.State, &address.PostalCode, &address.Country, &address.Phone,
		&address.IsDefault, &address.CreatedAt, &address.UpdatedAt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch address",
		})
	}

	// Get order items
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.product_id, oi.color_id, oi.size_id, oi.quantity, oi.price_per_unit,
			p.name, p.description,
			pc.color_name, pc.color_hex,
			ps.size_name,
			(SELECT image_url FROM product_images WHERE product_id = p.id AND is_primary = 1 LIMIT 1) as image_url
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		JOIN product_colors pc ON oi.color_id = pc.id
		JOIN product_sizes ps ON oi.size_id = ps.id
		WHERE oi.order_id = ?`,
		orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer rows.Close()

	var items []models.OrderItemResponse
	for rows.Next() {
		var item models.OrderItemResponse
		var pricePerUnit float64

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.ColorID, &item.SizeID, &item.Quantity, &pricePerUnit,
			&item.ProductName, &item.ProductDescription,
			&item.ColorName, &item.ColorHex,
			&item.SizeName,
			&item.ImageURL)
		if err != nil {
			continue
		}

		item.PricePerUnit = pricePerUnit
		item.SubTotal = pricePerUnit * float64(item.Quantity)
		items = append(items, item)
	}

	// Create order response
	orderResponse := models.OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		Address:       address,
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
		PaymentStatus: order.PaymentStatus,
		OrderStatus:   order.OrderStatus,
		Items:         items,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"order": orderResponse,
	})
}

// UpdateOrderStatus updates an order's status (admin only)
func UpdateOrderStatus(c *fiber.Ctx) error {
	// Get order ID from URL parameter
	id := c.Params("id")
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	// Parse request body
	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.OrderStatus == "" && req.PaymentStatus == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order status or payment status must be provided",
		})
	}

	// Validate order status values
	validOrderStatuses := map[string]bool{
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}
	validPaymentStatuses := map[string]bool{
		"pending":  true,
		"paid":     true,
		"failed":   true,
		"refunded": true,
	}

	if req.OrderStatus != "" && !validOrderStatuses[req.OrderStatus] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order status",
		})
	}

	if req.PaymentStatus != "" && !validPaymentStatuses[req.PaymentStatus] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment status",
		})
	}

	// Check if order exists
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE id = ?)", orderID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Prepare the update query and arguments
	query := "UPDATE orders SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.OrderStatus != "" {
		query += ", order_status = ?"
		args = append(args, req.OrderStatus)
	}

	if req.PaymentStatus != "" {
		query += ", payment_status = ?"
		args = append(args, req.PaymentStatus)
	}

	query += " WHERE id = ?"
	args = append(args, orderID)

	// Update the order
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order updated successfully",
	})
}
