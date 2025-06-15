package models

import "time"

// Product represents a product in the e-commerce system
type Product struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	CategoryID         int64     `json:"category_id"`
	BasePrice          float64   `json:"base_price"`
	DiscountPercentage float64   `json:"discount_percentage"`
	Featured           bool      `json:"featured"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ProductImage represents a product image
type ProductImage struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductColor represents a product color
type ProductColor struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ColorName string    `json:"color_name"`
	ColorHex  string    `json:"color_hex"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductSize represents a product size
type ProductSize struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	SizeName  string    `json:"size_name"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductInventory represents inventory for a specific product variant (color+size)
type ProductInventory struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ColorID   int64     `json:"color_id"`
	SizeID    int64     `json:"size_id"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductResponse represents a product with its associated data
type ProductResponse struct {
	ID                 int64           `json:"id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	CategoryID         int64           `json:"category_id"`
	CategoryName       string          `json:"category_name"`
	BasePrice          float64         `json:"base_price"`
	DiscountPercentage float64         `json:"discount_percentage"`
	FinalPrice         float64         `json:"final_price"`
	Featured           bool            `json:"featured"`
	Images             []ProductImage  `json:"images"`
	Colors             []ProductColor  `json:"colors"`
	Sizes              []ProductSize   `json:"sizes"`
	Inventory          []InventoryItem `json:"inventory"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// InventoryItem represents a simplified inventory item for the response
type InventoryItem struct {
	ColorID  int64 `json:"color_id"`
	SizeID   int64 `json:"size_id"`
	Quantity int   `json:"quantity"`
}

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	CategoryID         int64   `json:"category_id"`
	BasePrice          float64 `json:"base_price"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Featured           bool    `json:"featured"`
}

// Category represents a product category
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
} 