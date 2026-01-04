package dtos

import (
	"arabella-api/internal/app/models"
	"time"
)

type CurrencyResponseDto struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	IsActive  bool      `json:"is_active"`
}

type UpdateOrCreateDto struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	IsActive bool   `json:"is_active"`
}

func ToCurrencyResponse(currency *models.Currency) *CurrencyResponseDto {
	return &CurrencyResponseDto{
		ID:        currency.ID,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
		Code:      currency.Code,
		Name:      currency.Name,
		Symbol:    currency.Symbol,
		IsActive:  currency.IsActive,
	}
}

func ToCurrencyResponseList(currencies []models.Currency) []CurrencyResponseDto {
	result := make([]CurrencyResponseDto, len(currencies))
	for i, currency := range currencies {
		result[i] = *ToCurrencyResponse(&currency)
	}
	return result
}

// CurrencySummary represents a lightweight currency for use in other DTOs
type CurrencySummary struct {
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}
