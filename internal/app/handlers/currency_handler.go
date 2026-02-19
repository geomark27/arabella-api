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

// GetCurrencies godoc
// @Summary      Listar monedas disponibles
// @Description  Obtiene todas las monedas activas disponibles en el sistema (USD, EUR, MXN, etc.). Este endpoint es público y no requiere autenticación
// @Tags         Currencies
// @Produce      json
// @Success      200  {object}  object{data=[]dtos.CurrencyResponseDto,count=int}  "Lista de monedas activas"
// @Failure      500  {object}  dtos.ErrorResponse                                 "Error interno del servidor"
// @Router       /currencies [get]
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

// GetCurrencyByCode godoc
// @Summary      Obtener moneda por código ISO
// @Description  Obtiene el detalle de una moneda específica usando su código ISO 4217 (ej: USD, EUR, MXN). Este endpoint es público y no requiere autenticación
// @Tags         Currencies
// @Produce      json
// @Param        code  path      string                              true  "Código ISO 4217 de la moneda (ej: USD, EUR, MXN)"
// @Success      200   {object}  object{data=dtos.CurrencyResponseDto}  "Detalle de la moneda"
// @Failure      404   {object}  dtos.ErrorResponse                  "Moneda no encontrada"
// @Router       /currencies/{code} [get]
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
