package models

import "time"

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	ColorID   int64     `json:"color_id"`
	SizeID    int64     `json:"size_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CartItemRequest is the request format for adding/updating cart items
type CartItemRequest struct {
	ProductID int64 `json:"product_id"`
	ColorID   int64 `json:"color_id"`
	SizeID    int64 `json:"size_id"`
	Quantity  int   `json:"quantity"`
}

// CartItemResponse is the response format for cart items with product details
type CartItemResponse struct {
	ID                 int64   `json:"id"`
	ProductID          int64   `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	BasePrice          float64 `json:"base_price"`
	DiscountPercentage float64 `json:"discount_percentage"`
	FinalPrice         float64 `json:"final_price"`
	ColorID            int64   `json:"color_id"`
	ColorName          string  `json:"color_name"`
	ColorHex           string  `json:"color_hex"`
	SizeID             int64   `json:"size_id"`
	SizeName           string  `json:"size_name"`
	ImageURL           string  `json:"image_url"`
	Quantity           int     `json:"quantity"`
	InStock            int     `json:"in_stock"`
	SubTotal           float64 `json:"sub_total"`
}

// CartSummary represents a summary of the cart
type CartSummary struct {
	TotalItems     int     `json:"total_items"`
	SubTotal       float64 `json:"sub_total"`
	ShippingCost   float64 `json:"shipping_cost"`
	Tax            float64 `json:"tax"`
	Total          float64 `json:"total"`
	DiscountAmount float64 `json:"discount_amount"`
} 