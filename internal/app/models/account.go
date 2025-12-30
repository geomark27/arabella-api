package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	AccountType string `gorm:"size:50;not null" json:"account_type"`
	CurrencyID  uint   `gorm:"not null" json:"currency_id"`
	Balance     string `gorm:"type:decimal(11,4);default:0" json:"balance"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

func (Account) TableName() string {
	return "accounts"
}
