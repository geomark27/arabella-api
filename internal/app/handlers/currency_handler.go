package handlers

import (
	"arabella-api/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CurrencyHandler handles currency-related HTTP requests
type CurrencyHandler struct {
	currencyService services.CurrencyService
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(currencyService services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
	}
}

// GetCurrencies retrieves all active currencies
// GET /api/v1/currencies
func (h *CurrencyHandler) GetCurrencies(c *gin.Context) {
	currencies, err := h.currencyService.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve currencies",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  currencies,
		"count": len(currencies),
	})
}

// GetCurrencyByCode retrieves a currency by its code (e.g., "USD")
// GET /api/v1/currencies/:code
func (h *CurrencyHandler) GetCurrencyByCode(c *gin.Context) {
	code := c.Param("code")

	currency, err := h.currencyService.GetByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Currency not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": currency,
	})
}
