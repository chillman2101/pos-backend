package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryID      *uuid.UUID     `gorm:"type:uuid" json:"category_id"`
	Category        *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name            string         `gorm:"not null;size:255" json:"name"`
	SKU             string         `gorm:"uniqueIndex;not null;size:100" json:"sku"`
	Description     string         `gorm:"type:text" json:"description"`
	Price           float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Cost            float64        `gorm:"type:decimal(15,2);default:0" json:"cost"`
	Stock           int            `gorm:"default:0" json:"stock"`
	MinStock        int            `gorm:"default:0" json:"min_stock"`
	StockVersion    int            `gorm:"default:0" json:"stock_version"`
	LastStockUpdate *time.Time     `json:"last_stock_update"`
	ImageURL        string         `gorm:"size:500" json:"image_url"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID                  uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionCode     string            `gorm:"uniqueIndex;not null;size:50" json:"transaction_code"`
	ClientTransactionID string            `gorm:"uniqueIndex;size:100" json:"client_transaction_id"` // For idempotency
	UserID              *uuid.UUID        `gorm:"type:uuid" json:"user_id"`
	User                *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalAmount         float64           `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	DiscountAmount      float64           `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TaxAmount           float64           `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	FinalAmount         float64           `gorm:"type:decimal(15,2);not null" json:"final_amount"`
	PaymentMethod       string            `gorm:"not null;size:50" json:"payment_method"`        // cash, card, qris
	PaymentStatus       string            `gorm:"size:50;default:pending" json:"payment_status"` // pending, completed, cancelled
	CustomerName        string            `gorm:"size:255" json:"customer_name"`
	Notes               string            `gorm:"type:text" json:"notes"`
	Synced              bool              `gorm:"default:false" json:"synced"`
	SyncedAt            *time.Time        `json:"synced_at"`
	HasStockIssue       bool              `gorm:"default:false" json:"has_stock_issue"`
	StockIssueDetails   string            `gorm:"type:text" json:"stock_issue_details"`
	Items               []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	DeletedAt           gorm.DeletedAt    `gorm:"index" json:"-"`
}

type TransactionItem struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID  `gorm:"type:uuid;not null" json:"transaction_id"`
	ProductID     *uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Product       *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName   string     `gorm:"not null;size:255" json:"product_name"`
	ProductPrice  float64    `gorm:"type:decimal(15,2);not null" json:"product_price"`
	Quantity      int        `gorm:"not null" json:"quantity"`
	Subtotal      float64    `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	CreatedAt     time.Time  `json:"created_at"`
}

type InventoryMovement struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID     uuid.UUID  `gorm:"type:uuid;not null" json:"product_id"`
	Product       *Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	MovementType  string     `gorm:"not null;size:50" json:"movement_type"` // in, out, adjustment
	Quantity      int        `gorm:"not null" json:"quantity"`
	ReferenceType string     `gorm:"size:50" json:"reference_type"` // transaction, purchase, adjustment
	ReferenceID   *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	Notes         string     `gorm:"type:text" json:"notes"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User          *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
