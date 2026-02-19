package handlers

import (
	"net/http"

	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register godoc
// @Summary      Registrar nuevo usuario
// @Description  Crea una nueva cuenta de usuario y devuelve tokens de acceso JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.RegisterAuthDTO     true  "Datos de registro"
// @Success      201   {object}  dtos.LoginResponseDTO    "Usuario registrado exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse       "Datos inválidos"
// @Failure      409   {object}  dtos.ErrorResponse       "El email ya está registrado"
// @Failure      500   {object}  dtos.ErrorResponse       "Error interno del servidor"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto dtos.RegisterAuthDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(&dto)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica al usuario con email y contraseña, devuelve tokens de acceso y refresco JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.LoginDTO          true  "Credenciales de acceso"
// @Success      200   {object}  dtos.LoginResponseDTO  "Autenticación exitosa"
// @Failure      400   {object}  dtos.ErrorResponse     "Datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse     "Credenciales incorrectas"
// @Failure      500   {object}  dtos.ErrorResponse     "Error interno del servidor"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto dtos.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(&dto)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary      Refrescar tokens
// @Description  Genera un nuevo par de tokens (acceso + refresco) a partir de un refresh token válido
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.RefreshTokenDTO          true  "Refresh token"
// @Success      200   {object}  dtos.RefreshTokenResponseDTO  "Tokens renovados exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse            "Datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse            "Refresh token inválido o expirado"
// @Failure      500   {object}  dtos.ErrorResponse            "Error interno del servidor"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var dto dtos.RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RefreshToken(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword godoc
// @Summary      Cambiar contraseña
// @Description  Cambia la contraseña del usuario autenticado verificando la contraseña actual
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.ChangePasswordDTO  true  "Contraseñas actual y nueva"
// @Success      200   {object}  dtos.SuccessResponse    "Contraseña actualizada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse      "Datos inválidos o contraseña actual incorrecta"
// @Failure      401   {object}  dtos.ErrorResponse      "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse      "Error interno del servidor"
// @Security     BearerAuth
// @Router       /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var dto dtos.ChangePasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ChangePassword(userID.(uint), &dto)
	if err != nil {
		if err == services.ErrInvalidPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
