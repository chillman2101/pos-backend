package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportsHandler struct {
	db *gorm.DB
}

func NewReportsHandler(db *gorm.DB) *ReportsHandler {
	return &ReportsHandler{db: db}
}

// Sales Summary Response
type SalesSummaryResponse struct {
	TotalRevenue       float64 `json:"total_revenue"`
	TotalTransactions  int64   `json:"total_transactions"`
	TotalProductsSold  int64   `json:"total_products_sold"`
	AverageOrderValue  float64 `json:"average_order_value"`
}

// Top Product Response
type TopProductResponse struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	SKU           string  `json:"sku"`
	TotalQuantity int64   `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
	AvgPrice      float64 `json:"avg_price"`
}

// Sales by Payment Method Response
type PaymentMethodResponse struct {
	PaymentMethod    string  `json:"payment_method"`
	TotalAmount      float64 `json:"total_amount"`
	TransactionCount int64   `json:"transaction_count"`
	Percentage       float64 `json:"percentage"`
}

// Daily Sales Response
type DailySalesResponse struct {
	Date             string  `json:"date"`
	TotalRevenue     float64 `json:"total_revenue"`
	TransactionCount int64   `json:"transaction_count"`
}

// GetSalesSummary returns overall sales summary
func (h *ReportsHandler) GetSalesSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := h.db.Table("transactions").
		Where("payment_status = ?", "completed").
		Where("deleted_at IS NULL")

	// Apply date filters
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	var summary SalesSummaryResponse

	// Get total revenue and transaction count
	var result struct {
		TotalRevenue      float64
		TotalTransactions int64
	}

	query.Select("COALESCE(SUM(final_amount), 0) as total_revenue, COUNT(*) as total_transactions").
		Scan(&result)

	summary.TotalRevenue = result.TotalRevenue
	summary.TotalTransactions = result.TotalTransactions

	// Get total products sold
	itemQuery := h.db.Table("transaction_items").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.payment_status = ?", "completed").
		Where("transactions.deleted_at IS NULL")

	if startDate != "" {
		itemQuery = itemQuery.Where("DATE(transactions.created_at) >= ?", startDate)
	}
	if endDate != "" {
		itemQuery = itemQuery.Where("DATE(transactions.created_at) <= ?", endDate)
	}

	itemQuery.Select("COALESCE(SUM(quantity), 0)").Scan(&summary.TotalProductsSold)

	// Calculate average order value
	if summary.TotalTransactions > 0 {
		summary.AverageOrderValue = summary.TotalRevenue / float64(summary.TotalTransactions)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetTopProducts returns best selling products
func (h *ReportsHandler) GetTopProducts(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit := c.DefaultQuery("limit", "10")

	query := h.db.Table("transaction_items").
		Select(`
			transaction_items.product_id,
			transaction_items.product_name,
			products.sku,
			SUM(transaction_items.quantity) as total_quantity,
			SUM(transaction_items.subtotal) as total_revenue,
			AVG(transaction_items.product_price) as avg_price
		`).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("LEFT JOIN products ON products.id = transaction_items.product_id").
		Where("transactions.payment_status = ?", "completed").
		Where("transactions.deleted_at IS NULL").
		Group("transaction_items.product_id, transaction_items.product_name, products.sku").
		Order("total_quantity DESC").
		Limit(10)

	// Apply date filters
	if startDate != "" {
		query = query.Where("DATE(transactions.created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(transactions.created_at) <= ?", endDate)
	}

	// Override limit if provided
	if limit != "10" {
		query = query.Limit(10) // You can parse limit string to int if needed
	}

	var topProducts []TopProductResponse
	if err := query.Scan(&topProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch top products",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    topProducts,
	})
}

// GetSalesByPaymentMethod returns sales breakdown by payment method
func (h *ReportsHandler) GetSalesByPaymentMethod(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := h.db.Table("transactions").
		Select(`
			payment_method,
			SUM(final_amount) as total_amount,
			COUNT(*) as transaction_count
		`).
		Where("payment_status = ?", "completed").
		Where("deleted_at IS NULL").
		Group("payment_method").
		Order("total_amount DESC")

	// Apply date filters
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
	}

	var paymentMethods []PaymentMethodResponse
	if err := query.Scan(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch payment methods data",
			"error":   err.Error(),
		})
		return
	}

	// Calculate total for percentage
	var total float64
	for _, pm := range paymentMethods {
		total += pm.TotalAmount
	}

	// Calculate percentage
	for i := range paymentMethods {
		if total > 0 {
			paymentMethods[i].Percentage = (paymentMethods[i].TotalAmount / total) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    paymentMethods,
	})
}

// GetDailySales returns daily sales for charts
func (h *ReportsHandler) GetDailySales(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Default to last 30 days if no dates provided
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	query := h.db.Table("transactions").
		Select(`
			DATE(created_at) as date,
			SUM(final_amount) as total_revenue,
			COUNT(*) as transaction_count
		`).
		Where("payment_status = ?", "completed").
		Where("deleted_at IS NULL").
		Where("DATE(created_at) >= ?", startDate).
		Where("DATE(created_at) <= ?", endDate).
		Group("DATE(created_at)").
		Order("date ASC")

	var dailySales []DailySalesResponse
	if err := query.Scan(&dailySales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch daily sales",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dailySales,
	})
}
