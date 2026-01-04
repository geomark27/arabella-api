package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// UserService defines the business logic for users
type UserService interface {
	CreateUser(dto *dtos.CreateUserDTO) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id uint, dto *dtos.UpdateUserDTO) (*models.User, error)
	DeleteUser(id uint) error
	ValidatePassword(user *models.User, password string) error
}

// userServiceImpl implements UserService
type userServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userServiceImpl) CreateUser(dto *dtos.CreateUserDTO) (*models.User, error) {
	// Check if user with email already exists
	existingUser, _ := s.userRepo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create default password hash (should be changed by user)
	// In production, this would come from the DTO
	defaultPassword := "changeme123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user model
	user := &models.User{
		UserName:     dto.Email, // Use email as username by default
		Email:        dto.Email,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID gets a user by ID
func (s *userServiceImpl) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail gets a user by email
func (s *userServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetAllUsers gets all users
func (s *userServiceImpl) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.FindAll()
}

// UpdateUser updates a user
func (s *userServiceImpl) UpdateUser(id uint, dto *dtos.UpdateUserDTO) (*models.User, error) {
	// Find existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if dto.FirstName != nil {
		user.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.Email != nil {
		// Check if new email is already in use by another user
		existingUser, _ := s.userRepo.FindByEmail(*dto.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use by another user")
		}
		user.Email = *dto.Email
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userServiceImpl) DeleteUser(id uint) error {
	// TODO: Check if user has any accounts or transactions before deleting
	// For now, we'll just do a soft delete
	return s.userRepo.Delete(id)
}

// ValidatePassword validates a user's password
func (s *userServiceImpl) ValidatePassword(user *models.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
