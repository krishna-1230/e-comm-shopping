package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	// Open database connection
	dbDir := "../database/data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		fmt.Println("Database directory not found. Make sure the database is initialized.")
		return
	}

	dbPath := filepath.Join(dbDir, "ecommerce.db")
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database successfully")

	// Start seeding data
	fmt.Println("Starting data seeding...")

	// Seed users
	seedUsers()

	// Seed categories
	seedCategories()

	// Seed products
	seedProducts()

	fmt.Println("Data seeding completed successfully!")
}

func seedUsers() {
	fmt.Println("Seeding users...")

	// Clear existing users
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Warning: Failed to clear users table: %v", err)
	}

	// Create admin user
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err = db.Exec(
		"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
		"Admin User", "admin@example.com", string(adminPassword), "admin",
	)
	if err != nil {
		log.Printf("Failed to create admin user: %v", err)
	} else {
		fmt.Println("Admin user created")
	}

	// Create regular users
	users := []struct {
		name     string
		email    string
		password string
	}{
		{"John Doe", "john@example.com", "password123"},
		{"Jane Smith", "jane@example.com", "password123"},
		{"Bob Johnson", "bob@example.com", "password123"},
	}

	for _, user := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		_, err := db.Exec(
			"INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)",
			user.name, user.email, string(hashedPassword), "customer",
		)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.email, err)
		} else {
			fmt.Printf("User created: %s\n", user.email)
		}
	}
}

func seedCategories() {
	fmt.Println("Seeding categories...")

	// Clear existing categories
	_, err := db.Exec("DELETE FROM categories")
	if err != nil {
		log.Printf("Warning: Failed to clear categories table: %v", err)
	}

	categories := []struct {
		name        string
		description string
		imageUrl    string
	}{
		{
			"Men's Clothing",
			"Quality clothing for men including shirts, trousers, and jackets",
			"https://images.unsplash.com/photo-1602810318383-e386cc2a3ccf",
		},
		{
			"Women's Clothing",
			"Stylish clothing for women including dresses, tops, and skirts",
			"https://images.unsplash.com/photo-1567401893414-76b7b1e5a7a5",
		},
		{
			"Accessories",
			"Fashion accessories including bags, hats, and jewelry",
			"https://images.unsplash.com/photo-1576053139778-7e32f2ae3cfd",
		},
		{
			"Footwear",
			"Quality footwear for all occasions",
			"https://images.unsplash.com/photo-1549298916-b41d501d3772",
		},
	}

	for _, category := range categories {
		_, err := db.Exec(
			"INSERT INTO categories (name, description, image_url) VALUES (?, ?, ?)",
			category.name, category.description, category.imageUrl,
		)
		if err != nil {
			log.Printf("Failed to create category %s: %v", category.name, err)
		} else {
			fmt.Printf("Category created: %s\n", category.name)
		}
	}
}

