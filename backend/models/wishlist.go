package models

import "time"

// WishlistItem represents an item in the user's wishlist
type WishlistItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

// WishlistItemResponse is the response format for wishlist items with product details
type WishlistItemResponse struct {
	ID                 int64   `json:"id"`
	ProductID          int64   `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	BasePrice          float64 `json:"base_price"`
	DiscountPercentage float64 `json:"discount_percentage"`
	FinalPrice         float64 `json:"final_price"`
	ImageURL           string  `json:"image_url"`
	InStock            bool    `json:"in_stock"`
	CreatedAt          time.Time `json:"created_at"`
} 