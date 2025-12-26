package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(transaction *Transaction) error
	FindByID(id uuid.UUID) (*Transaction, error)
	FindByTransactionCode(code string) (*Transaction, error)
	FindByClientTransactionID(clientTxID string) (*Transaction, error)
	FindAll(page, limit int, filters TransactionFilters) ([]Transaction, int64, error)
	Update(transaction *Transaction) error
	Delete(id uuid.UUID) error
}

type TransactionFilters struct {
	UserID        *uuid.UUID
	PaymentMethod string
	PaymentStatus string
	StartDate     *time.Time
	EndDate       *time.Time
	HasStockIssue *bool
}
