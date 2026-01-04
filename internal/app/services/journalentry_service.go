package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/repositories"
)

// JournalEntryService handles journal entry queries (read-only)
// Journal entries are created automatically by the Accounting Engine
type JournalEntryService interface {
	GetByTransaction(transactionID uint) ([]dtos.JournalEntryResponse, error)
	GetByUser(userID uint, filters dtos.JournalEntryFilters) (*dtos.JournalEntryListResponse, error)
	VerifyTransactionBalance(transactionID uint) (*dtos.BalanceVerificationResponse, error)
}

type journalEntryService struct {
	journalEntryRepo repositories.JournalEntryRepository
}

// NewJournalEntryService creates a new journal entry service
func NewJournalEntryService(journalEntryRepo repositories.JournalEntryRepository) JournalEntryService {
	return &journalEntryService{
		journalEntryRepo: journalEntryRepo,
	}
}

// GetByTransaction retrieves all journal entries for a transaction
func (s *journalEntryService) GetByTransaction(transactionID uint) ([]dtos.JournalEntryResponse, error) {
	entries, err := s.journalEntryRepo.FindByTransaction(transactionID)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.JournalEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dtos.FromModelToJournalEntryResponse(entry)
	}

	return responses, nil
}

// GetByUser retrieves journal entries for a user with filters (audit trail)
func (s *journalEntryService) GetByUser(userID uint, filters dtos.JournalEntryFilters) (*dtos.JournalEntryListResponse, error) {
	entries, total, err := s.journalEntryRepo.FindByUser(userID, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.JournalEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = dtos.FromModelToJournalEntryResponse(entry)
	}

	// Calculate pagination
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 50
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dtos.JournalEntryListResponse{
		Entries:    responses,
		Total:      total,
		Page:       filters.Page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// VerifyTransactionBalance verifies that debits = credits for a transaction
func (s *journalEntryService) VerifyTransactionBalance(transactionID uint) (*dtos.BalanceVerificationResponse, error) {
	totalDebits, totalCredits, err := s.journalEntryRepo.VerifyTransactionBalance(transactionID)
	if err != nil {
		return nil, err
	}

	difference := totalDebits.Sub(totalCredits)
	isBalanced := difference.IsZero()

	message := "✅ Transaction is balanced (Debits = Credits)"
	if !isBalanced {
		message = "⚠️ WARNING: Transaction is NOT balanced!"
	}

	return &dtos.BalanceVerificationResponse{
		TransactionID: &transactionID,
		TotalDebits:   totalDebits,
		TotalCredits:  totalCredits,
		Difference:    difference,
		IsBalanced:    isBalanced,
		Message:       message,
	}, nil
}
