package repositories

import (
	"arabella-api/internal/app/models"
	"errors"

	"gorm.io/gorm"
)

// SystemValueRepository defines the interface for system value data access
type SystemValueRepository interface {
	Create(sv *models.SystemValue) error
	FindByID(id uint) (*models.SystemValue, error)
	FindByCatalogType(catalogType string) ([]*models.SystemValue, error)
	FindByCatalogTypeAndValue(catalogType, value string) (*models.SystemValue, error)
	GetAccountTypes() ([]*models.SystemValue, error)
	GetTransactionTypes() ([]*models.SystemValue, error)
	Update(sv *models.SystemValue) error
	Delete(id uint) error
}

// systemValueRepositoryImpl implements SystemValueRepository using GORM
type systemValueRepositoryImpl struct {
	db *gorm.DB
}

// NewSystemValueRepository creates a new system value repository
func NewSystemValueRepository(db *gorm.DB) SystemValueRepository {
	return &systemValueRepositoryImpl{db: db}
}

// Create creates a new system value
func (r *systemValueRepositoryImpl) Create(sv *models.SystemValue) error {
	return r.db.Create(sv).Error
}

// FindByID finds a system value by ID
func (r *systemValueRepositoryImpl) FindByID(id uint) (*models.SystemValue, error) {
	var sv models.SystemValue

	err := r.db.First(&sv, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("system value not found")
		}
		return nil, err
	}

	return &sv, nil
}

// FindByCatalogType finds all system values for a catalog type
func (r *systemValueRepositoryImpl) FindByCatalogType(catalogType string) ([]*models.SystemValue, error) {
	var values []*models.SystemValue

	err := r.db.
		Where("catalog_type = ? AND is_active = ?", catalogType, true).
		Order("display_order ASC, label ASC").
		Find(&values).Error

	if err != nil {
		return nil, err
	}

	return values, nil
}

// FindByCatalogTypeAndValue finds a specific system value by catalog type and value
func (r *systemValueRepositoryImpl) FindByCatalogTypeAndValue(catalogType, value string) (*models.SystemValue, error) {
	var sv models.SystemValue

	err := r.db.
		Where("catalog_type = ? AND value = ?", catalogType, value).
		First(&sv).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("system value not found")
		}
		return nil, err
	}

	return &sv, nil
}

// GetAccountTypes retrieves all account type system values
func (r *systemValueRepositoryImpl) GetAccountTypes() ([]*models.SystemValue, error) {
	return r.FindByCatalogType("ACCOUNT_TYPE")
}

// GetTransactionTypes retrieves all transaction type system values
func (r *systemValueRepositoryImpl) GetTransactionTypes() ([]*models.SystemValue, error) {
	return r.FindByCatalogType("TRANSACTION_TYPE")
}

// Update updates an existing system value
func (r *systemValueRepositoryImpl) Update(sv *models.SystemValue) error {
	return r.db.Save(sv).Error
}

// Delete soft deletes a system value
func (r *systemValueRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.SystemValue{}, id).Error
}
