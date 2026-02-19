package dtos

import (
	"arabella-api/internal/app/models"
	"time"

	"github.com/shopspring/decimal"
)

type AccountFiltreDTO struct {
	Name        string `form:"name" json:"name"`
	AccountType string `form:"account_type" json:"account_type"`
	CurrencyID  uint   `form:"currency_id" json:"currency_id"`
	IsActive    *bool  `form:"is_active" json:"is_active"`
}

type AccountResponseDTO struct {
	ID          uint             `json:"id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Name        string           `json:"name"`
	AccountType string           `json:"account_type"`
	CurrencyID  *uint            `json:"currency_id"`
	Currency    *CurrencySummary `json:"currency,omitempty"`
	Balance     decimal.Decimal  `json:"balance"` // Changed to decimal.Decimal to match model
	IsActive    bool             `json:"is_active"`
}

type CreateAccountDTO struct {
	UserID      uint            `json:"-"` // Set from JWT context, not from request body
	Name        string          `json:"name" binding:"required"`
	AccountType string          `json:"account_type" binding:"required"`
	CurrencyID  *uint           `json:"currency_id" binding:"required"`
	Balance     decimal.Decimal `json:"balance"`
	IsActive    *bool           `json:"is_active"`
}

type UpdateAccountDTO struct {
	Name        *string          `json:"name"`
	AccountType *string          `json:"account_type"`
	CurrencyID  *uint            `json:"currency_id"`
	Balance     *decimal.Decimal `json:"balance"` // Changed to *decimal.Decimal for optional updates
	IsActive    *bool            `json:"is_active"`
}

func ToAccountResponse(account *models.Account) *AccountResponseDTO {
	dto := &AccountResponseDTO{
		ID:          account.ID,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
		Name:        account.Name,
		AccountType: account.AccountType,
		CurrencyID:  account.CurrencyID,
		Balance:     account.Balance,
		IsActive:    account.IsActive,
	}
	if account.Currency != nil {
		dto.Currency = &CurrencySummary{
			ID:     account.Currency.ID,
			Code:   account.Currency.Code,
			Symbol: account.Currency.Symbol,
		}
	}
	return dto
}

func ToAccountResponseList(accounts []models.Account) []AccountResponseDTO {
	result := make([]AccountResponseDTO, len(accounts))
	for i, account := range accounts {
		result[i] = *ToAccountResponse(&account)
	}
	return result
}

// AccountSummary represents a lightweight account for use in other DTOs
type AccountSummary struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Balance  interface{}      `json:"balance"` // Can be string or decimal.Decimal depending on model version
	Currency *CurrencySummary `json:"currency,omitempty"`
}
