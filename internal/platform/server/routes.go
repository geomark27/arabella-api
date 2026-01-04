package server

import (
	"arabella-api/internal/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerRoutes registers all application routes
func registerRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	userHandler *handlers.UserHandler,
	accountHandler *handlers.AccountHandler,
	transactionHandler *handlers.TransactionHandler,
	categoryHandler *handlers.CategoryHandler,
	currencyHandler *handlers.CurrencyHandler,
	systemValueHandler *handlers.SystemValueHandler,
	journalEntryHandler *handlers.JournalEntryHandler,
	dashboardHandler *handlers.DashboardHandler,
) {
	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to Arabella Financial OS API!",
			"status":      "success",
			"version":     "v1.0.0 - Phase 1",
			"description": "Personal Financial Management System with Double-Entry Bookkeeping",
			"endpoints": gin.H{
				"health":          "/api/v1/health",
				"dashboard":       "/api/v1/dashboard",
				"users":           "/api/v1/users",
				"accounts":        "/api/v1/accounts",
				"transactions":    "/api/v1/transactions",
				"categories":      "/api/v1/categories",
				"currencies":      "/api/v1/currencies",
				"system_values":   "/api/v1/system-values",
				"journal_entries": "/api/v1/journal-entries",
				"docs":            "/docs",
			},
		})
	})

	// API v1 group
	api := router.Group("/api/v1")
	{
		// Health routes
		api.GET("/health", healthHandler.Health)
		api.GET("/health/ready", healthHandler.Ready)

		// Dashboard routes (Feature Star: Runway Calculation)
		api.GET("/dashboard", dashboardHandler.GetDashboard)
		api.GET("/dashboard/runway", dashboardHandler.GetRunway)
		api.GET("/dashboard/monthly-stats", dashboardHandler.GetMonthlyStats)

		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Account routes
		accounts := api.Group("/accounts")
		{
			accounts.GET("", accountHandler.GetAccounts)
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:id", accountHandler.GetAccountByID)
			accounts.PUT("/:id", accountHandler.UpdateAccount)
			accounts.DELETE("/:id", accountHandler.DeleteAccount)
		}

		// Transaction routes (Goes through Accounting Engine)
		transactions := api.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Currency routes (Read-only for now)
		currencies := api.Group("/currencies")
		{
			currencies.GET("", currencyHandler.GetCurrencies)
			currencies.GET("/:code", currencyHandler.GetCurrencyByCode)
		}

		// System Value routes (Catalogs - Read-only)
		systemValues := api.Group("/system-values")
		{
			systemValues.GET("/catalog/:catalogType", systemValueHandler.GetByCatalogType)
			systemValues.GET("/account-types", systemValueHandler.GetAccountTypes)
			systemValues.GET("/account-classifications", systemValueHandler.GetAccountClassifications)
			systemValues.GET("/transaction-types", systemValueHandler.GetTransactionTypes)
			systemValues.GET("/category-types", systemValueHandler.GetCategoryTypes)
		}

		// Journal Entry routes (Read-only, audit trail)
		journalEntries := api.Group("/journal-entries")
		{
			journalEntries.GET("", journalEntryHandler.GetJournalEntries)
			journalEntries.GET("/transaction/:id", journalEntryHandler.GetJournalEntriesByTransaction)
			journalEntries.GET("/verify/:id", journalEntryHandler.VerifyTransactionBalance)
		}
	}
}
