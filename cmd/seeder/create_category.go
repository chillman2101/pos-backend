package main

import (
	"fmt"
	"log"

	"pos-backend/internal/config"
	"pos-backend/internal/database"
	"pos-backend/internal/domain"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create default category
	categoryID := uuid.MustParse("0383a023-96b6-4d52-bfb2-332cb6bec573")
	category := domain.Category{
		ID:          categoryID,
		Name:        "General Products",
		Description: "Kategori umum untuk semua produk",
	}

	result := db.FirstOrCreate(&category, domain.Category{ID: categoryID})
	if result.Error != nil {
		log.Fatalf("Failed to create category: %v", result.Error)
	}

	fmt.Println("✅ Category 'General Products' ready!")
	fmt.Printf("   ID: %s\n", category.ID)
}
