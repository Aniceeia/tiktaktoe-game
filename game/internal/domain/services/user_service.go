package services

import (
	"context"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// GetUserByID получает пользователя по ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *AuthService) isUserExists(ctx context.Context, username string) bool {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && user != nil {
		return true
	}
	return false
}

func (s *AuthService) Register(ctx context.Context, username, password string) (*entities.User, error) {
	//check if user exists
	if s.isUserExists(ctx, username) {
		return nil, entities.ErrUserExists
	}
	//hash with base64
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errBcryptGenerate
	}
	//write username and hashed password
	user := entities.NewUser(username, string(hashedPassword))
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errCreateUser
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, entities.ErrInvalidCredentials
	}

	return user, nil
}
