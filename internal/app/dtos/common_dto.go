package dtos

// ErrorResponse represents a standard error response
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error   string `json:"error" example:"error message"`
	Details string `json:"details,omitempty" example:"additional error details"`
}

// SuccessResponse represents a standard success response with a message
// swagger:model SuccessResponse
type SuccessResponse struct {
	Message string `json:"message" example:"operation completed successfully"`
}

// PaginationMeta contains pagination metadata returned in list responses
// swagger:model PaginationMeta
type PaginationMeta struct {
	Total      int64 `json:"total" example:"100"`
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"20"`
	TotalPages int   `json:"total_pages" example:"5"`
}
