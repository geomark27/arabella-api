package models

import (
	"gorm.io/gorm"
)

// SystemValue represents system catalogs/lookups for dynamic select options
type SystemValue struct {
	gorm.Model
	CatalogType  string  `gorm:"type:varchar(50);not null;index:idx_system_values_catalog_value,priority:1;column:catalog_type" json:"catalog_type"`
	Value        string  `gorm:"type:varchar(100);not null;index:idx_system_values_catalog_value,priority:2" json:"value"`
	Label        string  `gorm:"type:varchar(200);not null" json:"label"`
	Description  *string `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int     `gorm:"default:0" json:"display_order"`
	IsActive     bool    `gorm:"default:true" json:"is_active"`
}

// TableName overrides the table name
func (SystemValue) TableName() string {
	return "system_values"
}
