package dtos

import (
	"arabella-api/internal/app/models"
)

// CreateCategoryRequest represents the request payload for creating a category
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Type string `json:"type" validate:"required,oneof=INCOME EXPENSE"`
}

// UpdateCategoryRequest represents the request payload for updating a category
type UpdateCategoryRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=1,max=100"`
	IsActive *bool   `json:"is_active" validate:"omitempty"`
}

// CategoryResponse represents the full category response
type CategoryResponse struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

// CategorySummary represents a lightweight category for use in other DTOs
type CategorySummary struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// CategoryListResponse represents a list of categories
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int64              `json:"total"`
}

// ToModel converts CreateCategoryRequest to models.Category
func (r *CreateCategoryRequest) ToModel(userID uint) *models.Category {
	return &models.Category{
		UserID:   userID,
		Name:     r.Name,
		Type:     r.Type,
		IsActive: true, // Default to active
	}
}

// FromModelToCategoryResponse converts models.Category to CategoryResponse
func FromModelToCategoryResponse(c *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:       c.ID,
		UserID:   c.UserID,
		Name:     c.Name,
		Type:     c.Type,
		IsActive: c.IsActive,
	}
}

// FromModelToCategorySummary converts models.Category to CategorySummary
func FromModelToCategorySummary(c *models.Category) CategorySummary {
	return CategorySummary{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
	}
}
