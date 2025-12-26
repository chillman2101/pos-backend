package repository

import (
	"pos-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Preload("Items.Product").Preload("User").First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByTransactionCode(code string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Preload("Items.Product").Preload("User").Where("transaction_code = ?", code).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByClientTransactionID(clientTxID string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.Where("client_transaction_id = ?", clientTxID).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindAll(page, limit int, filters domain.TransactionFilters) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var count int64

	query := r.db.Model(&domain.Transaction{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.PaymentMethod != "" {
		query = query.Where("payment_method = ?", filters.PaymentMethod)
	}
	if filters.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filters.PaymentStatus)
	}
	if filters.StartDate != nil {
		query = query.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}
	if filters.HasStockIssue != nil {
		query = query.Where("has_stock_issue = ?", *filters.HasStockIssue)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items.Product").Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, count, nil
}

func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Transaction{}, id).Error
}
