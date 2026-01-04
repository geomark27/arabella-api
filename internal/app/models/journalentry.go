package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// JournalEntry represents a single entry in the double-entry bookkeeping system
// Every transaction creates at least 2 journal entries (Debit + Credit)
type JournalEntry struct {
	gorm.Model
	UserID        uint            `gorm:"index;not null" json:"user_id"`
	TransactionID uint            `gorm:"index;not null" json:"transaction_id"`
	AccountID     uint            `gorm:"index;not null" json:"account_id"`                                                     // Account where the movement is registered
	DebitOrCredit string          `gorm:"size:10;not null;check:debit_or_credit IN ('DEBIT', 'CREDIT')" json:"debit_or_credit"` // "DEBIT" or "CREDIT"
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	EntryDate     time.Time       `gorm:"index;not null" json:"entry_date"`
	Description   string          `gorm:"size:255" json:"description"`

	// Relationships
	Account     Account     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}

// TableName overrides the table name
func (JournalEntry) TableName() string {
	return "journal_entries"
}

// Validate performs business rule validation on the JournalEntry
func (j *JournalEntry) Validate() error {
	if j.UserID == 0 {
		return errors.New("user_id is required")
	}

	if j.TransactionID == 0 {
		return errors.New("transaction_id is required")
	}

	if j.AccountID == 0 {
		return errors.New("account_id is required")
	}

	// Validate DebitOrCredit
	if j.DebitOrCredit != "DEBIT" && j.DebitOrCredit != "CREDIT" {
		return fmt.Errorf("debit_or_credit must be DEBIT or CREDIT, got: %s", j.DebitOrCredit)
	}

	// Validate amount is positive
	if j.Amount.IsZero() || j.Amount.IsNegative() {
		return fmt.Errorf("amount must be positive, got: %s", j.Amount.String())
	}

	// Validate entry date
	if j.EntryDate.IsZero() {
		return errors.New("entry_date is required")
	}

	return nil
}

// IsDebit returns true if this entry is a debit
func (j *JournalEntry) IsDebit() bool {
	return j.DebitOrCredit == "DEBIT"
}

// IsCredit returns true if this entry is a credit
func (j *JournalEntry) IsCredit() bool {
	return j.DebitOrCredit == "CREDIT"
}
