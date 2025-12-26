package service

import (
	"errors"
	"pos-backend/internal/domain"
	"pos-backend/internal/dto"

	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(page, limit int) ([]*dto.ProductResponse, int64, error)
	GetByCategory(categoryID string, page, limit int) ([]*dto.ProductResponse, int64, error)
	Create(req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetByID(id string) (*dto.ProductResponse, error)
	GetBySKU(sku string) (*dto.ProductResponse, error)
	Update(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	Delete(id string) error
}

type productService struct {
	productRepo domain.ProductRepository
}

func NewProductService(productRepo domain.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) Create(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Check if SKU already exists
	existingProduct, _ := s.productRepo.FindBySKU(req.SKU)
	if existingProduct != nil {
		return nil, errors.New("product with this SKU already exists")
	}

	product := domain.Product{
		Name:        req.Name,
		SKU:         req.SKU,
		Description: req.Description,
		Price:       req.Price,
		Cost:        req.Cost,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	}

	// Set category ID if provided
	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID format")
		}
		product.CategoryID = &categoryID
	}

	err := s.productRepo.Create(&product)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(&product), nil
}

func (s *productService) GetByID(id string) (*dto.ProductResponse, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) GetBySKU(sku string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindBySKU(sku)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) Update(id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product ID format")
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, err
	}

	// Check if SKU is being changed and if it already exists
	if product.SKU != req.SKU {
		existingProduct, _ := s.productRepo.FindBySKU(req.SKU)
		if existingProduct != nil {
			return nil, errors.New("product with this SKU already exists")
		}
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.Description = req.Description
	product.Price = req.Price
	product.Cost = req.Cost
	product.Stock = req.Stock
	product.MinStock = req.MinStock
	product.ImageURL = req.ImageURL
	product.IsActive = req.IsActive

	// Update category ID if provided
	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, errors.New("invalid category ID format")
		}
		product.CategoryID = &categoryID
	} else {
		product.CategoryID = nil
	}

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

func (s *productService) Delete(id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product ID format")
	}
	return s.productRepo.Delete(productID)
}

func (s *productService) GetAll(page, limit int) ([]*dto.ProductResponse, int64, error) {
	products, totalData, err := s.productRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ProductResponse
	for _, product := range products {
		responses = append(responses, s.toProductResponse(&product))
	}

	return responses, totalData, nil
}

func (s *productService) GetByCategory(categoryID string, page, limit int) ([]*dto.ProductResponse, int64, error) {
	catID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, 0, errors.New("invalid category ID format")
	}

	products, totalData, err := s.productRepo.FindByCategory(catID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.ProductResponse
	for _, product := range products {
		responses = append(responses, s.toProductResponse(&product))
	}

	return responses, totalData, nil
}

// Helper function to convert domain.Product to dto.ProductResponse
func (s *productService) toProductResponse(product *domain.Product) *dto.ProductResponse {
	response := &dto.ProductResponse{
		ID:           product.ID.String(),
		Name:         product.Name,
		SKU:          product.SKU,
		Description:  product.Description,
		Price:        product.Price,
		Cost:         product.Cost,
		Stock:        product.Stock,
		MinStock:     product.MinStock,
		StockVersion: product.StockVersion,
		ImageURL:     product.ImageURL,
		IsActive:     product.IsActive,
	}

	if product.CategoryID != nil {
		response.CategoryID = product.CategoryID.String()
		if product.Category != nil {
			response.CategoryName = product.Category.Name
		}
	}

	if product.LastStockUpdate != nil {
		response.LastStockUpdate = product.LastStockUpdate.Format("2006-01-02T15:04:05Z07:00")
	}

	return response
}
