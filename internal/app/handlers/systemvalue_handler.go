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

// GetByCatalogType retrieves all system values for a catalog type
// GET /api/v1/system-values/catalog/:catalogType
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

// GetAccountTypes retrieves all account types
// GET /api/v1/system-values/account-types
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

// GetAccountClassifications retrieves all account classifications
// GET /api/v1/system-values/account-classifications
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

// GetTransactionTypes retrieves all transaction types
// GET /api/v1/system-values/transaction-types
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

// GetCategoryTypes retrieves all category types
// GET /api/v1/system-values/category-types
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
