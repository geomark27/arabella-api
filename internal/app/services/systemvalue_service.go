package services

import (
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"fmt"
)

// SystemValueService handles system catalog values
type SystemValueService interface {
	GetByCatalogType(catalogType string) ([]*models.SystemValue, error)
	GetAccountTypes() ([]*models.SystemValue, error)
	GetAccountClassifications() ([]*models.SystemValue, error)
	GetTransactionTypes() ([]*models.SystemValue, error)
	GetCategoryTypes() ([]*models.SystemValue, error)
	GetJournalEntryTypes() ([]*models.SystemValue, error)
	ValidateValue(catalogType, value string) error
}

type systemValueService struct {
	systemValueRepo repositories.SystemValueRepository
}

// NewSystemValueService creates a new system value service
func NewSystemValueService(systemValueRepo repositories.SystemValueRepository) SystemValueService {
	return &systemValueService{
		systemValueRepo: systemValueRepo,
	}
}

// GetByCatalogType retrieves all system values for a catalog type
func (s *systemValueService) GetByCatalogType(catalogType string) ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType(catalogType)
}

// GetAccountTypes retrieves all account types
func (s *systemValueService) GetAccountTypes() ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType("ACCOUNT_TYPE")
}

// GetAccountClassifications retrieves all account classifications
func (s *systemValueService) GetAccountClassifications() ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType("ACCOUNT_CLASSIFICATION")
}

// GetTransactionTypes retrieves all transaction types
func (s *systemValueService) GetTransactionTypes() ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType("TRANSACTION_TYPE")
}

// GetCategoryTypes retrieves all category types
func (s *systemValueService) GetCategoryTypes() ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType("CATEGORY_TYPE")
}

// GetJournalEntryTypes retrieves all journal entry types
func (s *systemValueService) GetJournalEntryTypes() ([]*models.SystemValue, error) {
	return s.systemValueRepo.FindByCatalogType("JOURNAL_ENTRY_TYPE")
}

// ValidateValue checks if a value exists and is active for a given catalog type
func (s *systemValueService) ValidateValue(catalogType, value string) error {
	values, err := s.systemValueRepo.FindByCatalogType(catalogType)
	if err != nil {
		return err
	}

	for _, sv := range values {
		if sv.Value == value && sv.IsActive {
			return nil
		}
	}

	return fmt.Errorf("invalid %s: %s", catalogType, value)
}
