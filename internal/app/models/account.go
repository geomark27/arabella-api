package models

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID      uint            `gorm:"not null;index" json:"user_id"`
	Name        string          `gorm:"size:100;not null" json:"name"`
	AccountType string          `gorm:"size:50;not null" json:"account_type"` // BANK, CASH, CREDIT_CARD, SAVINGS, INVESTMENT
	CurrencyID  *uint           `gorm:"not null" json:"currency_id"`
	Balance     decimal.Decimal `gorm:"type:decimal(19,4);default:0" json:"balance"` // Upgraded precision per DOC.md
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	Currency    *Currency       `gorm:"foreignKey:CurrencyID;references:ID" json:"currency,omitempty"`
}

func (Account) TableName() string {
	return "accounts"
}

// Validate performs business rule validation on the Account
func (a *Account) Validate() error {
	if a.UserID == 0 {
		return errors.New("user_id is required")
	}

	if a.Name == "" {
		return errors.New("account name is required")
	}

	if len(a.Name) > 100 {
		return errors.New("account name cannot exceed 100 characters")
	}

	if a.AccountType == "" {
		return errors.New("account_type is required")
	}

	// Basic validation - detailed validation should be done at service layer with SystemValue
	// Valid types are: BANK, CASH, CREDIT_CARD, SAVINGS, INVESTMENT, CATEGORY
	// (These should match SystemValue.Category='ACCOUNT_TYPE')

	if a.CurrencyID == nil {
		return errors.New("currency_id is required")
	}

	return nil
}

// UpdateBalance updates the account balance by adding the given amount
// Use positive amounts to increase balance, negative to decrease
func (a *Account) UpdateBalance(amount decimal.Decimal) {
	a.Balance = a.Balance.Add(amount)
}

// IsLiquidAsset returns true if this account represents a liquid asset (BANK or CASH)
func (a *Account) IsLiquidAsset() bool {
	return a.AccountType == "BANK" || a.AccountType == "CASH"
}

// IsLiability returns true if this account represents a liability (CREDIT_CARD)
func (a *Account) IsLiability() bool {
	return a.AccountType == "CREDIT_CARD"
}
