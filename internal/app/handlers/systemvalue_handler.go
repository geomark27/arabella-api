package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SystemValueHandler handles system value HTTP requests
type SystemValueHandler struct {
	systemValueService services.SystemValueService
}

// NewSystemValueHandler creates a new system value handler
func NewSystemValueHandler(systemValueService services.SystemValueService) *SystemValueHandler {
	return &SystemValueHandler{
		systemValueService: systemValueService,
	}
}

// GetByCatalogType godoc
// @Summary      Obtener valores por tipo de catálogo
// @Description  Obtiene todos los valores activos de un catálogo específico del sistema. Catálogos disponibles: ACCOUNT_TYPE, ACCOUNT_CLASSIFICATION, TRANSACTION_TYPE, CATEGORY_TYPE, JOURNAL_ENTRY_TYPE. Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Param        catalogType  path      string  true  "Tipo de catálogo (ej: ACCOUNT_TYPE, TRANSACTION_TYPE)"
// @Success      200          {object}  dtos.SystemValueListResponseDTO  "Lista de valores del catálogo"
// @Failure      500          {object}  dtos.ErrorResponse               "Error interno del servidor"
// @Router       /system-values/catalog/{catalogType} [get]
func (h *SystemValueHandler) GetByCatalogType(c *gin.Context) {
	catalogType := c.Param("catalogType")

	values, err := h.systemValueService.GetByCatalogType(catalogType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve system values",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SystemValueListResponseDTO{
		Data:  values,
		Count: len(values),
	})
}

// GetAccountTypes godoc
// @Summary      Obtener tipos de cuenta
// @Description  Obtiene todos los tipos de cuenta disponibles: BANK (banco), CASH (efectivo), CREDIT_CARD (tarjeta de crédito), SAVINGS (ahorro), INVESTMENT (inversión). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  dtos.SystemValueListResponseDTO  "Lista de tipos de cuenta"
// @Failure      500  {object}  dtos.ErrorResponse               "Error interno del servidor"
// @Router       /system-values/account-types [get]
func (h *SystemValueHandler) GetAccountTypes(c *gin.Context) {
	values, err := h.systemValueService.GetAccountTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve account types",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SystemValueListResponseDTO{
		Data:  values,
		Count: len(values),
	})
}

// GetAccountClassifications godoc
// @Summary      Obtener clasificaciones de cuenta
// @Description  Obtiene todas las clasificaciones contables disponibles: ASSET (activo), LIABILITY (pasivo), EQUITY (patrimonio), INCOME (ingreso), EXPENSE (gasto). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  dtos.SystemValueListResponseDTO  "Lista de clasificaciones de cuenta"
// @Failure      500  {object}  dtos.ErrorResponse               "Error interno del servidor"
// @Router       /system-values/account-classifications [get]
func (h *SystemValueHandler) GetAccountClassifications(c *gin.Context) {
	values, err := h.systemValueService.GetAccountClassifications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve account classifications",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SystemValueListResponseDTO{
		Data:  values,
		Count: len(values),
	})
}

// GetTransactionTypes godoc
// @Summary      Obtener tipos de transacción
// @Description  Obtiene todos los tipos de transacción válidos: INCOME (requiere category_id), EXPENSE (requiere category_id), TRANSFER (requiere account_to_id). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  dtos.SystemValueListResponseDTO  "Lista de tipos de transacción"
// @Failure      500  {object}  dtos.ErrorResponse               "Error interno del servidor"
// @Router       /system-values/transaction-types [get]
func (h *SystemValueHandler) GetTransactionTypes(c *gin.Context) {
	values, err := h.systemValueService.GetTransactionTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve transaction types",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SystemValueListResponseDTO{
		Data:  values,
		Count: len(values),
	})
}

// GetCategoryTypes godoc
// @Summary      Obtener tipos de categoría
// @Description  Obtiene todos los tipos de categoría disponibles: INCOME (para ingresos) y EXPENSE (para gastos). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  dtos.SystemValueListResponseDTO  "Lista de tipos de categoría"
// @Failure      500  {object}  dtos.ErrorResponse               "Error interno del servidor"
// @Router       /system-values/category-types [get]
func (h *SystemValueHandler) GetCategoryTypes(c *gin.Context) {
	values, err := h.systemValueService.GetCategoryTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve category types",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SystemValueListResponseDTO{
		Data:  values,
		Count: len(values),
	})
}
