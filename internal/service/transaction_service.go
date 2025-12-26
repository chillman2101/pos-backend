package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"pos-backend/internal/domain"
	"pos-backend/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionService interface {
	Create(req *dto.CreateTransactionRequest, userID uuid.UUID) (*dto.TransactionResponse, []dto.StockWarning, error)
	BulkSync(req *dto.BulkSyncTransactionRequest, userID uuid.UUID) (*dto.BulkSyncResponse, error)
	GetByID(id string) (*dto.TransactionResponse, error)
	GetAll(page, limit int, filters domain.TransactionFilters) ([]*dto.TransactionResponse, int64, error)
	Cancel(id string) error
}

type transactionService struct {
	transactionRepo domain.TransactionRepository
	productRepo     domain.ProductRepository
	db              *gorm.DB
}

func NewTransactionService(
	transactionRepo domain.TransactionRepository,
	productRepo domain.ProductRepository,
	db *gorm.DB,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		db:              db,
	}
}

func (s *transactionService) Create(req *dto.CreateTransactionRequest, userID uuid.UUID) (*dto.TransactionResponse, []dto.StockWarning, error) {
	// Check for duplicate client transaction ID (idempotency)
	if req.ClientTransactionID != "" {
		existing, _ := s.transactionRepo.FindByClientTransactionID(req.ClientTransactionID)
		if existing != nil {
			return s.toTransactionResponse(existing), nil, nil
		}
	}

	var warnings []dto.StockWarning
	var stockIssueDetails []string

	// Start database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validate and prepare transaction items
	var transactionItems []domain.TransactionItem
	var totalAmount float64

	for _, itemReq := range req.Items {
		productID, err := uuid.Parse(itemReq.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("invalid product ID: %s", itemReq.ProductID)
		}

		// Get product with lock (prevent race condition)
		var product domain.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, productID).Error; err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("product not found: %s", itemReq.ProductID)
		}

		// Check stock availability and create warning if needed
		if product.Stock < itemReq.Quantity {
			shortage := itemReq.Quantity - product.Stock
			warning := dto.StockWarning{
				ProductID:      product.ID.String(),
				ProductName:    product.Name,
				SoldQuantity:   itemReq.Quantity,
				AvailableStock: product.Stock,
				Shortage:       shortage,
				Message:        fmt.Sprintf("Stock shortage detected for %s. Available: %d, Requested: %d, Short: %d", product.Name, product.Stock, itemReq.Quantity, shortage),
			}
			warnings = append(warnings, warning)
			stockIssueDetails = append(stockIssueDetails, warning.Message)
		}

		// Update product stock (allow negative)
		product.Stock -= itemReq.Quantity
		product.StockVersion++
		now := time.Now()
		product.LastStockUpdate = &now

		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("failed to update stock: %v", err)
		}

		// Create transaction item
		subtotal := itemReq.Price * float64(itemReq.Quantity)
		transactionItems = append(transactionItems, domain.TransactionItem{
			ProductID:    &productID,
			ProductName:  product.Name,
			ProductPrice: itemReq.Price,
			Quantity:     itemReq.Quantity,
			Subtotal:     subtotal,
		})

		totalAmount += subtotal

		// Create inventory movement record
		inventoryMovement := domain.InventoryMovement{
			ProductID:     productID,
			MovementType:  "out",
			Quantity:      -itemReq.Quantity, // Negative for outgoing
			ReferenceType: "transaction",
			UserID:        &userID,
		}
		if err := tx.Create(&inventoryMovement).Error; err != nil {
			tx.Rollback()
			return nil, nil, fmt.Errorf("failed to create inventory movement: %v", err)
		}
	}

	// Calculate final amount
	finalAmount := totalAmount - req.DiscountAmount + req.TaxAmount

	// Generate transaction code
	transactionCode := s.generateTransactionCode()

	// Create transaction
	transaction := domain.Transaction{
		TransactionCode:     transactionCode,
		ClientTransactionID: req.ClientTransactionID,
		UserID:              &userID,
		TotalAmount:         totalAmount,
		DiscountAmount:      req.DiscountAmount,
		TaxAmount:           req.TaxAmount,
		FinalAmount:         finalAmount,
		PaymentMethod:       req.PaymentMethod,
		PaymentStatus:       "completed",
		CustomerName:        req.CustomerName,
		Notes:               req.Notes,
		Synced:              true,
		HasStockIssue:       len(warnings) > 0,
		Items:               transactionItems,
	}

	// Set synced timestamp
	now := time.Now()
	transaction.SyncedAt = &now

	// Add stock issue details if any
	if len(stockIssueDetails) > 0 {
		issuesJSON, _ := json.Marshal(stockIssueDetails)
		transaction.StockIssueDetails = string(issuesJSON)
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	// Update inventory movement reference IDs
	for i := range transactionItems {
		transactionItems[i].TransactionID = transaction.ID
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Reload transaction with relations
	createdTransaction, err := s.transactionRepo.FindByID(transaction.ID)
	if err != nil {
		return nil, nil, err
	}

	return s.toTransactionResponse(createdTransaction), warnings, nil
}

func (s *transactionService) BulkSync(req *dto.BulkSyncTransactionRequest, userID uuid.UUID) (*dto.BulkSyncResponse, error) {
	response := &dto.BulkSyncResponse{
		Transactions: []dto.TransactionResponse{},
		Warnings:     []dto.StockWarning{},
		Errors:       []string{},
	}

	for _, txReq := range req.Transactions {
		txResponse, warnings, err := s.Create(&txReq, userID)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, err.Error())
			continue
		}

		response.SuccessCount++
		response.Transactions = append(response.Transactions, *txResponse)

		if len(warnings) > 0 {
			response.Warnings = append(response.Warnings, warnings...)
		}
	}

	return response, nil
}

