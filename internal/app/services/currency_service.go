package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/repositories"
)

// CurrencyService handles currency-related business logic
type CurrencyService interface {
	GetAll() ([]dtos.CurrencyResponseDto, error)
	GetActive() ([]dtos.CurrencyResponseDto, error)
	GetByCode(code string) (*dtos.CurrencyResponseDto, error)
}

type currencyService struct {
	currencyRepo repositories.CurrencyRepository
}

// NewCurrencyService creates a new currency service
func NewCurrencyService(currencyRepo repositories.CurrencyRepository) CurrencyService {
	return &currencyService{
		currencyRepo: currencyRepo,
	}
}

// GetAll retrieves all currencies
func (s *currencyService) GetAll() ([]dtos.CurrencyResponseDto, error) {
	currencies, err := s.currencyRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.CurrencyResponseDto, len(currencies))
	for i, curr := range currencies {
		responses[i] = *dtos.ToCurrencyResponse(curr)
	}

	return responses, nil
}

// GetActive retrieves all active currencies
func (s *currencyService) GetActive() ([]dtos.CurrencyResponseDto, error) {
	currencies, err := s.currencyRepo.GetActiveCurrencies()
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.CurrencyResponseDto, len(currencies))
	for i, curr := range currencies {
		responses[i] = *dtos.ToCurrencyResponse(curr)
	}

	return responses, nil
}

// GetByCode retrieves a currency by its code (e.g., "USD")
func (s *currencyService) GetByCode(code string) (*dtos.CurrencyResponseDto, error) {
	currency, err := s.currencyRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	response := dtos.ToCurrencyResponse(currency)
	return response, nil
}
