package validators

import (
	"arabella-api/internal/app/models"
	"fmt"

	"gorm.io/gorm"
)

// SystemValueValidator provides validation against system values
type SystemValueValidator struct {
	db *gorm.DB
}

// NewSystemValueValidator creates a new system value validator
func NewSystemValueValidator(db *gorm.DB) *SystemValueValidator {
	return &SystemValueValidator{db: db}
}

// ValidateValue checks if a value exists and is active for a given catalog type
func (v *SystemValueValidator) ValidateValue(catalogType, value string) error {
	var sv models.SystemValue

	err := v.db.Where("catalog_type = ? AND value = ? AND is_active = ?", catalogType, value, true).First(&sv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid %s: %s (not found in system values)", catalogType, value)
		}
		return fmt.Errorf("error validating %s: %w", catalogType, err)
	}

	return nil
}

// GetValidValues returns all valid values for a catalog type
func (v *SystemValueValidator) GetValidValues(catalogType string) ([]string, error) {
	var values []*models.SystemValue

	err := v.db.Where("catalog_type = ? AND is_active = ?", catalogType, true).
		Order("display_order ASC").
		Find(&values).Error

	if err != nil {
		return nil, err
	}

	result := make([]string, len(values))
	for i, sv := range values {
		result[i] = sv.Value
	}

	return result, nil
}
