package domain

import "github.com/google/uuid"

type ProductRepository interface {
	Create(product *Product) error
	FindByID(id uuid.UUID) (*Product, error)
	FindBySKU(sku string) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
	FindAll(page, limit int) ([]Product, int64, error)
	FindByCategory(categoryID uuid.UUID, page, limit int) ([]Product, int64, error)
}
