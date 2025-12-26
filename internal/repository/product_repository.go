package repository

import (
	"pos-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) FindAll(page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64
	if err := r.db.Model(&domain.Product{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Preload("Category").Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, count, nil
}

func (r *productRepository) FindByCategory(categoryID uuid.UUID, page, limit int) ([]domain.Product, int64, error) {
	var products []domain.Product
	var count int64
	query := r.db.Model(&domain.Product{}).Where("category_id = ?", categoryID)
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("Category").Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, count, nil
}
