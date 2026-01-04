package models

import "gorm.io/gorm"

// Currency represents the currency entity in the database
type Currency struct {
	gorm.Model
	Code     string `gorm:"size:10;not null;uniqueIndex" json:"code"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Symbol   string `gorm:"size:10;not null" json:"symbol"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// TableName overrides the table name (optional)
func (Currency) TableName() string {
	return "currencies"
}
