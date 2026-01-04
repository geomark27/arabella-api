package dtos

// CreateUserDTO representa los datos para crear un usuario
type CreateUserDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

// UpdateUserDTO representa los datos para actualizar un usuario
type UpdateUserDTO struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
}

// UserResponseDTO representa la respuesta de un usuario (sin datos sensibles)
type UserResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
