package repositories

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// JournalEntryRepository defines the interface for journal entry data access
type JournalEntryRepository interface {
	CreateBatch(entries []*models.JournalEntry) error
	FindByID(id uint) (*models.JournalEntry, error)
	FindByTransaction(transactionID uint) ([]*models.JournalEntry, error)
	FindByAccount(accountID uint, filters dtos.JournalEntryFilters) ([]*models.JournalEntry, int64, error)
	FindByUser(userID uint, filters dtos.JournalEntryFilters) ([]*models.JournalEntry, int64, error)
	VerifyTransactionBalance(transactionID uint) (totalDebit, totalCredit decimal.Decimal, err error)
	GetAccountBalance(accountID uint, asOf *time.Time) (decimal.Decimal, error)
	GetBalanceSheet(userID uint, asOf time.Time) (map[uint]decimal.Decimal, error)
}

// journalEntryRepositoryImpl implements JournalEntryRepository using GORM
type journalEntryRepositoryImpl struct {
	db *gorm.DB
}

// NewJournalEntryRepository creates a new journal entry repository
func NewJournalEntryRepository(db *gorm.DB) JournalEntryRepository {
	return &journalEntryRepositoryImpl{db: db}
}

// CreateBatch creates multiple journal entries in a single transaction
// This is critical for double-entry bookkeeping - entries must be created atomically
func (r *journalEntryRepositoryImpl) CreateBatch(entries []*models.JournalEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Validate all entries first
		for _, entry := range entries {
			if err := entry.Validate(); err != nil {
				return err
			}
		}

		// Create all entries
		return tx.Create(&entries).Error
	})
}

// FindByID finds a journal entry by ID
func (r *journalEntryRepositoryImpl) FindByID(id uint) (*models.JournalEntry, error) {
	var entry models.JournalEntry

	err := r.db.
		Preload("Account.Currency").
		Preload("Transaction").
		First(&entry, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("journal entry not found")
		}
		return nil, err
	}

	return &entry, nil
}

// FindByTransaction finds all journal entries for a transaction
func (r *journalEntryRepositoryImpl) FindByTransaction(transactionID uint) ([]*models.JournalEntry, error) {
	var entries []*models.JournalEntry

	err := r.db.
		Preload("Account.Currency").
		Where("transaction_id = ?", transactionID).
		Order("debit_or_credit DESC"). // DEBIT first, then CREDIT
		Find(&entries).Error

	if err != nil {
		return nil, err
	}

	return entries, nil
}

// FindByAccount finds journal entries for a specific account with filters
func (r *journalEntryRepositoryImpl) FindByAccount(accountID uint, filters dtos.JournalEntryFilters) ([]*models.JournalEntry, int64, error) {
	var entries []*models.JournalEntry
	var total int64

	query := r.db.Model(&models.JournalEntry{}).Where("account_id = ?", accountID)

	// Apply filters
	if filters.DebitOrCredit != nil {
		query = query.Where("debit_or_credit = ?", *filters.DebitOrCredit)
	}

	if filters.StartDate != nil {
		query = query.Where("entry_date >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("entry_date <= ?", *filters.EndDate)
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
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get entries
	err := query.
		Preload("Account.Currency").
		Preload("Transaction").
		Order("entry_date DESC, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// FindByUser finds journal entries for a user with filters (audit trail)
func (r *journalEntryRepositoryImpl) FindByUser(userID uint, filters dtos.JournalEntryFilters) ([]*models.JournalEntry, int64, error) {
	var entries []*models.JournalEntry
	var total int64

	query := r.db.Model(&models.JournalEntry{}).Where("user_id = ?", userID)

	// Apply filters
	if filters.TransactionID != nil {
		query = query.Where("transaction_id = ?", *filters.TransactionID)
	}

	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}

	if filters.DebitOrCredit != nil {
		query = query.Where("debit_or_credit = ?", *filters.DebitOrCredit)
	}

	if filters.StartDate != nil {
		query = query.Where("entry_date >= ?", *filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("entry_date <= ?", *filters.EndDate)
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
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get entries
	err := query.
		Preload("Account.Currency").
		Preload("Transaction").
		Order("entry_date DESC, created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// VerifyTransactionBalance verifies that debits equal credits for a transaction
func (r *journalEntryRepositoryImpl) VerifyTransactionBalance(transactionID uint) (totalDebit, totalCredit decimal.Decimal, err error) {
	// Sum all debits
	err = r.db.Model(&models.JournalEntry{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("transaction_id = ? AND debit_or_credit = ?", transactionID, "DEBIT").
		Scan(&totalDebit).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	// Sum all credits
	err = r.db.Model(&models.JournalEntry{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("transaction_id = ? AND debit_or_credit = ?", transactionID, "CREDIT").
		Scan(&totalCredit).Error

	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	return totalDebit, totalCredit, nil
}

// GetAccountBalance calculates account balance from journal entries up to a point in time
func (r *journalEntryRepositoryImpl) GetAccountBalance(accountID uint, asOf *time.Time) (decimal.Decimal, error) {
	query := r.db.Model(&models.JournalEntry{}).Where("account_id = ?", accountID)

	if asOf != nil {
		query = query.Where("entry_date <= ?", *asOf)
	}

	// Get total debits
	var totalDebits decimal.Decimal
	err := query.
		Where("debit_or_credit = ?", "DEBIT").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalDebits).Error

	if err != nil {
		return decimal.Zero, err
	}

	// Get total credits
	var totalCredits decimal.Decimal
	err = query.
		Where("debit_or_credit = ?", "CREDIT").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalCredits).Error

	if err != nil {
		return decimal.Zero, err
	}

	// For asset accounts, balance = Debits - Credits
	// For liability/equity accounts, balance = Credits - Debits
	// We'll return Debits - Credits and let the service layer handle the sign
	balance := totalDebits.Sub(totalCredits)

	return balance, nil
}

// GetBalanceSheet returns balances for all accounts at a point in time
func (r *journalEntryRepositoryImpl) GetBalanceSheet(userID uint, asOf time.Time) (map[uint]decimal.Decimal, error) {
	type AccountBalance struct {
		AccountID    uint
		TotalDebits  decimal.Decimal
		TotalCredits decimal.Decimal
	}

	var results []AccountBalance

	err := r.db.Raw(`
		SELECT 
			account_id,
			COALESCE(SUM(CASE WHEN debit_or_credit = 'DEBIT' THEN amount ELSE 0 END), 0) as total_debits,
			COALESCE(SUM(CASE WHEN debit_or_credit = 'CREDIT' THEN amount ELSE 0 END), 0) as total_credits
		FROM journal_entries
		WHERE user_id = ? AND entry_date <= ?
		GROUP BY account_id
	`, userID, asOf).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to map
	balances := make(map[uint]decimal.Decimal)
	for _, result := range results {
		balances[result.AccountID] = result.TotalDebits.Sub(result.TotalCredits)
	}

	return balances, nil
}
