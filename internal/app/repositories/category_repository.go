package repositories

import (
	"arabella-api/internal/app/models"
	"errors"

	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(category *models.Category) error
	FindByID(id uint) (*models.Category, error)
	FindByUserID(userID uint) ([]*models.Category, error)
	FindByUserAndType(userID uint, categoryType string) ([]*models.Category, error)
	Update(category *models.Category) error
	Delete(id uint) error
}

// categoryRepositoryImpl implements CategoryRepository using GORM
type categoryRepositoryImpl struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db: db}
}

// Create creates a new category
func (r *categoryRepositoryImpl) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID
func (r *categoryRepositoryImpl) FindByID(id uint) (*models.Category, error) {
	var category models.Category

	err := r.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return &category, nil
}

// FindByUserID finds all active categories for a user
func (r *categoryRepositoryImpl) FindByUserID(userID uint) ([]*models.Category, error) {
	var categories []*models.Category

	err := r.db.
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("type ASC, name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// FindByUserAndType finds categories by user ID and type (INCOME or EXPENSE)
func (r *categoryRepositoryImpl) FindByUserAndType(userID uint, categoryType string) ([]*models.Category, error) {
	var categories []*models.Category

	err := r.db.
		Where("user_id = ? AND type = ? AND is_active = ?", userID, categoryType, true).
		Order("name ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Update updates an existing category
func (r *categoryRepositoryImpl) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes a category
func (r *categoryRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
