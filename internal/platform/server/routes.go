package server

import (
	"arabella-api/internal/app/handlers"
	"arabella-api/internal/shared/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swaggerUI5HTML es la plantilla HTML que carga Swagger UI 5.x desde CDN.
// Esta versión incluye el toggle de dark/light mode nativo (ícono 💡 arriba a la derecha).
// El spec se carga desde /swagger/doc.json generado por swag init.
const swaggerUI5HTML = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Arabella Financial OS — API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    /* Ocultar el topbar azul oscuro por defecto de StandaloneLayout */
    .swagger-ui .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        persistAuthorization: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

// registerRoutes registers all application routes
func registerRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	accountHandler *handlers.AccountHandler,
	transactionHandler *handlers.TransactionHandler,
	categoryHandler *handlers.CategoryHandler,
	currencyHandler *handlers.CurrencyHandler,
	systemValueHandler *handlers.SystemValueHandler,
	journalEntryHandler *handlers.JournalEntryHandler,
	dashboardHandler *handlers.DashboardHandler,
) {
	// ------------------------------------------------------------------
	// Swagger UI clásico  →  /swagger/index.html  (servido por swaggo)
	// Swagger UI 5.x CDN  →  /docs               (con dark mode nativo)
	// ------------------------------------------------------------------
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUI5HTML)
	})
	// Redirección conveniente: /swagger → /docs
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Welcome to Arabella Financial OS API!",
			"status":      "success",
			"version":     "v1.0.0 - Phase 1",
			"description": "Personal Financial Management System with Double-Entry Bookkeeping",
			"endpoints": gin.H{
				"docs":            "/docs",
				"health":          "/api/v1/health",
				"auth":            "/api/v1/auth",
				"dashboard":       "/api/v1/dashboard",
				"users":           "/api/v1/users",
				"accounts":        "/api/v1/accounts",
				"transactions":    "/api/v1/transactions",
				"categories":      "/api/v1/categories",
				"currencies":      "/api/v1/currencies",
				"system_values":   "/api/v1/system-values",
				"journal_entries": "/api/v1/journal-entries",
			},
		})
	})

	// ============================================================
	// PUBLIC ROUTES — No authentication required
	// ============================================================
	public := router.Group("/api/v1")
	{
		// Health routes
		public.GET("/health", healthHandler.Health)
		public.GET("/health/ready", healthHandler.Ready)

		// Auth routes (register & login are always public)
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// System Value routes (public catalog data — read-only)
		systemValues := public.Group("/system-values")
		{
			systemValues.GET("/catalog/:catalogType", systemValueHandler.GetByCatalogType)
			systemValues.GET("/account-types", systemValueHandler.GetAccountTypes)
			systemValues.GET("/account-classifications", systemValueHandler.GetAccountClassifications)
			systemValues.GET("/transaction-types", systemValueHandler.GetTransactionTypes)
			systemValues.GET("/category-types", systemValueHandler.GetCategoryTypes)
		}

		// Currency routes (public catalog data — read-only)
		currencies := public.Group("/currencies")
		{
			currencies.GET("", currencyHandler.GetCurrencies)
			currencies.GET("/:code", currencyHandler.GetCurrencyByCode)
		}
	}

	// ============================================================
	// PROTECTED ROUTES — JWT authentication required
	// ============================================================
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.RequireAuth())
	{
		// Auth — actions that require an active session
		protected.PUT("/auth/change-password", authHandler.ChangePassword)

		// Dashboard routes (Feature Star: Runway Calculation)
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("", dashboardHandler.GetDashboard)
			dashboard.GET("/runway", dashboardHandler.GetRunway)
			dashboard.GET("/monthly-stats", dashboardHandler.GetMonthlyStats)
		}

		// User routes
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUserByID)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Account routes
		accounts := protected.Group("/accounts")
		{
			accounts.GET("", accountHandler.GetAccounts)
			accounts.POST("", accountHandler.CreateAccount)
			accounts.GET("/:id", accountHandler.GetAccountByID)
			accounts.PUT("/:id", accountHandler.UpdateAccount)
			accounts.DELETE("/:id", accountHandler.DeleteAccount)
		}

		// Transaction routes (Goes through Accounting Engine)
		transactions := protected.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("/:id", transactionHandler.GetTransactionByID)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}

		// Category routes
		categories := protected.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Journal Entry routes (Read-only, audit trail)
		journalEntries := protected.Group("/journal-entries")
		{
			journalEntries.GET("", journalEntryHandler.GetJournalEntries)
			journalEntries.GET("/transaction/:id", journalEntryHandler.GetJournalEntriesByTransaction)
			journalEntries.GET("/verify/:id", journalEntryHandler.VerifyTransactionBalance)
		}
	}
}
