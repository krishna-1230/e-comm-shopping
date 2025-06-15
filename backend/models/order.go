package models

import "time"

// Order represents an order in the system
type Order struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	AddressID     int64     `json:"address_id"`
	TotalAmount   float64   `json:"total_amount"`
	PaymentMethod string    `json:"payment_method"`
	PaymentStatus string    `json:"payment_status"`
	OrderStatus   string    `json:"order_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           int64   `json:"id"`
	OrderID      int64   `json:"order_id"`
	ProductID    int64   `json:"product_id"`
	ColorID      int64   `json:"color_id"`
	SizeID       int64   `json:"size_id"`
	Quantity     int     `json:"quantity"`
	PricePerUnit float64 `json:"price_per_unit"`
}

// OrderItemResponse is the response format for order items with product details
type OrderItemResponse struct {
	ID                 int64   `json:"id"`
	ProductID          int64   `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	ColorID            int64   `json:"color_id"`
	ColorName          string  `json:"color_name"`
	ColorHex           string  `json:"color_hex"`
	SizeID             int64   `json:"size_id"`
	SizeName           string  `json:"size_name"`
	ImageURL           string  `json:"image_url"`
	Quantity           int     `json:"quantity"`
	PricePerUnit       float64 `json:"price_per_unit"`
	SubTotal           float64 `json:"sub_total"`
}

// OrderRequest is the request format for creating an order
type OrderRequest struct {
	AddressID     int64  `json:"address_id"`
	PaymentMethod string `json:"payment_method"`
}

// OrderResponse is the response format for orders with address and items
type OrderResponse struct {
	ID            int64              `json:"id"`
	UserID        int64              `json:"user_id"`
	Address       Address            `json:"address"`
	TotalAmount   float64            `json:"total_amount"`
	PaymentMethod string             `json:"payment_method"`
	PaymentStatus string             `json:"payment_status"`
	OrderStatus   string             `json:"order_status"`
	Items         []OrderItemResponse `json:"items"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// UpdateOrderStatusRequest is the request format for updating order status
type UpdateOrderStatusRequest struct {
	OrderStatus   string `json:"order_status"`
	PaymentStatus string `json:"payment_status"`
}

// Address model represents a user's address
type Address struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Name       string    `json:"name"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	Phone      string    `json:"phone"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AddressRequest is the request format for creating/updating an address
type AddressRequest struct {
	Name       string `json:"name"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Phone      string `json:"phone"`
	IsDefault  bool   `json:"is_default"`
} 