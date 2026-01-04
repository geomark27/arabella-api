package repositories

import (
	"arabella-api/internal/app/models"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccountRepository defines the interface for account data access
type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id uint) (*models.Account, error)
	FindByUserID(userID uint) ([]*models.Account, error)
	FindByUserAndType(userID uint, accountType string) ([]*models.Account, error)
	Update(account *models.Account) error
	Delete(id uint) error
	UpdateBalance(accountID uint, amount decimal.Decimal) error
	GetTotalAssets(userID uint) (decimal.Decimal, error)
	GetTotalLiabilities(userID uint) (decimal.Decimal, error)
	GetLiquidAssets(userID uint) (decimal.Decimal, error)
}

// accountRepositoryImpl implements AccountRepository using GORM
type accountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepositoryImpl{db: db}
}

// Create creates a new account
func (r *accountRepositoryImpl) Create(account *models.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	return r.db.Create(account).Error
}

// FindByID finds an account by ID with currency preloaded
func (r *accountRepositoryImpl) FindByID(id uint) (*models.Account, error) {
	var account models.Account

	err := r.db.Preload("Currency").First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	return &account, nil
}

// FindByUserID finds all accounts for a user
func (r *accountRepositoryImpl) FindByUserID(userID uint) ([]*models.Account, error) {
	var accounts []*models.Account

	err := r.db.
		Preload("Currency").
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindByUserAndType finds accounts by user ID and account type
func (r *accountRepositoryImpl) FindByUserAndType(userID uint, accountType string) ([]*models.Account, error) {
	var accounts []*models.Account

	err := r.db.
		Preload("Currency").
		Where("user_id = ? AND account_type = ? AND is_active = ?", userID, accountType, true).
		Order("created_at DESC").
		Find(&accounts).Error

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// Update updates an existing account
func (r *accountRepositoryImpl) Update(account *models.Account) error {
	if err := account.Validate(); err != nil {
		return err
	}

	return r.db.Save(account).Error
}

// Delete soft deletes an account
func (r *accountRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}

// UpdateBalance updates the account balance atomically using optimistic locking
func (r *accountRepositoryImpl) UpdateBalance(accountID uint, amount decimal.Decimal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var account models.Account

		// Find the account within the transaction
		if err := tx.Where("id = ?", accountID).First(&account).Error; err != nil {
			return err
		}

		// Update balance
		account.UpdateBalance(amount)

		// Save the updated balance
		return tx.Save(&account).Error
	})
}

// GetTotalAssets calculates total assets (BANK + CASH + SAVINGS + INVESTMENT)
func (r *accountRepositoryImpl) GetTotalAssets(userID uint) (decimal.Decimal, error) {
	var total decimal.Decimal

	err := r.db.Model(&models.Account{}).
		Select("COALESCE(SUM(balance), 0)").
		Where("user_id = ? AND account_type IN (?, ?, ?, ?) AND is_active = ?",
			userID, "BANK", "CASH", "SAVINGS", "INVESTMENT", true).
		Scan(&total).Error

	if err != nil {
		return decimal.Zero, err
	}

	return total, nil
}

// GetTotalLiabilities calculates total liabilities (CREDIT_CARD)
func (r *accountRepositoryImpl) GetTotalLiabilities(userID uint) (decimal.Decimal, error) {
	var total decimal.Decimal

	err := r.db.Model(&models.Account{}).
		Select("COALESCE(SUM(balance), 0)").
		Where("user_id = ? AND account_type = ? AND is_active = ?",
			userID, "CREDIT_CARD", true).
		Scan(&total).Error

	if err != nil {
		return decimal.Zero, err
	}

	// Credit card balances are positive (we owe that amount), so return as is
	return total, nil
}

// GetLiquidAssets calculates liquid assets (BANK + CASH only)
func (r *accountRepositoryImpl) GetLiquidAssets(userID uint) (decimal.Decimal, error) {
	var total decimal.Decimal

	err := r.db.Model(&models.Account{}).
		Select("COALESCE(SUM(balance), 0)").
		Where("user_id = ? AND account_type IN (?, ?) AND is_active = ?",
			userID, "BANK", "CASH", true).
		Scan(&total).Error

	if err != nil {
		return decimal.Zero, err
	}

	return total, nil
}
