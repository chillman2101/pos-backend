package dto

type TransactionItemRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,gt=0"`
	Price     float64 `json:"price" validate:"required,gte=0"`
}

type TransactionItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id,omitempty"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`
}

type CreateTransactionRequest struct {
	ClientTransactionID string                   `json:"client_transaction_id"` // For offline sync idempotency
	Items               []TransactionItemRequest `json:"items" validate:"required,min=1,dive"`
	PaymentMethod       string                   `json:"payment_method" validate:"required,oneof=cash card qris"`
	CustomerName        string                   `json:"customer_name"`
	DiscountAmount      float64                  `json:"discount_amount" validate:"gte=0"`
	TaxAmount           float64                  `json:"tax_amount" validate:"gte=0"`
	Notes               string                   `json:"notes"`
}

type BulkSyncTransactionRequest struct {
	Transactions []CreateTransactionRequest `json:"transactions" validate:"required,min=1,dive"`
}

type StockWarning struct {
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	SoldQuantity    int    `json:"sold_quantity"`
	AvailableStock  int    `json:"available_stock"`
	Shortage        int    `json:"shortage"`
	Message         string `json:"message"`
}

type TransactionResponse struct {
	ID                string                    `json:"id"`
	TransactionCode   string                    `json:"transaction_code"`
	UserID            string                    `json:"user_id,omitempty"`
	Username          string                    `json:"username,omitempty"`
	Items             []TransactionItemResponse `json:"items"`
	TotalAmount       float64                   `json:"total_amount"`
	DiscountAmount    float64                   `json:"discount_amount"`
	TaxAmount         float64                   `json:"tax_amount"`
	FinalAmount       float64                   `json:"final_amount"`
	PaymentMethod     string                    `json:"payment_method"`
	PaymentStatus     string                    `json:"payment_status"`
	CustomerName      string                    `json:"customer_name,omitempty"`
	Notes             string                    `json:"notes,omitempty"`
	HasStockIssue     bool                      `json:"has_stock_issue"`
	StockIssueDetails string                    `json:"stock_issue_details,omitempty"`
	CreatedAt         string                    `json:"created_at"`
}

type BulkSyncResponse struct {
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
	Warnings     []StockWarning        `json:"warnings,omitempty"`
	Errors       []string              `json:"errors,omitempty"`
	Transactions []TransactionResponse `json:"transactions"`
}
