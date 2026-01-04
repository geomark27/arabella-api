package dtos

import (
	"arabella-api/internal/app/models"
	"time"

	"github.com/shopspring/decimal"
)

// CreateTransactionRequest represents the request payload for creating a transaction
type CreateTransactionRequest struct {
	Type            string          `json:"type" validate:"required,oneof=INCOME EXPENSE TRANSFER"`
	Description     string          `json:"description" validate:"required,min=1,max=255"`
	Amount          decimal.Decimal `json:"amount" validate:"required"`
	AccountFromID   uint            `json:"account_from_id" validate:"required,gt=0"`
	AccountToID     *uint           `json:"account_to_id" validate:"omitempty,gt=0"`
	CategoryID      *uint           `json:"category_id" validate:"omitempty,gt=0"`
	TransactionDate string          `json:"transaction_date" validate:"required"` // ISO 8601 format
	Notes           string          `json:"notes" validate:"omitempty,max=1000"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate" validate:"omitempty,gt=0"`
}

// UpdateTransactionRequest represents the request payload for updating a transaction
type UpdateTransactionRequest struct {
	Description     *string          `json:"description" validate:"omitempty,min=1,max=255"`
	Amount          *decimal.Decimal `json:"amount" validate:"omitempty,gt=0"`
	TransactionDate *string          `json:"transaction_date" validate:"omitempty"` // ISO 8601 format
	Notes           *string          `json:"notes" validate:"omitempty,max=1000"`
	IsReconciled    *bool            `json:"is_reconciled" validate:"omitempty"`
}

// TransactionResponse represents the full transaction response with all relationships
type TransactionResponse struct {
	ID              uint            `json:"id"`
	UserID          uint            `json:"user_id"`
	Type            string          `json:"type"`
	Description     string          `json:"description"`
	Amount          decimal.Decimal `json:"amount"`
	AmountInUSD     decimal.Decimal `json:"amount_in_usd"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate"`
	TransactionDate time.Time       `json:"transaction_date"`
	Notes           string          `json:"notes"`
	IsReconciled    bool            `json:"is_reconciled"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Relationships
	AccountFrom AccountSummary   `json:"account_from"`
	AccountTo   *AccountSummary  `json:"account_to,omitempty"`
	Category    *CategorySummary `json:"category,omitempty"`
}

// TransactionSummary represents a lightweight transaction for list views
type TransactionSummary struct {
	ID              uint            `json:"id"`
	Type            string          `json:"type"`
	Description     string          `json:"description"`
	Amount          decimal.Decimal `json:"amount"`
	TransactionDate time.Time       `json:"transaction_date"`
	AccountFromName string          `json:"account_from_name"`
	CategoryName    *string         `json:"category_name,omitempty"`
}

// TransactionListResponse represents a paginated list of transactions
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
}

// TransactionFilters represents query parameters for filtering transactions
type TransactionFilters struct {
	Type         string     `json:"type" validate:"omitempty,oneof=INCOME EXPENSE TRANSFER"`
	AccountID    *uint      `json:"account_id" validate:"omitempty,gt=0"`
	CategoryID   *uint      `json:"category_id" validate:"omitempty,gt=0"`
	StartDate    *time.Time `json:"start_date" validate:"omitempty"`
	EndDate      *time.Time `json:"end_date" validate:"omitempty"`
	IsReconciled *bool      `json:"is_reconciled" validate:"omitempty"`
	Page         int        `json:"page" validate:"omitempty,gte=1"`
	PageSize     int        `json:"page_size" validate:"omitempty,gte=1,lte=100"`
}

// ToModel converts CreateTransactionRequest to models.Transaction
func (r *CreateTransactionRequest) ToModel(userID uint) (*models.Transaction, error) {
	transactionDate, err := time.Parse(time.RFC3339, r.TransactionDate)
	if err != nil {
		return nil, err
	}

	tx := &models.Transaction{
		UserID:          userID,
		Type:            r.Type,
		Description:     r.Description,
		Amount:          r.Amount,
		AccountFromID:   r.AccountFromID,
		AccountToID:     r.AccountToID,
		CategoryID:      r.CategoryID,
		TransactionDate: transactionDate,
		Notes:           r.Notes,
		ExchangeRate:    r.ExchangeRate,
	}

	// Set default exchange rate if not provided
	if tx.ExchangeRate.IsZero() {
		tx.ExchangeRate = decimal.NewFromInt(1)
	}

	// Calculate USD amount
	tx.AmountInUSD = tx.CalculateAmountInUSD()

	return tx, nil
}

// FromModel converts models.Transaction to TransactionResponse
func FromModelToTransactionResponse(tx *models.Transaction) TransactionResponse {
	resp := TransactionResponse{
		ID:              tx.ID,
		UserID:          tx.UserID,
		Type:            tx.Type,
		Description:     tx.Description,
		Amount:          tx.Amount,
		AmountInUSD:     tx.AmountInUSD,
		ExchangeRate:    tx.ExchangeRate,
		TransactionDate: tx.TransactionDate,
		Notes:           tx.Notes,
		IsReconciled:    tx.IsReconciled,
		CreatedAt:       tx.CreatedAt,
		UpdatedAt:       tx.UpdatedAt,
	}

	// Include related data if loaded
	if tx.AccountFrom.ID != 0 {
		resp.AccountFrom = AccountSummary{
			ID:       tx.AccountFrom.ID,
			Name:     tx.AccountFrom.Name,
			Type:     tx.AccountFrom.AccountType,
			Balance:  tx.AccountFrom.Balance,
			Currency: nil, // Can be populated if needed
		}

		if tx.AccountFrom.Currency != nil {
			currencySummary := CurrencySummary{
				ID:     tx.AccountFrom.Currency.ID,
				Code:   tx.AccountFrom.Currency.Code,
				Symbol: tx.AccountFrom.Currency.Symbol,
			}
			resp.AccountFrom.Currency = &currencySummary
		}
	}

	if tx.AccountTo != nil && tx.AccountTo.ID != 0 {
		accountTo := AccountSummary{
			ID:       tx.AccountTo.ID,
			Name:     tx.AccountTo.Name,
			Type:     tx.AccountTo.AccountType,
			Balance:  tx.AccountTo.Balance,
			Currency: nil,
		}

		if tx.AccountTo.Currency != nil {
			currencySummary := CurrencySummary{
				ID:     tx.AccountTo.Currency.ID,
				Code:   tx.AccountTo.Currency.Code,
				Symbol: tx.AccountTo.Currency.Symbol,
			}
			accountTo.Currency = &currencySummary
		}

		resp.AccountTo = &accountTo
	}

	if tx.Category != nil && tx.Category.ID != 0 {
		resp.Category = &CategorySummary{
			ID:   tx.Category.ID,
			Name: tx.Category.Name,
			Type: tx.Category.Type,
		}
	}

	return resp
}

// FromModelToTransactionSummary converts models.Transaction to TransactionSummary
func FromModelToTransactionSummary(tx *models.Transaction) TransactionSummary {
	summary := TransactionSummary{
		ID:              tx.ID,
		Type:            tx.Type,
		Description:     tx.Description,
		Amount:          tx.Amount,
		TransactionDate: tx.TransactionDate,
		AccountFromName: tx.AccountFrom.Name,
		CategoryName:    nil,
	}

	if tx.Category != nil && tx.Category.ID != 0 {
		summary.CategoryName = &tx.Category.Name
	}

	return summary
}
