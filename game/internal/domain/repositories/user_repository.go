package repositories

import (
	"context"
	"game/internal/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	IncreaseScore(ctx context.Context, id string) error
}
