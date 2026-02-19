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

// GetTransactions godoc
// @Summary      Listar transacciones
// @Description  Obtiene una lista paginada de transacciones del usuario autenticado con filtros opcionales. Cada transacción pasa por el Motor Contable de doble partida
// @Tags         Transactions
// @Produce      json
// @Param        type         query     string  false  "Tipo de transacción (INCOME, EXPENSE, TRANSFER)"
// @Param        account_id   query     int     false  "Filtrar por ID de cuenta"
// @Param        category_id  query     int     false  "Filtrar por ID de categoría"
// @Param        page         query     int     false  "Número de página (default: 1)"
// @Param        page_size    query     int     false  "Elementos por página (default: 20, máx: 100)"
// @Success      200  {object}  dtos.TransactionListResponse  "Lista de transacciones paginada"
// @Failure      401  {object}  dtos.ErrorResponse            "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse            "Error interno del servidor"
// @Security     BearerAuth
// @Router       /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

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

// GetTransactionByID godoc
// @Summary      Obtener transacción por ID
// @Description  Obtiene el detalle completo de una transacción incluyendo sus asientos contables relacionados (account_from, account_to, categoría)
// @Tags         Transactions
// @Produce      json
// @Param        id   path      int                                       true  "ID de la transacción"
// @Success      200  {object}  object{data=dtos.TransactionResponse}     "Detalle de la transacción"
// @Failure      400  {object}  dtos.ErrorResponse                        "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                        "No autenticado"
// @Failure      404  {object}  dtos.ErrorResponse                        "Transacción no encontrada"
// @Security     BearerAuth
// @Router       /transactions/{id} [get]
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

// CreateTransaction godoc
// @Summary      Crear transacción
// @Description  Crea una nueva transacción procesándola a través del Motor Contable de doble partida. Genera automáticamente los asientos de débito y crédito y actualiza los saldos de las cuentas involucradas. Tipos: INCOME (requiere category_id), EXPENSE (requiere category_id), TRANSFER (requiere account_to_id)
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.CreateTransactionRequest                          true  "Datos de la transacción"
// @Success      201   {object}  object{message=string,data=dtos.TransactionResponse}  "Transacción creada y contabilizada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse                                    "Datos inválidos o regla de negocio violada"
// @Failure      401   {object}  dtos.ErrorResponse                                    "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse                                    "Error interno del servidor"
// @Security     BearerAuth
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dtos.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

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

// UpdateTransaction godoc
// @Summary      Actualizar transacción
// @Description  Actualiza los campos editables de una transacción existente (descripción, notas, fecha, estado de conciliación). El monto y las cuentas no pueden modificarse directamente; para ello, elimina la transacción y crea una nueva
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id    path      int                          true  "ID de la transacción"
// @Param        body  body      dtos.UpdateTransactionRequest  true  "Campos a actualizar"
// @Success      200   {object}  dtos.SuccessResponse         "Transacción actualizada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse           "ID o datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse           "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse           "Error interno del servidor"
// @Security     BearerAuth
// @Router       /transactions/{id} [put]
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

// DeleteTransaction godoc
// @Summary      Eliminar transacción (reversión contable)
// @Description  Elimina una transacción creando asientos de reversión en el diario contable (DEBIT↔CREDIT invertidos). Nunca borra registros del diario; garantiza la integridad del historial contable. Los saldos de las cuentas se restauran automáticamente
// @Tags         Transactions
// @Produce      json
// @Param        id   path      int                   true  "ID de la transacción a revertir"
// @Success      200  {object}  dtos.SuccessResponse  "Transacción revertida exitosamente"
// @Failure      400  {object}  dtos.ErrorResponse    "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse    "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse    "Error interno del servidor"
// @Security     BearerAuth
// @Router       /transactions/{id} [delete]
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
