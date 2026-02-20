package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/repositories"
	"fmt"
)

// SystemValueService handles system catalog values
type SystemValueService interface {
	GetByCatalogType(catalogType string) ([]dtos.SystemValueResponseDTO, error)
	GetAccountTypes() ([]dtos.SystemValueResponseDTO, error)
	GetAccountClassifications() ([]dtos.SystemValueResponseDTO, error)
	GetTransactionTypes() ([]dtos.SystemValueResponseDTO, error)
	GetCategoryTypes() ([]dtos.SystemValueResponseDTO, error)
	GetJournalEntryTypes() ([]dtos.SystemValueResponseDTO, error)
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
func (s *systemValueService) GetByCatalogType(catalogType string) ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType(catalogType)
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
}

// GetAccountTypes retrieves all account types
func (s *systemValueService) GetAccountTypes() ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType("ACCOUNT_TYPE")
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
}

// GetAccountClassifications retrieves all account classifications
func (s *systemValueService) GetAccountClassifications() ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType("ACCOUNT_CLASSIFICATION")
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
}

// GetTransactionTypes retrieves all transaction types
func (s *systemValueService) GetTransactionTypes() ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType("TRANSACTION_TYPE")
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
}

// GetCategoryTypes retrieves all category types
func (s *systemValueService) GetCategoryTypes() ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType("CATEGORY_TYPE")
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
}

// GetJournalEntryTypes retrieves all journal entry types
func (s *systemValueService) GetJournalEntryTypes() ([]dtos.SystemValueResponseDTO, error) {
	values, err := s.systemValueRepo.FindByCatalogType("JOURNAL_ENTRY_TYPE")
	if err != nil {
		return nil, err
	}
	return dtos.ToSystemValueResponseList(values), nil
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
