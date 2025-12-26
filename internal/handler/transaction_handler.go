package handler

import (
	"math"
	"pos-backend/internal/domain"
	"pos-backend/internal/dto"
	"pos-backend/internal/service"
	"pos-backend/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Get user ID from JWT token
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	transaction, warnings, err := h.transactionService.Create(&req, userID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Return success with warnings if any
	if len(warnings) > 0 {
		response.SuccessWithWarnings(c, "Transaction created with stock warnings", transaction, warnings)
		return
	}

	response.Success(c, "Transaction created successfully", transaction)
}

func (h *TransactionHandler) BulkSync(c *gin.Context) {
	var req dto.BulkSyncTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Get user ID from JWT token
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	result, err := h.transactionService.BulkSync(&req, userID)
	if err != nil {
		response.InternalServerError(c, "Failed to sync transactions", err.Error())
		return
	}

	response.Success(c, "Bulk sync completed", result)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionService.GetByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Transaction retrieved successfully", transaction)
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Build filters
	filters := domain.TransactionFilters{}

	// Filter by user ID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			filters.UserID = &userID
		}
	}

	// Filter by payment method
	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		filters.PaymentMethod = paymentMethod
	}

	// Filter by payment status
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters.PaymentStatus = paymentStatus
	}

	// Filter by date range
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filters.StartDate = &startDate
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.EndDate = &endDate
		}
	}

	// Filter by stock issue
	if hasStockIssueStr := c.Query("has_stock_issue"); hasStockIssueStr != "" {
		hasStockIssue := hasStockIssueStr == "true"
		filters.HasStockIssue = &hasStockIssue
	}

	transactions, total, err := h.transactionService.GetAll(page, limit, filters)
	if err != nil {
		response.InternalServerError(c, "Failed to get transactions", err.Error())
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response.SuccessWithPagination(c, "Transactions retrieved successfully", transactions, response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalRows:  total,
		TotalPages: totalPages,
	})
}

func (h *TransactionHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	if err := h.transactionService.Cancel(id); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Transaction cancelled successfully", nil)
}
