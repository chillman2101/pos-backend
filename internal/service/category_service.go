package service

import (
	"errors"
	"pos-backend/internal/domain"
	"pos-backend/internal/dto"

	"github.com/google/uuid"
)

type CategoryService interface {
	GetAll(page, limit int) ([]*dto.CategoryResponse, int64, error)
	Create(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(id string) (*dto.CategoryResponse, error)
	GetByName(name string) (*dto.CategoryResponse, error)
	Update(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(id string) error
}

type categoryService struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryService(categoryRepo domain.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) Create(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.categoryRepo.Create(&category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) GetByID(id string) (*dto.CategoryResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	category, err := s.categoryRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) GetByName(name string) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) Update(id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	category, err := s.categoryRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (s *categoryService) Delete(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	return s.categoryRepo.Delete(userID)
}

func (s *categoryService) GetAll(page, limit int) ([]*dto.CategoryResponse, int64, error) {
	categories, totalData, err := s.categoryRepo.FindAll(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.CategoryResponse
	for _, category := range categories {
		responses = append(responses, &dto.CategoryResponse{
			ID:          category.ID.String(),
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return responses, totalData, nil
}
