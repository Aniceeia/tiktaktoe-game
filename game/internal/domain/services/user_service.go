package services

import (
	"context"
	"game/internal/domain/entities"
	"game/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) isUserExists(ctx context.Context, login string) bool {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err == nil && user != nil {
		return true
	}
	return false
}

func (s *UserService) Register(ctx context.Context, login, password string) (*entities.User, error) {
	if s.isUserExists(ctx, login) {
		return nil, entities.ErrUserExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errBcryptGenerate
	}
	user := entities.NewUser(login, string(hashedPassword))
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errCreateUser
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, login, password string) (*entities.User, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, entities.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, entities.ErrInvalidCredentials
	}

	return user, nil
}