func (s *transactionService) GetByID(id string) (*dto.TransactionResponse, error) {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid transaction ID format")
	}

	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, err
	}

	return s.toTransactionResponse(transaction), nil
}

func (s *transactionService) GetAll(page, limit int, filters domain.TransactionFilters) ([]*dto.TransactionResponse, int64, error) {
	transactions, totalData, err := s.transactionRepo.FindAll(page, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.TransactionResponse
	for _, transaction := range transactions {
		responses = append(responses, s.toTransactionResponse(&transaction))
	}

	return responses, totalData, nil
}

func (s *transactionService) Cancel(id string) error {
	transactionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid transaction ID format")
	}

	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return err
	}

	if transaction.PaymentStatus == "cancelled" {
		return errors.New("transaction already cancelled")
	}

	// Start database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Restore stock for each item
	for _, item := range transaction.Items {
		if item.ProductID != nil {
			var product domain.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("product not found: %v", err)
			}

			product.Stock += item.Quantity
			product.StockVersion++
			now := time.Now()
			product.LastStockUpdate = &now

			if err := tx.Save(&product).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to restore stock: %v", err)
			}

			// Create inventory movement record
			inventoryMovement := domain.InventoryMovement{
				ProductID:     *item.ProductID,
				MovementType:  "in",
				Quantity:      item.Quantity,
				ReferenceType: "transaction_cancel",
				ReferenceID:   &transaction.ID,
			}
			if err := tx.Create(&inventoryMovement).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create inventory movement: %v", err)
			}
		}
	}

	// Update transaction status
	transaction.PaymentStatus = "cancelled"
	if err := tx.Save(transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to cancel transaction: %v", err)
	}

	return tx.Commit().Error
}

// Helper functions

func (s *transactionService) generateTransactionCode() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	timeStr := now.Format("150405")
	return fmt.Sprintf("TRX-%s-%s", dateStr, timeStr)
}

func (s *transactionService) toTransactionResponse(transaction *domain.Transaction) *dto.TransactionResponse {
	response := &dto.TransactionResponse{
		ID:                transaction.ID.String(),
		TransactionCode:   transaction.TransactionCode,
		TotalAmount:       transaction.TotalAmount,
		DiscountAmount:    transaction.DiscountAmount,
		TaxAmount:         transaction.TaxAmount,
		FinalAmount:       transaction.FinalAmount,
		PaymentMethod:     transaction.PaymentMethod,
		PaymentStatus:     transaction.PaymentStatus,
		CustomerName:      transaction.CustomerName,
		Notes:             transaction.Notes,
		HasStockIssue:     transaction.HasStockIssue,
		StockIssueDetails: transaction.StockIssueDetails,
		CreatedAt:         transaction.CreatedAt.Format(time.RFC3339),
	}

	if transaction.UserID != nil {
		response.UserID = transaction.UserID.String()
		if transaction.User != nil {
			response.Username = transaction.User.Username
		}
	}

	// Convert items
	for _, item := range transaction.Items {
		itemResponse := dto.TransactionItemResponse{
			ID:           item.ID.String(),
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
		}
		if item.ProductID != nil {
			itemResponse.ProductID = item.ProductID.String()
		}
		response.Items = append(response.Items, itemResponse)
	}

	return response
}
