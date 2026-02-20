package dtos

// CreateUserDTO representa los datos para crear un usuario
type CreateUserDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"user_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

// UpdateUserDTO representa los datos para actualizar un usuario
type UpdateUserDTO struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	UserName  *string `json:"user_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UserResponseDTO representa la respuesta de un usuario (sin datos sensibles)
type UserResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ChangePassword struct {
	NewPassword string `json:"new_password" validate:"required"`
}
