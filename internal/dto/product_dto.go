package dto

type ProductResponse struct {
	ID              string  `json:"id"`
	CategoryID      string  `json:"category_id,omitempty"`
	CategoryName    string  `json:"category_name,omitempty"`
	Name            string  `json:"name"`
	SKU             string  `json:"sku"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Cost            float64 `json:"cost"`
	Stock           int     `json:"stock"`
	MinStock        int     `json:"min_stock"`
	StockVersion    int     `json:"stock_version"`
	LastStockUpdate string  `json:"last_stock_update,omitempty"`
	ImageURL        string  `json:"image_url,omitempty"`
	IsActive        bool    `json:"is_active"`
}

type CreateProductRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Cost        float64 `json:"cost" validate:"gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	MinStock    int     `json:"min_stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name" validate:"required"`
	SKU         string  `json:"sku" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	Cost        float64 `json:"cost" validate:"gte=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	MinStock    int     `json:"min_stock" validate:"gte=0"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}
