package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	accountService services.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// GetAccounts godoc
// @Summary      Listar cuentas
// @Description  Obtiene todas las cuentas financieras del usuario autenticado (BANK, CASH, CREDIT_CARD, SAVINGS, INVESTMENT)
// @Tags         Accounts
// @Produce      json
// @Success      200  {object}  object{data=[]dtos.AccountResponseDTO,count=int}  "Lista de cuentas"
// @Failure      401  {object}  dtos.ErrorResponse                                "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                                "Error interno del servidor"
// @Security     BearerAuth
// @Router       /accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

	accounts, err := h.accountService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve accounts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  accounts,
		"count": len(accounts),
	})
}

// GetAccountByID godoc
// @Summary      Obtener cuenta por ID
// @Description  Obtiene el detalle completo de una cuenta financiera por su ID, incluyendo saldo y moneda
// @Tags         Accounts
// @Produce      json
// @Param        id   path      int                        true  "ID de la cuenta"
// @Success      200  {object}  object{data=dtos.AccountResponseDTO}  "Detalle de la cuenta"
// @Failure      400  {object}  dtos.ErrorResponse                    "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                    "No autenticado"
// @Failure      404  {object}  dtos.ErrorResponse                    "Cuenta no encontrada"
// @Security     BearerAuth
// @Router       /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	account, err := h.accountService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Account not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

// CreateAccount godoc
// @Summary      Crear cuenta
// @Description  Crea una nueva cuenta financiera para el usuario autenticado. El campo user_id se asigna automáticamente desde el JWT. Tipos válidos: BANK, CASH, CREDIT_CARD, SAVINGS, INVESTMENT
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.CreateAccountDTO                        true  "Datos de la nueva cuenta"
// @Success      201   {object}  object{message=string,data=dtos.AccountResponseDTO}  "Cuenta creada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse                           "Datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse                           "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse                           "Error interno del servidor"
// @Security     BearerAuth
// @Router       /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dtos.CreateAccountDTO

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
	req.UserID = userID

	account, err := h.accountService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"data":    account,
	})
}

// UpdateAccount godoc
// @Summary      Actualizar cuenta
// @Description  Actualiza los datos de una cuenta financiera existente. Solo se actualizan los campos enviados (PATCH semántico)
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "ID de la cuenta"
// @Param        body  body      dtos.UpdateAccountDTO     true  "Campos a actualizar"
// @Success      200   {object}  dtos.SuccessResponse      "Cuenta actualizada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse        "ID o datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse        "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse        "Error interno del servidor"
// @Security     BearerAuth
// @Router       /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	var req dtos.UpdateAccountDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.accountService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
	})
}

// DeleteAccount godoc
// @Summary      Eliminar cuenta
// @Description  Realiza un borrado lógico (soft delete) de una cuenta financiera. Los registros históricos se conservan
// @Tags         Accounts
// @Produce      json
// @Param        id   path      int                   true  "ID de la cuenta"
// @Success      200  {object}  dtos.SuccessResponse  "Cuenta eliminada exitosamente"
// @Failure      400  {object}  dtos.ErrorResponse    "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse    "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse    "Error interno del servidor"
// @Security     BearerAuth
// @Router       /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	if err := h.accountService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete account",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
