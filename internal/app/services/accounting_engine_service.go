package services

import (
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// AccountingEngineService is the CORE of the entire system
// It implements double-entry bookkeeping logic
// Every transaction MUST go through this engine to maintain data integrity
type AccountingEngineService interface {
	ProcessTransaction(tx *models.Transaction) error
	ReverseTransaction(transactionID uint) error
	VerifyTransactionBalance(transactionID uint) (bool, error)
}

type accountingEngineService struct {
	db                     *gorm.DB
	journalEntryRepository repositories.JournalEntryRepository
	accountRepository      repositories.AccountRepository
	transactionRepository  repositories.TransactionRepository
}

// NewAccountingEngineService creates a new accounting engine service
func NewAccountingEngineService(
	db *gorm.DB,
	journalEntryRepo repositories.JournalEntryRepository,
	accountRepo repositories.AccountRepository,
	transactionRepo repositories.TransactionRepository,
) AccountingEngineService {
	return &accountingEngineService{
		db:                     db,
		journalEntryRepository: journalEntryRepo,
		accountRepository:      accountRepo,
		transactionRepository:  transactionRepo,
	}
}

// ProcessTransaction is the main entry point for creating a transaction
// It performs the following steps:
// 1. Validates the transaction
// 2. Generates journal entries (Debit + Credit)
// 3. Verifies balance (SUM(Debits) = SUM(Credits))
// 4. Saves everything in a database transaction (atomic)
// 5. Updates account balances
func (s *accountingEngineService) ProcessTransaction(tx *models.Transaction) error {
	// Step 1: Validate the transaction
	if err := tx.Validate(); err != nil {
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	// Step 2: Start a database transaction (everything or nothing)
	return s.db.Transaction(func(dbTx *gorm.DB) error {
		// Step 2a: Save the transaction first to get its ID
		if err := dbTx.Create(tx).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Step 3: Generate journal entries based on transaction type
		entries, err := s.generateJournalEntries(tx)
		if err != nil {
			return fmt.Errorf("failed to generate journal entries: %w", err)
		}

		// Step 4: Validate that debits = credits
		if err := s.validateBalance(entries); err != nil {
			return fmt.Errorf("balance validation failed: %w", err)
		}

		// Step 5: Save journal entries
		if err := dbTx.Create(&entries).Error; err != nil {
			return fmt.Errorf("failed to save journal entries: %w", err)
		}

		// Step 6: Update account balances
		if err := s.updateAccountBalances(dbTx, entries); err != nil {
			return fmt.Errorf("failed to update account balances: %w", err)
		}

		return nil
	})
}

// generateJournalEntries creates the debit and credit entries based on transaction type
func (s *accountingEngineService) generateJournalEntries(tx *models.Transaction) ([]*models.JournalEntry, error) {
	entries := []*models.JournalEntry{}

	switch tx.Type {
	case "EXPENSE":
		// EXPENSE: Money leaves an account and goes to an expense category
		// Debit: Expense Category (increases expense)
		// Credit: Account (decreases asset)

		if tx.CategoryID == nil {
			return nil, errors.New("category_id is required for EXPENSE transactions")
		}

		// Expense category receives a DEBIT (expense increases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     *tx.CategoryID, // Using CategoryID as a virtual account
			DebitOrCredit: "DEBIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Expense: %s", tx.Description),
		})

		// Bank/Cash account receives a CREDIT (asset decreases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     tx.AccountFromID,
			DebitOrCredit: "CREDIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Payment: %s", tx.Description),
		})

	case "INCOME":
		// INCOME: Money enters an account from an income category
		// Debit: Account (increases asset)
		// Credit: Income Category (increases income)

		if tx.CategoryID == nil {
			return nil, errors.New("category_id is required for INCOME transactions")
		}

		// Bank/Cash account receives a DEBIT (asset increases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     tx.AccountFromID,
			DebitOrCredit: "DEBIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Income: %s", tx.Description),
		})

		// Income category receives a CREDIT (income increases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     *tx.CategoryID,
			DebitOrCredit: "CREDIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Revenue: %s", tx.Description),
		})

	case "TRANSFER":
		// TRANSFER: Money moves from one account to another
		// Debit: Account To (destination asset increases)
		// Credit: Account From (source asset decreases)

		if tx.AccountToID == nil {
			return nil, errors.New("account_to_id is required for TRANSFER transactions")
		}

		// Destination account receives a DEBIT (asset increases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     *tx.AccountToID,
			DebitOrCredit: "DEBIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Transfer in: %s", tx.Description),
		})

		// Source account receives a CREDIT (asset decreases)
		entries = append(entries, &models.JournalEntry{
			UserID:        tx.UserID,
			TransactionID: tx.ID,
			AccountID:     tx.AccountFromID,
			DebitOrCredit: "CREDIT",
			Amount:        tx.Amount,
			EntryDate:     tx.TransactionDate,
			Description:   fmt.Sprintf("Transfer out: %s", tx.Description),
		})

	default:
		return nil, fmt.Errorf("unsupported transaction type: %s", tx.Type)
	}

	return entries, nil
}

