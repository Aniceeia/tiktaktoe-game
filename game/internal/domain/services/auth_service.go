package services

import (
	"context"
	"game/internal/domain/entities"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (*entities.User, error) {
	return s.userService.Register(ctx, login, password)
}

func (s *AuthService) Login(ctx context.Context, login, password string) (*entities.User, error) {
	return s.userService.Login(ctx, login, password)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	return s.userService.GetUserByID(ctx, userID)
}
