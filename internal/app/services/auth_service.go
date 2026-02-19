package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)

type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
	db         *gorm.DB
}

func NewAuthService(userRepo repositories.UserRepository, jwtService JWTService, db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
		db:         db,
	}
}

// Register creates a new user account
func (s *AuthService) Register(dto *dtos.RegisterAuthDTO) (*dtos.LoginResponseDTO, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        dto.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		UserName:     dto.Email, // Use email as username by default
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
	}, nil
}

// Login authenticates a user with email and password
func (s *AuthService) Login(dto *dtos.LoginDTO) (*dtos.LoginResponseDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
	}, nil
}

// RefreshToken generates new access and refresh tokens
func (s *AuthService) RefreshToken(dto *dtos.RefreshTokenDTO) (*dtos.RefreshTokenResponseDTO, error) {
	// Validate refresh token using the dedicated refresh-token validator
	// (signed with refreshSecret, only stores userID in Subject claim)
	userID, err := s.jwtService.ValidateRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user from the extracted userID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ChangePassword changes the user's password
func (s *AuthService) ChangePassword(userID uint, dto *dtos.ChangePasswordDTO) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.OldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}
