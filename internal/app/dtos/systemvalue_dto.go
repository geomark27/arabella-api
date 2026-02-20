package dtos

import "arabella-api/internal/app/models"

// SystemValueResponseDTO representa la respuesta de un valor de catálogo del sistema
type SystemValueResponseDTO struct {
	ID           uint    `json:"id"`
	CatalogType  string  `json:"catalog_type"`
	Value        string  `json:"value"`
	Label        string  `json:"label"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder int     `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}

// SystemValueListResponseDTO representa la respuesta de una lista de valores de catálogo
type SystemValueListResponseDTO struct {
	Data  []SystemValueResponseDTO `json:"data"`
	Count int                      `json:"count"`
}

// ToSystemValueResponse convierte un models.SystemValue a SystemValueResponseDTO
func ToSystemValueResponse(sv *models.SystemValue) SystemValueResponseDTO {
	return SystemValueResponseDTO{
		ID:           sv.ID,
		CatalogType:  sv.CatalogType,
		Value:        sv.Value,
		Label:        sv.Label,
		Description:  sv.Description,
		DisplayOrder: sv.DisplayOrder,
		IsActive:     sv.IsActive,
	}
}

// ToSystemValueResponseList convierte un slice de *models.SystemValue a []SystemValueResponseDTO
func ToSystemValueResponseList(svs []*models.SystemValue) []SystemValueResponseDTO {
	result := make([]SystemValueResponseDTO, len(svs))
	for i, sv := range svs {
		result[i] = ToSystemValueResponse(sv)
	}
	return result
}
