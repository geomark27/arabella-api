package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// JournalEntryHandler handles journal entry-related HTTP requests
// Note: Journal entries are READ-ONLY (created automatically by Accounting Engine)
type JournalEntryHandler struct {
	journalEntryService services.JournalEntryService
}

// NewJournalEntryHandler creates a new journal entry handler
func NewJournalEntryHandler(journalEntryService services.JournalEntryService) *JournalEntryHandler {
	return &JournalEntryHandler{
		journalEntryService: journalEntryService,
	}
}

// GetJournalEntries retrieves journal entries (audit trail)
// GET /api/v1/journal-entries?transaction_id=123&page=1
func (h *JournalEntryHandler) GetJournalEntries(c *gin.Context) {
	// TODO: Get userID from authentication context
	userID := uint(1)

	filters := dtos.JournalEntryFilters{
		Page:     parseIntParam(c, "page", 1),
		PageSize: parseIntParam(c, "page_size", 50),
	}

	if txIDStr := c.Query("transaction_id"); txIDStr != "" {
		txID, _ := strconv.ParseUint(txIDStr, 10, 32)
		id := uint(txID)
		filters.TransactionID = &id
	}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		accountID, _ := strconv.ParseUint(accountIDStr, 10, 32)
		id := uint(accountID)
		filters.AccountID = &id
	}

	if debitCredit := c.Query("debit_or_credit"); debitCredit != "" {
		filters.DebitOrCredit = &debitCredit
	}

	result, err := h.journalEntryService.GetByUser(userID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve journal entries",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetJournalEntriesByTransaction retrieves journal entries for a specific transaction
// GET /api/v1/journal-entries/transaction/:id
func (h *JournalEntryHandler) GetJournalEntriesByTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	entries, err := h.journalEntryService.GetByTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve journal entries",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  entries,
		"count": len(entries),
	})
}

// VerifyTransactionBalance verifies that debits = credits for a transaction
// GET /api/v1/journal-entries/verify/:id
func (h *JournalEntryHandler) VerifyTransactionBalance(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	verification, err := h.journalEntryService.VerifyTransactionBalance(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to verify transaction balance",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": verification,
	})
}