func seedProducts() {
	fmt.Println("Seeding products...")

	// Clear existing product data
	tables := []string{
		"product_inventory",
		"product_sizes",
		"product_colors",
		"product_images",
		"products",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Warning: Failed to clear %s table: %v", table, err)
		}
	}

	// Get category IDs
	categoryIDs := make(map[string]int)
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		log.Fatalf("Failed to get categories: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("Failed to scan category row: %v", err)
		}
		categoryIDs[name] = id
	}

	// Seed products
	products := []struct {
		name              string
		description       string
		categoryName      string
		price             float64
		discountPercent   float64
		featured          bool
		colors            []string
		colorHexes        []string
		sizes             []string
		images            []string
		primaryImageIndex int
	}{
		{
			name:            "Classic White T-Shirt",
			description:     "A comfortable white t-shirt made from 100% cotton. Perfect for everyday casual wear.",
			categoryName:    "Men's Clothing",
			price:           24.99,
			discountPercent: 0,
			featured:        true,
			colors:          []string{"White", "Black", "Gray"},
			colorHexes:      []string{"#FFFFFF", "#000000", "#808080"},
			sizes:           []string{"S", "M", "L", "XL"},
			images: []string{
				"https://images.unsplash.com/photo-1521572163474-6864f9cf17ab",
				"https://images.unsplash.com/photo-1622445275576-721325763afe",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Summer Floral Dress",
			description:     "A beautiful floral summer dress perfect for sunny days and casual outings.",
			categoryName:    "Women's Clothing",
			price:           49.99,
			discountPercent: 10,
			featured:        true,
			colors:          []string{"Blue", "Red"},
			colorHexes:      []string{"#0000FF", "#FF0000"},
			sizes:           []string{"S", "M", "L"},
			images: []string{
				"https://images.unsplash.com/photo-1612722432474-b971cdcea546",
				"https://images.unsplash.com/photo-1583496661160-fb5886a773ba",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Casual Denim Jacket",
			description:     "A classic denim jacket that never goes out of style. Perfect for layering in any season.",
			categoryName:    "Men's Clothing",
			price:           79.99,
			discountPercent: 0,
			featured:        true,
			colors:          []string{"Blue"},
			colorHexes:      []string{"#0000AA"},
			sizes:           []string{"M", "L", "XL"},
			images: []string{
				"https://images.unsplash.com/photo-1601333144130-8cbb312386b6",
				"https://images.unsplash.com/photo-1542272604-787c3835535d",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Leather Crossbody Bag",
			description:     "A stylish leather crossbody bag with multiple compartments. Perfect for keeping your essentials organized.",
			categoryName:    "Accessories",
			price:           89.99,
			discountPercent: 15,
			featured:        true,
			colors:          []string{"Brown", "Black"},
			colorHexes:      []string{"#964B00", "#000000"},
			sizes:           []string{},
			images: []string{
				"https://images.unsplash.com/photo-1598532163257-ae3c6b2524b6",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Running Sneakers",
			description:     "Lightweight and comfortable running sneakers with excellent cushioning and support.",
			categoryName:    "Footwear",
			price:           119.99,
			discountPercent: 0,
			featured:        true,
			colors:          []string{"White", "Black", "Red"},
			colorHexes:      []string{"#FFFFFF", "#000000", "#FF0000"},
			sizes:           []string{"7", "8", "9", "10", "11"},
			images: []string{
				"https://images.unsplash.com/photo-1597248881519-d8a9190fa51c",
				"https://images.unsplash.com/photo-1607522370275-f14206abe5d3",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Slim Fit Chinos",
			description:     "Classic slim fit chinos for a smart casual look. Made from comfortable and durable cotton blend.",
			categoryName:    "Men's Clothing",
			price:           59.99,
			discountPercent: 0,
			featured:        true,
			colors:          []string{"Beige", "Navy", "Olive"},
			colorHexes:      []string{"#F5F5DC", "#000080", "#808000"},
			sizes:           []string{"28", "30", "32", "34", "36"},
			images: []string{
				"https://images.unsplash.com/photo-1473966968600-fa801b869a1a",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Oversized Knit Sweater",
			description:     "A cozy oversized knit sweater perfect for chilly days. Features a stylish pattern and ribbed cuffs.",
			categoryName:    "Women's Clothing",
			price:           64.99,
			discountPercent: 20,
			featured:        true,
			colors:          []string{"Cream", "Gray", "Pink"},
			colorHexes:      []string{"#FFFDD0", "#808080", "#FFC0CB"},
			sizes:           []string{"S", "M", "L"},
			images: []string{
				"https://images.unsplash.com/photo-1434389677669-e08b4cac3105",
			},
			primaryImageIndex: 0,
		},
		{
			name:            "Stainless Steel Watch",
			description:     "An elegant stainless steel watch with a classic design. Water resistant and built to last.",
			categoryName:    "Accessories",
			price:           149.99,
			discountPercent: 0,
			featured:        true,
			colors:          []string{"Silver", "Gold", "Rose Gold"},
			colorHexes:      []string{"#C0C0C0", "#FFD700", "#B76E79"},
			sizes:           []string{},
			images: []string{
				"https://images.unsplash.com/photo-1524592094714-0f0654e20314",
				"https://images.unsplash.com/photo-1522312346375-d1a52e2b99b3",
			},
			primaryImageIndex: 0,
		},
	}

	for _, product := range products {
		// Insert product
		categoryID := categoryIDs[product.categoryName]

		result, err := db.Exec(
			"INSERT INTO products (name, description, category_id, base_price, discount_percentage, featured) VALUES (?, ?, ?, ?, ?, ?)",
			product.name, product.description, categoryID, product.price, product.discountPercent, product.featured,
		)
		if err != nil {
			log.Printf("Failed to create product %s: %v", product.name, err)
			continue
		}

		productID, _ := result.LastInsertId()
		fmt.Printf("Product created: %s (ID: %d)\n", product.name, productID)

		// Insert product images
		for i, imageUrl := range product.images {
			isPrimary := i == product.primaryImageIndex
			_, err := db.Exec(
				"INSERT INTO product_images (product_id, image_url, is_primary) VALUES (?, ?, ?)",
				productID, imageUrl, isPrimary,
			)
			if err != nil {
				log.Printf("Failed to add image for product %s: %v", product.name, err)
			}
		}

		// Insert product colors
		colorIDs := make([]int64, len(product.colors))
		for i, color := range product.colors {
			result, err := db.Exec(
				"INSERT INTO product_colors (product_id, color_name, color_hex) VALUES (?, ?, ?)",
				productID, color, product.colorHexes[i],
			)
			if err != nil {
				log.Printf("Failed to add color %s for product %s: %v", color, product.name, err)
				continue
			}
			colorID, _ := result.LastInsertId()
			colorIDs[i] = colorID
		}

		// Insert product sizes
		sizeIDs := make([]int64, len(product.sizes))
		for i, size := range product.sizes {
			result, err := db.Exec(
				"INSERT INTO product_sizes (product_id, size_name) VALUES (?, ?)",
				productID, size,
			)
			if err != nil {
				log.Printf("Failed to add size %s for product %s: %v", size, product.name, err)
				continue
			}
			sizeID, _ := result.LastInsertId()
			sizeIDs[i] = sizeID
		}

		// Insert product inventory
		if len(colorIDs) > 0 && len(sizeIDs) > 0 {
			for _, colorID := range colorIDs {
				for _, sizeID := range sizeIDs {
					_, err := db.Exec(
						"INSERT INTO product_inventory (product_id, color_id, size_id, quantity) VALUES (?, ?, ?, ?)",
						productID, colorID, sizeID, 100, // Default quantity of 100
					)
					if err != nil {
						log.Printf("Failed to add inventory for product %s: %v", product.name, err)
					}
				}
			}
		}
	}
}
