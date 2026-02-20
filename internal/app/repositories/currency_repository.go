package repositories

import (
	"arabella-api/internal/app/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// CurrencyRepository defines the interface for currency data access
type CurrencyRepository interface {
	Create(currency *models.Currency) error
	FindByID(id uint) (*models.Currency, error)
	FindByCode(code string) (*models.Currency, error)
	GetAll() ([]*models.Currency, error)
	GetActiveCurrencies() ([]*models.Currency, error)
	Update(currency *models.Currency) error
	Delete(id uint) error
}

// currencyRepositoryImpl implements CurrencyRepository using GORM
type currencyRepositoryImpl struct {
	db *gorm.DB
}

// NewCurrencyRepository creates a new currency repository
func NewCurrencyRepository(db *gorm.DB) CurrencyRepository {
	return &currencyRepositoryImpl{db: db}
}

// Create creates a new currency
func (r *currencyRepositoryImpl) Create(currency *models.Currency) error {
	return r.db.Create(currency).Error
}

// FindByID finds a currency by ID
func (r *currencyRepositoryImpl) FindByID(id uint) (*models.Currency, error) {
	var currency models.Currency

	err := r.db.First(&currency, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// FindByCode finds a currency by its ISO code (e.g., "USD", "usd", "Usd").
// The lookup is case-insensitive: the code is normalized to uppercase before querying.
func (r *currencyRepositoryImpl) FindByCode(code string) (*models.Currency, error) {
	var currency models.Currency

	err := r.db.Where("code = ?", strings.ToUpper(strings.TrimSpace(code))).First(&currency).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}

	return &currency, nil
}

// GetAll retrieves all currencies
func (r *currencyRepositoryImpl) GetAll() ([]*models.Currency, error) {
	var currencies []*models.Currency

	err := r.db.Order("code ASC").Find(&currencies).Error
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

// GetActiveCurrencies retrieves all active currencies
func (r *currencyRepositoryImpl) GetActiveCurrencies() ([]*models.Currency, error) {
	var currencies []*models.Currency

	err := r.db.
		Where("is_active = ?", true).
		Order("code ASC").
		Find(&currencies).Error

	if err != nil {
		return nil, err
	}

	return currencies, nil
}

// Update updates an existing currency
func (r *currencyRepositoryImpl) Update(currency *models.Currency) error {
	return r.db.Save(currency).Error
}

// Delete soft deletes a currency
func (r *currencyRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Currency{}, id).Error
}
