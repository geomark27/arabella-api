package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
// @Summary      Listar usuarios
// @Description  Obtiene la lista de todos los usuarios registrados en el sistema. Endpoint de administración
// @Tags         Users
// @Produce      json
// @Success      200  {object}  object{data=[]dtos.UserResponseDTO,count=int}  "Lista de usuarios"
// @Failure      401  {object}  dtos.ErrorResponse                             "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse                             "Error interno del servidor"
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"count": len(users),
	})
}

// GetUserByID godoc
// @Summary      Obtener usuario por ID
// @Description  Obtiene el detalle de un usuario específico por su ID. No devuelve datos sensibles como el hash de contraseña
// @Tags         Users
// @Produce      json
// @Param        id   path      int                                    true  "ID del usuario"
// @Success      200  {object}  object{data=dtos.UserResponseDTO}      "Detalle del usuario"
// @Failure      400  {object}  dtos.ErrorResponse                     "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                     "No autenticado"
// @Failure      404  {object}  dtos.ErrorResponse                     "Usuario no encontrado"
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// CreateUser godoc
// @Summary      Crear usuario
// @Description  Crea un nuevo usuario en el sistema. Para el registro con autenticación JWT usa POST /auth/register en su lugar
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.CreateUserDTO                              true  "Datos del nuevo usuario"
// @Success      201   {object}  object{message=string,data=dtos.UserResponseDTO}  "Usuario creado exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse                              "Datos inválidos o email duplicado"
// @Failure      401   {object}  dtos.ErrorResponse                              "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse                              "Error interno del servidor"
// @Security     BearerAuth
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var dto dtos.CreateUserDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser godoc
// @Summary      Actualizar usuario
// @Description  Actualiza los datos de un usuario existente. Solo se modifican los campos enviados en el body (PATCH semántico)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int                 true  "ID del usuario"
// @Param        body  body      dtos.UpdateUserDTO  true  "Campos a actualizar"
// @Success      200   {object}  object{message=string,data=dtos.UserResponseDTO}  "Usuario actualizado exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse  "ID o datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse  "No autenticado"
// @Failure      404   {object}  dtos.ErrorResponse  "Usuario no encontrado"
// @Failure      500   {object}  dtos.ErrorResponse  "Error interno del servidor"
// @Security     BearerAuth
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var dto dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser godoc
// @Summary      Eliminar usuario
// @Description  Realiza un borrado lógico (soft delete) del usuario. Sus cuentas, transacciones y asientos contables se conservan en la base de datos
// @Tags         Users
// @Produce      json
// @Param        id   path      int                   true  "ID del usuario"
// @Success      200  {object}  dtos.SuccessResponse  "Usuario eliminado exitosamente"
// @Failure      400  {object}  dtos.ErrorResponse    "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse    "No autenticado"
// @Failure      404  {object}  dtos.ErrorResponse    "Usuario no encontrado"
// @Failure      500  {object}  dtos.ErrorResponse    "Error interno del servidor"
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
