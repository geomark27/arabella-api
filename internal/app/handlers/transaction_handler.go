package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionService services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetTransactions retrieves transactions with optional filters
// GET /api/v1/transactions?type=EXPENSE&page=1&page_size=20
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	// TODO: Get userID from authentication context
	userID := uint(1)

	// Parse filters from query params
	filters := dtos.TransactionFilters{
		Type:     c.Query("type"),
		Page:     parseIntParam(c, "page", 1),
		PageSize: parseIntParam(c, "page_size", 20),
	}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)
		id := uint(accountID)
		filters.AccountID = &id
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
		id := uint(categoryID)
		filters.CategoryID = &id
	}

	result, err := h.transactionService.GetByUser(userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTransactionByID retrieves a specific transaction by ID
// GET /api/v1/transactions/:id
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	transaction, err := h.transactionService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Transaction not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

// CreateTransaction creates a new transaction
// POST /api/v1/transactions
// This goes through the Accounting Engine!
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dtos.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// TODO: Get userID from authentication context
	userID := uint(1)

	// This will process through the Accounting Engine
	transaction, err := h.transactionService.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction created successfully",
		"data":    transaction,
	})
}

// UpdateTransaction updates an existing transaction
// PUT /api/v1/transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	var req dtos.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.transactionService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated successfully",
	})
}

// DeleteTransaction deletes a transaction (by reversing it)
// DELETE /api/v1/transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	// This creates reversing journal entries (proper accounting)
	if err := h.transactionService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully (reversed)",
	})
}

// Helper function to parse int query params with default
func parseIntParam(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
