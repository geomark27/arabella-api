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

// GetJournalEntries godoc
// @Summary      Listar asientos contables (pista de auditoría)
// @Description  Obtiene una lista paginada de asientos del diario contable del usuario autenticado. Los asientos son generados automáticamente por el Motor Contable al crear o revertir transacciones; nunca se crean ni modifican manualmente. Se pueden filtrar por transacción, cuenta o tipo (DEBIT/CREDIT)
// @Tags         Journal Entries
// @Produce      json
// @Param        transaction_id  query     int     false  "Filtrar por ID de transacción"
// @Param        account_id      query     int     false  "Filtrar por ID de cuenta"
// @Param        debit_or_credit query     string  false  "Filtrar por tipo de asiento (DEBIT, CREDIT)"
// @Param        page            query     int     false  "Número de página (default: 1)"
// @Param        page_size       query     int     false  "Elementos por página (default: 50, máx: 100)"
// @Success      200  {object}  dtos.JournalEntryListResponse  "Lista paginada de asientos contables"
// @Failure      401  {object}  dtos.ErrorResponse             "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse             "Error interno del servidor"
// @Security     BearerAuth
// @Router       /journal-entries [get]
func (h *JournalEntryHandler) GetJournalEntries(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

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

// GetJournalEntriesByTransaction godoc
// @Summary      Asientos contables de una transacción
// @Description  Obtiene todos los asientos de débito y crédito generados para una transacción específica. Cada transacción produce exactamente 2 asientos (uno DEBIT y uno CREDIT) que deben sumar igual, garantizando el equilibrio de la doble partida
// @Tags         Journal Entries
// @Produce      json
// @Param        id   path      int                                                        true  "ID de la transacción"
// @Success      200  {object}  object{data=[]dtos.JournalEntryResponse,count=int}         "Asientos de la transacción"
// @Failure      400  {object}  dtos.ErrorResponse                                         "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                                         "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                                         "Error interno del servidor"
// @Security     BearerAuth
// @Router       /journal-entries/transaction/{id} [get]
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

// VerifyTransactionBalance godoc
// @Summary      Verificar equilibrio contable de una transacción
// @Description  Verifica que la suma de todos los asientos DEBIT sea igual a la suma de los asientos CREDIT para una transacción dada. Esto garantiza la regla fundamental de la contabilidad de doble partida: ∑ Débitos = ∑ Créditos
// @Tags         Journal Entries
// @Produce      json
// @Param        id   path      int                                     true  "ID de la transacción a verificar"
// @Success      200  {object}  object{data=dtos.BalanceVerificationResponse}  "Resultado de la verificación con totales de débito y crédito"
// @Failure      400  {object}  dtos.ErrorResponse                      "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                      "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                      "Error al ejecutar la verificación"
// @Security     BearerAuth
// @Router       /journal-entries/verify/{id} [get]
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
