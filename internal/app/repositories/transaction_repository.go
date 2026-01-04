package repositories

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(tx *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByUser(userID uint, filters dtos.TransactionFilters) ([]*models.Transaction, int64, error)
	FindByAccount(accountID uint) ([]*models.Transaction, error)
	FindByDateRange(userID uint, startDate, endDate time.Time) ([]*models.Transaction, error)
	Update(tx *models.Transaction) error
	Delete(id uint) error
	GetMonthlyStats(userID uint, month, year int) (income, expenses decimal.Decimal, count int64, err error)
}

// transactionRepositoryImpl implements TransactionRepository using GORM
type transactionRepositoryImpl struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

// Create creates a new transaction
func (r *transactionRepositoryImpl) Create(tx *models.Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	return r.db.Create(tx).Error
}

// FindByID finds a transaction by ID with all relationships preloaded
func (r *transactionRepositoryImpl) FindByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction

	err := r.db.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Preload("Category").
		First(&tx, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &tx, nil
}

// FindByUser finds transactions for a user with filters and pagination
func (r *transactionRepositoryImpl) FindByUser(userID uint, filters dtos.TransactionFilters) ([]*models.Transaction, int64, error) {
	var transactions []*models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("user_id = ?", userID)

	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.AccountID != nil {
		query = query.Where("account_from_id = ? OR account_to_id = ?", *filters.AccountID, *filters.AccountID)
	}

	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}

	if filters.StartDate != nil {
		query = query.Where("transaction_date >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("transaction_date <= ?", *filters.EndDate)
	}

	if filters.IsReconciled != nil {
		query = query.Where("is_reconciled = ?", *filters.IsReconciled)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := filters.Page
	if page < 1 {
		page = 1
	}

	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get transactions with preloaded relationships
	err := query.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Preload("Category").
		Order("transaction_date DESC, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// FindByAccount finds all transactions for a specific account
func (r *transactionRepositoryImpl) FindByAccount(accountID uint) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	err := r.db.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Preload("Category").
		Where("account_from_id = ? OR account_to_id = ?", accountID, accountID).
		Order("transaction_date DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// FindByDateRange finds transactions within a date range
func (r *transactionRepositoryImpl) FindByDateRange(userID uint, startDate, endDate time.Time) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	err := r.db.
		Where("user_id = ? AND transaction_date >= ? AND transaction_date <= ?", userID, startDate, endDate).
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Preload("Category").
		Order("transaction_date ASC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Update updates an existing transaction
func (r *transactionRepositoryImpl) Update(tx *models.Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}

	return r.db.Save(tx).Error
}

// Delete soft deletes a transaction
func (r *transactionRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

// GetMonthlyStats calculates monthly income and expense statistics
func (r *transactionRepositoryImpl) GetMonthlyStats(userID uint, month, year int) (income, expenses decimal.Decimal, count int64, err error) {
	// Calculate start and end dates for the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Get total income
	var incomeTotal decimal.Decimal
	err = r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "INCOME", startDate, endDate).
		Scan(&incomeTotal).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, 0, err
	}

	// Get total expenses
	var expensesTotal decimal.Decimal
	err = r.db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, "EXPENSE", startDate, endDate).
		Scan(&expensesTotal).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, 0, err
	}

	// Get transaction count
	err = r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND transaction_date >= ? AND transaction_date <= ?",
			userID, startDate, endDate).
		Count(&count).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, 0, err
	}

	return incomeTotal, expensesTotal, count, nil
}
