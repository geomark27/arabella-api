package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction representa la entidad de transacción con lógica financiera robusta
type Transaction struct {
	gorm.Model
	UserID          uint            `gorm:"not null;index" json:"user_id"`
	Type            string          `gorm:"type:varchar(20);not null;check:type IN ('INCOME', 'EXPENSE', 'TRANSFER')" json:"type"`
	Description     string          `gorm:"size:255;not null" json:"description"`
	Amount          decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	AmountInUSD     decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"amount_in_usd"`
	ExchangeRate    decimal.Decimal `gorm:"type:decimal(19,6);default:1" json:"exchange_rate"`
	AccountFromID   uint            `gorm:"index;not null" json:"account_from_id"`
	AccountToID     *uint           `gorm:"index" json:"account_to_id"` // Puntero porque es opcional (solo TRANSFER)
	CategoryID      *uint           `gorm:"index" json:"category_id"`   // Puntero porque es opcional (solo INCOME/EXPENSE)
	TransactionDate time.Time       `gorm:"index;not null" json:"transaction_date"`
	Notes           string          `gorm:"type:text" json:"notes"`
	IsReconciled    bool            `gorm:"default:false" json:"is_reconciled"`
	AccountFrom     Account         `gorm:"foreignKey:AccountFromID" json:"account_from,omitempty"`
	AccountTo       *Account        `gorm:"foreignKey:AccountToID" json:"account_to,omitempty"`
	Category        *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TableName define el nombre de la tabla
func (Transaction) TableName() string {
	return "transactions"
}

// Validate performs business rule validation on the Transaction
func (t *Transaction) Validate() error {
	// Validate amount is positive
	if t.Amount.IsZero() || t.Amount.IsNegative() {
		return fmt.Errorf("amount must be positive, got: %s", t.Amount.String())
	}

	if t.Type == "" {
		return errors.New("transaction type is required")
	}

	// Basic validation - detailed type validation should be done at service layer with SystemValue
	// Valid types are: INCOME, EXPENSE, TRANSFER, DEBT_PAYMENT
	// (These should match SystemValue.Category='TRANSACTION_TYPE')

	// Type-specific validation
	switch t.Type {
	case "TRANSFER":
		if t.AccountToID == nil {
			return errors.New("account_to_id is required for TRANSFER transactions")
		}
		if t.CategoryID != nil {
			return errors.New("category_id should not be set for TRANSFER transactions")
		}
	case "INCOME", "EXPENSE":
		if t.CategoryID == nil {
			return fmt.Errorf("category_id is required for %s transactions", t.Type)
		}
		if t.AccountToID != nil {
			return fmt.Errorf("account_to_id should not be set for %s transactions", t.Type)
		}
	}

	// Validate description
	if t.Description == "" {
		return errors.New("description is required")
	}

	// Validate exchange rate is positive
	if !t.ExchangeRate.IsZero() && t.ExchangeRate.IsNegative() {
		return fmt.Errorf("exchange_rate must be positive, got: %s", t.ExchangeRate.String())
	}

	// Validate transaction date is set
	if t.TransactionDate.IsZero() {
		return errors.New("transaction_date is required")
	}

	return nil
}

// CalculateAmountInUSD calculates the USD equivalent amount using the exchange rate
// If exchange rate is not set (0 or 1), assumes the amount is already in the base currency
func (t *Transaction) CalculateAmountInUSD() decimal.Decimal {
	if t.ExchangeRate.IsZero() || t.ExchangeRate.Equal(decimal.NewFromInt(1)) {
		return t.Amount
	}
	return t.Amount.Mul(t.ExchangeRate)
}

// IsMultiCurrency returns true if this transaction involves currency conversion
func (t *Transaction) IsMultiCurrency() bool {
	return !t.ExchangeRate.IsZero() && !t.ExchangeRate.Equal(decimal.NewFromInt(1))
}
