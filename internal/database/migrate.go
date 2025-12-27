package database

import (
	"log"
	"pos-backend/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running auto migrations...")

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.Transaction{},
		&domain.TransactionItem{},
		&domain.InventoryMovement{},
		&domain.Setting{},
	)

	if err != nil {
		return err
	}

	log.Println("Auto migrations completed successfully")
	return nil
}
