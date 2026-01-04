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

// GetDashboard retrieves complete dashboard data including the Runway calculation
// GET /api/v1/dashboard
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// TODO: Get userID from authentication context
	userID := uint(1)

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

// GetRunway retrieves detailed runway calculation
// GET /api/v1/dashboard/runway
func (h *DashboardHandler) GetRunway(c *gin.Context) {
	// TODO: Get userID from authentication context
	userID := uint(1)

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

// GetMonthlyStats retrieves monthly statistics
// GET /api/v1/dashboard/monthly-stats?month=1&year=2026
func (h *DashboardHandler) GetMonthlyStats(c *gin.Context) {
	// TODO: Get userID from authentication context
	userID := uint(1)

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
