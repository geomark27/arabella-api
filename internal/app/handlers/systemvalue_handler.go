package handlers

import (
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
// @Description  Obtiene todos los valores activos de un catálogo específico del sistema. Los catálogos disponibles son: ACCOUNT_TYPE, ACCOUNT_CLASSIFICATION, TRANSACTION_TYPE, CATEGORY_TYPE. Este endpoint es público y no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Param        catalogType  path      string  true  "Tipo de catálogo (ej: ACCOUNT_TYPE, TRANSACTION_TYPE)"
// @Success      200          {object}  object{data=[]interface{},count=int}  "Lista de valores del catálogo"
// @Failure      500          {object}  dtos.ErrorResponse                    "Error interno del servidor"
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

	c.JSON(http.StatusOK, gin.H{
		"data":  values,
		"count": len(values),
	})
}

// GetAccountTypes godoc
// @Summary      Obtener tipos de cuenta
// @Description  Obtiene todos los tipos de cuenta disponibles en el sistema: BANK (banco), CASH (efectivo), CREDIT_CARD (tarjeta de crédito), SAVINGS (ahorro), INVESTMENT (inversión). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  object{data=[]interface{},count=int}  "Lista de tipos de cuenta"
// @Failure      500  {object}  dtos.ErrorResponse                    "Error interno del servidor"
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

	c.JSON(http.StatusOK, gin.H{
		"data":  values,
		"count": len(values),
	})
}

// GetAccountClassifications godoc
// @Summary      Obtener clasificaciones de cuenta
// @Description  Obtiene todas las clasificaciones contables de cuenta disponibles en el sistema (ej: ASSET, LIABILITY, EQUITY). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  object{data=[]interface{},count=int}  "Lista de clasificaciones de cuenta"
// @Failure      500  {object}  dtos.ErrorResponse                    "Error interno del servidor"
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

	c.JSON(http.StatusOK, gin.H{
		"data":  values,
		"count": len(values),
	})
}

// GetTransactionTypes godoc
// @Summary      Obtener tipos de transacción
// @Description  Obtiene todos los tipos de transacción válidos en el sistema: INCOME (ingreso, requiere category_id), EXPENSE (gasto, requiere category_id), TRANSFER (transferencia entre cuentas, requiere account_to_id). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  object{data=[]interface{},count=int}  "Lista de tipos de transacción"
// @Failure      500  {object}  dtos.ErrorResponse                    "Error interno del servidor"
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

	c.JSON(http.StatusOK, gin.H{
		"data":  values,
		"count": len(values),
	})
}

// GetCategoryTypes godoc
// @Summary      Obtener tipos de categoría
// @Description  Obtiene todos los tipos de categoría disponibles en el sistema: INCOME (para clasificar ingresos) y EXPENSE (para clasificar gastos). Endpoint público, no requiere autenticación
// @Tags         System Values
// @Produce      json
// @Success      200  {object}  object{data=[]interface{},count=int}  "Lista de tipos de categoría"
// @Failure      500  {object}  dtos.ErrorResponse                    "Error interno del servidor"
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

	c.JSON(http.StatusOK, gin.H{
		"data":  values,
		"count": len(values),
	})
}
