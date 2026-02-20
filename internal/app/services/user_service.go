package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
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

	existingUser, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user model
	user := &models.User{
		UserName:     dto.UserName,
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
	// 1. Find existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Validate email uniqueness before applying changes (business rule, handled separately)
	if dto.Email != nil && *dto.Email != user.Email {
		existingUser, _ := s.userRepo.FindByEmail(*dto.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use by another user")
		}
	}

	// 3. Copy only the non-nil fields from DTO onto the model.
	//    IgnoreEmpty: true skips nil pointers and zero values, so existing
	//    model fields are never overwritten by absent DTO fields.
	if err := copier.CopyWithOption(user, dto, copier.Option{IgnoreEmpty: true}); err != nil {
		return nil, fmt.Errorf("failed to apply updates: %w", err)
	}

	// 4. Save changes
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
