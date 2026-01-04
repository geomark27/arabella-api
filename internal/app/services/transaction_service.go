package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"fmt"
)

// TransactionService handles transaction-related business logic
// It coordinates with the AccountingEngine to ensure proper double-entry bookkeeping
type TransactionService interface {
	Create(req *dtos.CreateTransactionRequest, userID uint) (*models.Transaction, error)
	GetByID(id uint) (*dtos.TransactionResponse, error)
	GetByUser(userID uint, filters dtos.TransactionFilters) (*dtos.TransactionListResponse, error)
	Update(id uint, req *dtos.UpdateTransactionRequest) error
	Delete(id uint) error
}

type transactionService struct {
	transactionRepo  repositories.TransactionRepository
	accountingEngine AccountingEngineService
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	accountingEngine AccountingEngineService,
) TransactionService {
	return &transactionService{
		transactionRepo:  transactionRepo,
		accountingEngine: accountingEngine,
	}
}

// Create creates a new transaction and processes it through the accounting engine
func (s *transactionService) Create(req *dtos.CreateTransactionRequest, userID uint) (*models.Transaction, error) {
	// Convert DTO to model
	transaction, err := req.ToModel(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction data: %w", err)
	}

	// Process through accounting engine (creates journal entries, updates balances)
	if err := s.accountingEngine.ProcessTransaction(transaction); err != nil {
		return nil, fmt.Errorf("failed to process transaction: %w", err)
	}

	// Reload transaction with relationships
	return s.transactionRepo.FindByID(transaction.ID)
}

// GetByID retrieves a transaction by ID with all relationships loaded
func (s *transactionService) GetByID(id uint) (*dtos.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := dtos.FromModelToTransactionResponse(transaction)
	return &response, nil
}

// GetByUser retrieves transactions for a user with filters and pagination
func (s *transactionService) GetByUser(userID uint, filters dtos.TransactionFilters) (*dtos.TransactionListResponse, error) {
	transactions, total, err := s.transactionRepo.FindByUser(userID, filters)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]dtos.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = dtos.FromModelToTransactionResponse(tx)
	}

	// Calculate pagination
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &dtos.TransactionListResponse{
		Transactions: responses,
		Total:        total,
		Page:         filters.Page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}

// Update updates a transaction
// Note: In double-entry bookkeeping, we typically don't "update" transactions
// Instead, we reverse the old one and create a new one
func (s *transactionService) Update(id uint, req *dtos.UpdateTransactionRequest) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Apply updates
	if req.Description != nil {
		transaction.Description = *req.Description
	}
	if req.Notes != nil {
		transaction.Notes = *req.Notes
	}
	if req.IsReconciled != nil {
		transaction.IsReconciled = *req.IsReconciled
	}

	// Note: We don't allow updating amount or accounts
	// If those need to change, user should delete and create a new transaction

	return s.transactionRepo.Update(transaction)
}

// Delete deletes a transaction by reversing it (proper accounting practice)
func (s *transactionService) Delete(id uint) error {
	// Use the accounting engine to reverse the transaction
	return s.accountingEngine.ReverseTransaction(id)
}
