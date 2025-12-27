package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"pos-backend/internal/config"
	"pos-backend/internal/database"
	"pos-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env from root directory
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations first
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("🌱 Starting product seeder...")

	// Open CSV file
	file, err := os.Open("products_data.csv")
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV: %v", err)
	}

	// Skip header row
	records = records[1:]

	successCount := 0
	errorCount := 0

	// Insert each product
	for i, record := range records {
		if len(record) != 10 {
			log.Printf("Row %d: Invalid number of columns, skipping...", i+1)
			errorCount++
			continue
		}

		// Parse data
		categoryID, err := uuid.Parse(record[0])
		if err != nil {
			log.Printf("Row %d: Invalid category_id '%s', skipping...", i+1, record[0])
			errorCount++
			continue
		}

		price, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Row %d: Invalid price '%s', skipping...", i+1, record[4])
			errorCount++
			continue
		}

		cost, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Row %d: Invalid cost '%s', skipping...", i+1, record[5])
			errorCount++
			continue
		}

		stock, err := strconv.Atoi(record[6])
		if err != nil {
			log.Printf("Row %d: Invalid stock '%s', skipping...", i+1, record[6])
			errorCount++
			continue
		}

		minStock, err := strconv.Atoi(record[7])
		if err != nil {
			log.Printf("Row %d: Invalid min_stock '%s', skipping...", i+1, record[7])
			errorCount++
			continue
		}

		isActive := record[9] == "TRUE" || record[9] == "true" || record[9] == "1"

		// Create product
		product := domain.Product{
			ID:          uuid.New(),
			CategoryID:  &categoryID, // Pointer ke uuid
			Name:        record[1],
			SKU:         record[2],
			Description: record[3],
			Price:       price,
			Cost:        cost,
			Stock:       stock,
			MinStock:    minStock,
			ImageURL:    record[8],
			IsActive:    isActive,
		}

		// Insert to database
		result := db.Create(&product)
		if result.Error != nil {
			log.Printf("Row %d: Failed to insert product '%s' (SKU: %s): %v", i+1, product.Name, product.SKU, result.Error)
			errorCount++
			continue
		}

		fmt.Printf("✅ Inserted: %s (SKU: %s)\n", product.Name, product.SKU)
		successCount++
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("✅ Successfully inserted: %d products\n", successCount)
	fmt.Printf("❌ Failed to insert: %d products\n", errorCount)
	fmt.Printf("📊 Total processed: %d products\n", successCount+errorCount)
	fmt.Println(strings.Repeat("=", 60))
}
