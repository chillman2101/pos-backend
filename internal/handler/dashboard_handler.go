package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TodaySales          float64 `json:"today_sales"`
	YesterdaySales      float64 `json:"yesterday_sales"`
	SalesChange         float64 `json:"sales_change"`
	TodayTransactions   int64   `json:"today_transactions"`
	YesterdayTxns       int64   `json:"yesterday_transactions"`
	TransactionsChange  float64 `json:"transactions_change"`
	TodayProductsSold   int64   `json:"today_products_sold"`
	YesterdayProdsSold  int64   `json:"yesterday_products_sold"`
	ProductsSoldChange  float64 `json:"products_sold_change"`
	AvgPerTransaction   float64 `json:"avg_per_transaction"`
	YesterdayAvg        float64 `json:"yesterday_avg"`
	AvgChange           float64 `json:"avg_change"`
}

// RecentTransaction represents a recent transaction summary
type RecentTransaction struct {
	ID              string    `json:"id"`
	TransactionCode string    `json:"transaction_code"`
	CustomerName    string    `json:"customer_name"`
	FinalAmount     float64   `json:"final_amount"`
	PaymentStatus   string    `json:"payment_status"`
	CreatedAt       time.Time `json:"created_at"`
}

// LowStockProduct represents a product with low stock
type LowStockProduct struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Stock int    `json:"stock"`
}

// GetDashboardStats returns dashboard statistics
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	yesterdayStart := todayStart.Add(-24 * time.Hour)
	yesterdayEnd := todayStart

	stats := DashboardStats{}

	// Today's sales
	h.db.Table("transactions").
		Where("created_at >= ? AND created_at < ? AND payment_status = ?", todayStart, todayEnd, "completed").
		Select("COALESCE(SUM(final_amount), 0)").
		Scan(&stats.TodaySales)

	// Yesterday's sales
	h.db.Table("transactions").
		Where("created_at >= ? AND created_at < ? AND payment_status = ?", yesterdayStart, yesterdayEnd, "completed").
		Select("COALESCE(SUM(final_amount), 0)").
		Scan(&stats.YesterdaySales)

	// Calculate sales change percentage
	if stats.YesterdaySales > 0 {
		stats.SalesChange = ((stats.TodaySales - stats.YesterdaySales) / stats.YesterdaySales) * 100
	}

	// Today's transactions count
	h.db.Table("transactions").
		Where("created_at >= ? AND created_at < ? AND payment_status = ?", todayStart, todayEnd, "completed").
		Count(&stats.TodayTransactions)

	// Yesterday's transactions count
	h.db.Table("transactions").
		Where("created_at >= ? AND created_at < ? AND payment_status = ?", yesterdayStart, yesterdayEnd, "completed").
		Count(&stats.YesterdayTxns)

	// Calculate transactions change percentage
	if stats.YesterdayTxns > 0 {
		stats.TransactionsChange = ((float64(stats.TodayTransactions) - float64(stats.YesterdayTxns)) / float64(stats.YesterdayTxns)) * 100
	}

	// Today's products sold (sum of quantities from transaction_items)
	h.db.Table("transaction_items").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.created_at >= ? AND transactions.created_at < ? AND transactions.payment_status = ?", todayStart, todayEnd, "completed").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&stats.TodayProductsSold)

	// Yesterday's products sold
	h.db.Table("transaction_items").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.created_at >= ? AND transactions.created_at < ? AND transactions.payment_status = ?", yesterdayStart, yesterdayEnd, "completed").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&stats.YesterdayProdsSold)

	// Calculate products sold change percentage
	if stats.YesterdayProdsSold > 0 {
		stats.ProductsSoldChange = ((float64(stats.TodayProductsSold) - float64(stats.YesterdayProdsSold)) / float64(stats.YesterdayProdsSold)) * 100
	}

	// Average per transaction (today)
	if stats.TodayTransactions > 0 {
		stats.AvgPerTransaction = stats.TodaySales / float64(stats.TodayTransactions)
	}

	// Average per transaction (yesterday)
	if stats.YesterdayTxns > 0 {
		stats.YesterdayAvg = stats.YesterdaySales / float64(stats.YesterdayTxns)
	}

	// Calculate average change percentage
	if stats.YesterdayAvg > 0 {
		stats.AvgChange = ((stats.AvgPerTransaction - stats.YesterdayAvg) / stats.YesterdayAvg) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetRecentTransactions returns recent transactions (today)
func (h *DashboardHandler) GetRecentTransactions(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var transactions []RecentTransaction

	err := h.db.Table("transactions").
		Where("created_at >= ? AND payment_status = ?", todayStart, "completed").
		Order("created_at DESC").
		Limit(5).
		Select("id, transaction_code, customer_name, final_amount, payment_status, created_at").
		Scan(&transactions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch recent transactions",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
	})
}

// GetLowStockProducts returns products with low stock (stock <= 10)
func (h *DashboardHandler) GetLowStockProducts(c *gin.Context) {
	var products []LowStockProduct

	err := h.db.Table("products").
		Where("stock <= ? AND deleted_at IS NULL", 10).
		Order("stock ASC").
		Limit(10).
		Select("id, name, sku, stock").
		Scan(&products).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch low stock products",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    products,
	})
}
