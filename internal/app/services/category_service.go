package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
)

// CategoryService handles category-related business logic
type CategoryService interface {
	Create(req *dtos.CreateCategoryRequest, userID uint) (*models.Category, error)
	GetByID(id uint) (*dtos.CategoryResponse, error)
	GetByUser(userID uint) ([]dtos.CategoryResponse, error)
	GetByType(userID uint, categoryType string) ([]dtos.CategoryResponse, error)
	Update(id uint, req *dtos.UpdateCategoryRequest) error
	Delete(id uint) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// Create creates a new category
func (s *categoryService) Create(req *dtos.CreateCategoryRequest, userID uint) (*models.Category, error) {
	category := req.ToModel(userID)

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetByID retrieves a category by ID
func (s *categoryService) GetByID(id uint) (*dtos.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := dtos.FromModelToCategoryResponse(category)
	return &response, nil
}

// GetByUser retrieves all categories for a user
func (s *categoryService) GetByUser(userID uint) ([]dtos.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dtos.FromModelToCategoryResponse(cat)
	}

	return responses, nil
}

// GetByType retrieves categories by type (INCOME or EXPENSE)
func (s *categoryService) GetByType(userID uint, categoryType string) ([]dtos.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindByUserAndType(userID, categoryType)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dtos.FromModelToCategoryResponse(cat)
	}

	return responses, nil
}

// Update updates a category
func (s *categoryService) Update(id uint, req *dtos.UpdateCategoryRequest) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	return s.categoryRepo.Update(category)
}

// Delete soft deletes a category
func (s *categoryService) Delete(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	category.IsActive = false
	return s.categoryRepo.Update(category)
}
