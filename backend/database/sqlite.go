package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDatabase initializes the SQLite database
func InitDatabase() {
	// Create database directory if it doesn't exist
	dbDir := "./database/data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	dbPath := filepath.Join(dbDir, "ecommerce.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DB = db
	fmt.Println("Database connection established")

	// Initialize tables
	createTables()
}

// createTables creates all necessary tables for the e-commerce application
func createTables() {
	// Users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'customer',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Addresses table
	createAddressesTable := `
	CREATE TABLE IF NOT EXISTS addresses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		street TEXT NOT NULL,
		city TEXT NOT NULL,
		state TEXT NOT NULL,
		postal_code TEXT NOT NULL,
		country TEXT NOT NULL,
		phone TEXT NOT NULL,
		is_default BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	// Categories table
	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		description TEXT,
		image_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Products table
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		category_id INTEGER,
		base_price REAL NOT NULL,
		discount_percentage REAL DEFAULT 0,
		featured BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);`

	// Product Images table
	createProductImagesTable := `
	CREATE TABLE IF NOT EXISTS product_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		image_url TEXT NOT NULL,
		is_primary BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);`

	// Product Colors table
	createProductColorsTable := `
	CREATE TABLE IF NOT EXISTS product_colors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		color_name TEXT NOT NULL,
		color_hex TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);`

	// Product Sizes table
	createProductSizesTable := `
	CREATE TABLE IF NOT EXISTS product_sizes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		size_name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);`

	// Product Inventory table
	createProductInventoryTable := `
	CREATE TABLE IF NOT EXISTS product_inventory (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		color_id INTEGER NOT NULL,
		size_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL DEFAULT 0,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		FOREIGN KEY (color_id) REFERENCES product_colors(id) ON DELETE CASCADE,
		FOREIGN KEY (size_id) REFERENCES product_sizes(id) ON DELETE CASCADE,
		UNIQUE(product_id, color_id, size_id)
	);`

	// Orders table
	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		address_id INTEGER NOT NULL,
		total_amount REAL NOT NULL,
		payment_method TEXT NOT NULL,
		payment_status TEXT DEFAULT 'pending',
		order_status TEXT DEFAULT 'processing',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE RESTRICT
	);`

	// Order Items table
	createOrderItemsTable := `
	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		color_id INTEGER NOT NULL,
		size_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		price_per_unit REAL NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT,
		FOREIGN KEY (color_id) REFERENCES product_colors(id) ON DELETE RESTRICT,
		FOREIGN KEY (size_id) REFERENCES product_sizes(id) ON DELETE RESTRICT
	);`

	// Cart table
	createCartTable := `
	CREATE TABLE IF NOT EXISTS cart (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		color_id INTEGER NOT NULL,
		size_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		FOREIGN KEY (color_id) REFERENCES product_colors(id) ON DELETE CASCADE,
		FOREIGN KEY (size_id) REFERENCES product_sizes(id) ON DELETE CASCADE,
		UNIQUE(user_id, product_id, color_id, size_id)
	);`

	// Wishlist table
	createWishlistTable := `
	CREATE TABLE IF NOT EXISTS wishlist (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		UNIQUE(user_id, product_id)
	);`

	// Reviews table
	createReviewsTable := `
	CREATE TABLE IF NOT EXISTS reviews (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
		comment TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);`

	// Execute all create table statements
	tables := []string{
		createUsersTable,
		createAddressesTable,
		createCategoriesTable,
		createProductsTable,
		createProductImagesTable,
		createProductColorsTable,
		createProductSizesTable,
		createProductInventoryTable,
		createOrdersTable,
		createOrderItemsTable,
		createCartTable,
		createWishlistTable,
		createReviewsTable,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}

	log.Println("All tables created successfully")
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed")
	}
}
