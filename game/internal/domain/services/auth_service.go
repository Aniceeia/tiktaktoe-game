package services

import (
	"context"
	"errors"
	"game/internal/domain/entities"
	"game/internal/interfaces/http/dto"
)

type AuthService struct {
	userService *UserService
	jwtProvider *JwtProvider
}

func NewAuthService(userService *UserService, jwtProvider *JwtProvider) *AuthService {
	return &AuthService{userService: userService, jwtProvider: jwtProvider}
}

func (s *AuthService) Register(ctx context.Context, login, password string) (*entities.User, error) {
	return s.userService.Register(ctx, login, password)
}

// Authenticate: логин/пароль -> пара токенов
func (s *AuthService) Authenticate(ctx context.Context, req dto.JwtRequest) (*dto.JwtResponse, error) {
	user, err := s.userService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	access, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwtProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.JwtResponse{Type: "Bearer", AccessToken: access, RefreshToken: refresh}, nil
}

// RefreshAccessToken: по refresh -> новый access
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil || userID == "" {
		return nil, entities.ErrUnauthorized
	}
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	access, err := s.jwtProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.JwtResponse{Type: "Bearer", AccessToken: access, RefreshToken: refreshToken}, nil
}

// RefreshRefreshToken: по refresh -> новый refresh
func (s *AuthService) RefreshRefreshToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error) {
	userID, err := s.jwtProvider.ValidateRefreshToken(refreshToken)
	if err != nil || userID == "" {
		return nil, entities.ErrUnauthorized
	}
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	newRefresh, err := s.jwtProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}
	return &dto.JwtResponse{Type: "Bearer", AccessToken: "", RefreshToken: newRefresh}, nil
}

// ValidateAccessToken -> UUID пользователя
func (s *AuthService) ValidateAccessToken(token string) (string, error) {
	if token == "" {
		return "", errors.New("empty token")
	}
	return s.jwtProvider.ValidateAccessToken(token)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	return s.userService.GetUserByID(ctx, userID)
}