// validateBalance ensures that SUM(Debits) = SUM(Credits)
// This is the FUNDAMENTAL rule of double-entry bookkeeping
func (s *accountingEngineService) validateBalance(entries []*models.JournalEntry) error {
	var totalDebits decimal.Decimal
	var totalCredits decimal.Decimal

	for _, entry := range entries {
		if entry.DebitOrCredit == "DEBIT" {
			totalDebits = totalDebits.Add(entry.Amount)
		} else if entry.DebitOrCredit == "CREDIT" {
			totalCredits = totalCredits.Add(entry.Amount)
		}
	}

	if !totalDebits.Equal(totalCredits) {
		return fmt.Errorf(
			"debits and credits must be equal: debits=%s, credits=%s, difference=%s",
			totalDebits.String(),
			totalCredits.String(),
			totalDebits.Sub(totalCredits).String(),
		)
	}

	return nil
}

// updateAccountBalances updates the balances of affected accounts
// For EXPENSE/INCOME: Only update the real account (not categories)
// For TRANSFER: Update both accounts
func (s *accountingEngineService) updateAccountBalances(dbTx *gorm.DB, entries []*models.JournalEntry) error {
	accountUpdates := make(map[uint]decimal.Decimal)

	for _, entry := range entries {
		// Skip category entries (we only update real accounts)
		// In a more complete implementation, categories would be in a separate table
		// For now, we assume AccountID > 1000 means it's a category (convention)
		// TODO: Improve this logic when category accounts are properly separated

		amount := entry.Amount
		if entry.DebitOrCredit == "CREDIT" {
			amount = amount.Neg() // Credit decreases account balance
		}

		if existing, ok := accountUpdates[entry.AccountID]; ok {
			accountUpdates[entry.AccountID] = existing.Add(amount)
		} else {
			accountUpdates[entry.AccountID] = amount
		}
	}

	// Apply updates to each account
	for accountID, balanceChange := range accountUpdates {
		var account models.Account
		if err := dbTx.First(&account, accountID).Error; err != nil {
			// If account not found, it might be a category - skip for now
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return fmt.Errorf("failed to find account %d: %w", accountID, err)
		}

		account.UpdateBalance(balanceChange)

		if err := dbTx.Save(&account).Error; err != nil {
			return fmt.Errorf("failed to update account %d balance: %w", accountID, err)
		}
	}

	return nil
}

// ReverseTransaction creates reversing entries to cancel out a transaction
// This is used when a transaction needs to be "deleted" (we never truly delete in accounting)
func (s *accountingEngineService) ReverseTransaction(transactionID uint) error {
	return s.db.Transaction(func(dbTx *gorm.DB) error {
		// Get original transaction
		tx, err := s.transactionRepository.FindByID(transactionID)
		if err != nil {
			return fmt.Errorf("transaction not found: %w", err)
		}

		// Get original journal entries
		originalEntries, err := s.journalEntryRepository.FindByTransaction(transactionID)
		if err != nil {
			return fmt.Errorf("failed to find journal entries: %w", err)
		}

		// Create reversing entries (swap DEBIT <-> CREDIT)
		reversingEntries := []*models.JournalEntry{}
		now := time.Now()

		for _, original := range originalEntries {
			reversedType := "DEBIT"
			if original.DebitOrCredit == "DEBIT" {
				reversedType = "CREDIT"
			}

			reversingEntries = append(reversingEntries, &models.JournalEntry{
				UserID:        original.UserID,
				TransactionID: transactionID,
				AccountID:     original.AccountID,
				DebitOrCredit: reversedType,
				Amount:        original.Amount,
				EntryDate:     now,
				Description:   fmt.Sprintf("REVERSAL: %s", original.Description),
			})
		}

		// Save reversing entries
		if err := dbTx.Create(&reversingEntries).Error; err != nil {
			return fmt.Errorf("failed to create reversing entries: %w", err)
		}

		// Update account balances (reverse the changes)
		if err := s.updateAccountBalances(dbTx, reversingEntries); err != nil {
			return fmt.Errorf("failed to update balances during reversal: %w", err)
		}

		// Mark original transaction as reconciled (archived)
		tx.IsReconciled = true
		if err := dbTx.Save(tx).Error; err != nil {
			return fmt.Errorf("failed to mark transaction as reversed: %w", err)
		}

		return nil
	})
}

// VerifyTransactionBalance verifies that a transaction's journal entries balance
func (s *accountingEngineService) VerifyTransactionBalance(transactionID uint) (bool, error) {
	totalDebit, totalCredit, err := s.journalEntryRepository.VerifyTransactionBalance(transactionID)
	if err != nil {
		return false, err
	}

	return totalDebit.Equal(totalCredit), nil
}
