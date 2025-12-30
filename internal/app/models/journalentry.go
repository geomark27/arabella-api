package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// JournalEntry represents the journalentry entity in the database
type JournalEntry struct {
	gorm.Model
	UserID        uint            `gorm:"index;not null" json:"user_id"`
	TransactionID uint            `gorm:"index;not null" json:"transaction_id"`
	AccountID     uint            `gorm:"index;not null" json:"account_id"`        // Cuenta donde se registra el movimiento
	DebitOrCredit string          `gorm:"size:10;not null" json:"debit_or_credit"` // "DEBIT", "CREDIT"
	Amount        decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"amount"`
	EntryDate     time.Time       `gorm:"index" json:"entry_date"`
	Description   string          `gorm:"size:255" json:"description"`
}

// TableName overrides the table name (optional)
func (JournalEntry) TableName() string {
	return "journalentrys"
}
