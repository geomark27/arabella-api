package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"arabella-api/internal/app/handlers"
	"arabella-api/internal/app/repositories"
	"arabella-api/internal/app/services"
	"arabella-api/internal/platform/config"
	"arabella-api/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpServer *http.Server
}

// New creates a new server instance with all dependencies injected
func New(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ========================================
	// DEPENDENCY INJECTION
	// ========================================

	// Warn if JWT secrets are still using insecure defaults
	if cfg.JWTSecret == "your-default-secret-change-in-production" ||
		cfg.JWTRefreshSecret == "your-refresh-secret-change-in-production" {
		log.Println("⚠️  WARNING: JWT secrets are using insecure default values. Set JWT_SECRET and JWT_REFRESH_SECRET environment variables.")
	}

	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	journalEntryRepo := repositories.NewJournalEntryRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	currencyRepo := repositories.NewCurrencyRepository(db)
	systemValueRepo := repositories.NewSystemValueRepository(db)

	// Create services (injecting repositories)
	jwtService := services.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, jwtService, db)
	accountingEngine := services.NewAccountingEngineService(db, journalEntryRepo, accountRepo, transactionRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountingEngine)
	accountService := services.NewAccountService(accountRepo, systemValueRepo)
	dashboardService := services.NewDashboardService(accountRepo, transactionRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	currencyService := services.NewCurrencyService(currencyRepo)
	systemValueService := services.NewSystemValueService(systemValueRepo)
	journalEntryService := services.NewJournalEntryService(journalEntryRepo)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Create handlers (injecting services)
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	currencyHandler := handlers.NewCurrencyHandler(currencyService)
	systemValueHandler := handlers.NewSystemValueHandler(systemValueService)
	journalEntryHandler := handlers.NewJournalEntryHandler(journalEntryService)

	// Create Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(corsMiddleware(cfg.CorsAllowedOrigins))

	// Register routes
	registerRoutes(
		router,
		authMiddleware,
		healthHandler,
		authHandler,
		userHandler,
		accountHandler,
		transactionHandler,
		categoryHandler,
		currencyHandler,
		systemValueHandler,
		journalEntryHandler,
		dashboardHandler,
	)

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:     cfg,
		router:     router,
		httpServer: httpServer,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// corsMiddleware returns a Gin middleware for CORS
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Company-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
