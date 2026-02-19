package handlers

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryService services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategories godoc
// @Summary      Listar categorías
// @Description  Obtiene todas las categorías del usuario autenticado. Se puede filtrar por tipo (INCOME o EXPENSE)
// @Tags         Categories
// @Produce      json
// @Param        type  query     string  false  "Tipo de categoría (INCOME, EXPENSE)"
// @Success      200   {object}  object{data=[]dtos.CategoryResponse,count=int}  "Lista de categorías"
// @Failure      401   {object}  dtos.ErrorResponse                              "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse                              "Error interno del servidor"
// @Security     BearerAuth
// @Router       /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		respondUnauthorized(c)
		return
	}

	categoryType := c.Query("type")

	var categories []dtos.CategoryResponse
	var err error

	if categoryType != "" {
		categories, err = h.categoryService.GetByType(userID, categoryType)
	} else {
		categories, err = h.categoryService.GetByUser(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve categories",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"count": len(categories),
	})
}

// GetCategoryByID godoc
// @Summary      Obtener categoría por ID
// @Description  Obtiene el detalle de una categoría específica por su ID
// @Tags         Categories
// @Produce      json
// @Param        id   path      int                                      true  "ID de la categoría"
// @Success      200  {object}  object{data=dtos.CategoryResponse}       "Detalle de la categoría"
// @Failure      400  {object}  dtos.ErrorResponse                       "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse                       "No autenticado"
// @Failure      404  {object}  dtos.ErrorResponse                       "Categoría no encontrada"
// @Security     BearerAuth
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Category not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// CreateCategory godoc
// @Summary      Crear categoría
// @Description  Crea una nueva categoría para clasificar transacciones del usuario autenticado. Tipos válidos: INCOME (ingresos) o EXPENSE (gastos)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.CreateCategoryRequest                          true  "Datos de la nueva categoría"
// @Success      201   {object}  object{message=string,data=dtos.CategoryResponse}  "Categoría creada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse                                 "Datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse                                 "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse                                 "Error interno del servidor"
// @Security     BearerAuth
// @Router       /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req dtos.CreateCategoryRequest

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

	category, err := h.categoryService.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}

// UpdateCategory godoc
// @Summary      Actualizar categoría
// @Description  Actualiza el nombre o el estado activo de una categoría existente. Solo se modifican los campos enviados
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        id    path      int                         true  "ID de la categoría"
// @Param        body  body      dtos.UpdateCategoryRequest  true  "Campos a actualizar"
// @Success      200   {object}  dtos.SuccessResponse        "Categoría actualizada exitosamente"
// @Failure      400   {object}  dtos.ErrorResponse          "ID o datos inválidos"
// @Failure      401   {object}  dtos.ErrorResponse          "No autenticado"
// @Failure      500   {object}  dtos.ErrorResponse          "Error interno del servidor"
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	var req dtos.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.categoryService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
	})
}

// DeleteCategory godoc
// @Summary      Eliminar categoría
// @Description  Realiza un borrado lógico (soft delete) de una categoría. Las transacciones asociadas no se eliminan
// @Tags         Categories
// @Produce      json
// @Param        id   path      int                   true  "ID de la categoría"
// @Success      200  {object}  dtos.SuccessResponse  "Categoría eliminada exitosamente"
// @Failure      400  {object}  dtos.ErrorResponse    "ID inválido"
// @Failure      401  {object}  dtos.ErrorResponse    "No autenticado"
// @Failure      500  {object}  dtos.ErrorResponse    "Error interno del servidor"
// @Security     BearerAuth
// @Router       /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	if err := h.categoryService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete category",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
