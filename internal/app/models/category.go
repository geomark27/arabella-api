package models

import "gorm.io/gorm"

// Category represents the category entity in the database
type Category struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index" json:"user_id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Type     string `gorm:"size:50;not null" json:"type"` // e.g., "income" or "expense"
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// TableName overrides the table name (optional)
func (Category) TableName() string {
	return "categorys"
}
