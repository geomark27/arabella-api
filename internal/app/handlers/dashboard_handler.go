package handlers

import (
	"arabella-api/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardService services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboard godoc
// @Summary      Resumen financiero del usuario
// @Description  Obtiene el dashboard completo del usuario autenticado: patrimonio neto, activos, pasivos, flujo de caja mensual, y el cálculo de Runway (cuántos meses puede sobrevivir sin ingresos basado en el promedio de gastos de los últimos 3 meses)
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  object{data=dtos.DashboardResponse}  "Dashboard financiero completo"
// @Failure      401  {object}  dtos.ErrorResponse                   "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                   "Error interno del servidor"
// @Security     BearerAuth
// @Router       /dashboard [get]
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

	dashboard, err := h.dashboardService.GetDashboard(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve dashboard data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dashboard,
	})
}

// GetRunway godoc
// @Summary      Cálculo detallado de Runway
// @Description  Calcula cuántos meses y días puede el usuario sostenerse financieramente sin nuevos ingresos. Fórmula: (Activos líquidos - Pasivos a corto plazo) / Promedio de gastos mensuales (últimos 3 meses). Estados: HEALTHY (≥6 meses), WARNING (3-6 meses), CRITICAL (<3 meses)
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  object{data=dtos.RunwayCalculation}  "Cálculo de Runway detallado con desglose por tipo de cuenta"
// @Failure      401  {object}  dtos.ErrorResponse                   "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                   "Error interno del servidor"
// @Security     BearerAuth
// @Router       /dashboard/runway [get]
func (h *DashboardHandler) GetRunway(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

	runway, err := h.dashboardService.CalculateRunway(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to calculate runway",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": runway,
	})
}

// GetMonthlyStats godoc
// @Summary      Estadísticas mensuales de ingresos y gastos
// @Description  Obtiene el resumen de ingresos, gastos y flujo neto de caja para un mes y año específicos del usuario autenticado
// @Tags         Dashboard
// @Produce      json
// @Param        month  query     int  false  "Mes (1-12, default: mes actual)"
// @Param        year   query     int  false  "Año (ej: 2026, default: año actual)"
// @Success      200    {object}  object{data=dtos.MonthlyStats}  "Estadísticas del mes solicitado"
// @Failure      401    {object}  dtos.ErrorResponse              "No autenticado"
// @Failure      500    {object}  dtos.ErrorResponse              "Error interno del servidor"
// @Security     BearerAuth
// @Router       /dashboard/monthly-stats [get]
func (h *DashboardHandler) GetMonthlyStats(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

	// Default to current month/year if not provided
	month := parseIntParam(c, "month", 1)
	year := parseIntParam(c, "year", 2026)

	stats, err := h.dashboardService.GetMonthlyStats(userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve monthly stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}
