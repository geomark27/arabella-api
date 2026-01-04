package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"fmt"

	"github.com/shopspring/decimal"
)

// AccountService handles account-related business logic
type AccountService interface {
	Create(req *dtos.CreateAccountDTO) (*models.Account, error)
	GetByID(id uint) (*dtos.AccountResponseDTO, error)
	GetByUserID(userID uint) ([]dtos.AccountResponseDTO, error)
	Update(id uint, req *dtos.UpdateAccountDTO) error
	Delete(id uint) error
}

type accountService struct {
	accountRepo     repositories.AccountRepository
	systemValueRepo repositories.SystemValueRepository
}

// NewAccountService creates a new account service
func NewAccountService(
	accountRepo repositories.AccountRepository,
	systemValueRepo repositories.SystemValueRepository,
) AccountService {
	return &accountService{
		accountRepo:     accountRepo,
		systemValueRepo: systemValueRepo,
	}
}

// Create creates a new account
func (s *accountService) Create(req *dtos.CreateAccountDTO) (*models.Account, error) {
	// Validate account type against SystemValue
	if err := s.validateAccountType(req.AccountType); err != nil {
		return nil, err
	}

	// Set default balance if not provided
	balance := req.Balance
	if balance.IsZero() {
		balance = decimal.Zero
	}

	account := &models.Account{
		UserID:      req.UserID,
		Name:        req.Name,
		AccountType: req.AccountType,
		CurrencyID:  req.CurrencyID,
		Balance:     balance,
		IsActive:    true,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	// Reload with currency relationship
	return s.accountRepo.FindByID(account.ID)
}

// validateAccountType validates the account type against system values
func (s *accountService) validateAccountType(accountType string) error {
	_, err := s.systemValueRepo.FindByCatalogTypeAndValue("ACCOUNT_TYPE", accountType)
	if err != nil {
		return fmt.Errorf("invalid account type: %s", accountType)
	}
	return nil
}

// GetByID retrieves an account by ID
func (s *accountService) GetByID(id uint) (*dtos.AccountResponseDTO, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := dtos.ToAccountResponse(account)
	return response, nil
}

// GetByUserID retrieves all accounts for a user
func (s *accountService) GetByUserID(userID uint) ([]dtos.AccountResponseDTO, error) {
	accounts, err := s.accountRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.AccountResponseDTO, len(accounts))
	for i, account := range accounts {
		responses[i] = *dtos.ToAccountResponse(account)
	}

	return responses, nil
}

// Update updates an account
func (s *accountService) Update(id uint, req *dtos.UpdateAccountDTO) error {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Apply updates
	if req.Name != nil {
		account.Name = *req.Name
	}
	if req.AccountType != nil {
		account.AccountType = *req.AccountType
	}
	if req.CurrencyID != nil {
		account.CurrencyID = req.CurrencyID
	}
	if req.Balance != nil {
		account.Balance = *req.Balance
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}

	return s.accountRepo.Update(account)
}

// Delete soft deletes an account
func (s *accountService) Delete(id uint) error {
	// TODO: Check if account has transactions before deleting
	// For now, just mark as inactive instead of hard delete
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return err
	}

	account.IsActive = false
	return s.accountRepo.Update(account)
}
