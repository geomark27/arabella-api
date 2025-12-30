package models

import (
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
